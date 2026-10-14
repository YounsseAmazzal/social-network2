package utils

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type visitorState struct {
	lastSeenAt     time.Time
	requestHistory []time.Time
}

type RateLimiter struct {
	windowDuration  time.Duration
	thresholds      map[string]int
	visitorRegistry map[string]*visitorState
	registryMutex   sync.RWMutex
}

func NewRateLimiter(duration time.Duration) *RateLimiter {
	limiter := &RateLimiter{
		windowDuration: duration,
		thresholds: map[string]int{
			"normalapi":      1000,
			"authintication": 200,
		},
		visitorRegistry: make(map[string]*visitorState),
	}

	go limiter.evictExpiredVisitors()

	return limiter
}

func (rl *RateLimiter) ApplyMiddleware(category string, next http.Handler) http.Handler {
	maxRequests, exists := rl.thresholds[category]

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions || !exists {
			next.ServeHTTP(w, r)
			return
		}

		clientIP := extractIP(r)

		if !rl.isWithinLimit(clientIP, maxRequests) {
			w.Header().Set("X-RateLimit-Config", category)
			http.Error(w, "Rate limit exceeded. Too many requests.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// insure if the user passed the num on allowed requests
func (rl *RateLimiter) isWithinLimit(ip string, maxAllowed int) bool {
	rl.registryMutex.Lock()
	defer rl.registryMutex.Unlock()

	now := time.Now()
	expirationBoundary := now.Add(-rl.windowDuration)
	visitor, active := rl.visitorRegistry[ip]
	if !active {
		rl.visitorRegistry[ip] = &visitorState{
			lastSeenAt:     now,
			requestHistory: []time.Time{now},
		}
		return true
	}

	visitor.lastSeenAt = now

	var recentRequests []time.Time
	for _, timestamp := range visitor.requestHistory {
		if timestamp.After(expirationBoundary) {
			recentRequests = append(recentRequests, timestamp)
		}
	}

	if len(recentRequests) >= maxAllowed {
		visitor.requestHistory = recentRequests
		return false
	}

	visitor.requestHistory = append(recentRequests, now)
	return true
}

// delete the veseters not active to clean the memory
func (rl *RateLimiter) evictExpiredVisitors() {
	ticker := time.NewTicker(rl.windowDuration * 2)
	for currentTime := range ticker.C {
		rl.registryMutex.Lock()
		expirationBoundary := currentTime.Add(-rl.windowDuration)

		for ip, state := range rl.visitorRegistry {
			if state.lastSeenAt.Before(expirationBoundary) {
				delete(rl.visitorRegistry, ip)
			}
		}
		rl.registryMutex.Unlock()
	}
}

// bring the exact ip of the visetor
func extractIP(r *http.Request) string {
	// porrroxy first
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
