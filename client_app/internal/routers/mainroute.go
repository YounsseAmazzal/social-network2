package routers

import (
	"net/http"
	"os"
	"path/filepath"
	"social/internal/handlers"
	"social/internal/middleware"
	"social/internal/utils"
	"strings"
	"time"
)

func Colleectrouters(handler *handlers.Handlers, middleware *middleware.Miiiddeware,wsm *handlers.WSManager) http.Handler {
	servermux := http.DefaultServeMux
	serverstactic(servermux)
	Auhentication(servermux, handler, middleware)
	ProfileAndFollowers(servermux,handler,middleware)
	PooostsandComments(servermux,handler,middleware)
	Groupsandevents(servermux,handler,middleware)
	WsChatAndNothification(servermux,wsm,middleware ,handler)
	newratelimeter := utils.NewRateLimiter(time.Minute)
	return newratelimeter.ApplyMiddleware("normalapi",servermux)
}
func serverstactic(muxserve *http.ServeMux) {
	serveStaticDir(muxserve, "/socialup/", utils.Soccialuplodfile)

	serveStaticDir(muxserve, "/public/", "../client_web/public/") 
    
    // serveStaticDir(muxserve, "/icons/", "./public/icons/")
}
func serveStaticDir(mux *http.ServeMux, prefix string, dirPath string) {
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		// avoid p t
		cleanPath := filepath.Clean(strings.TrimPrefix(r.URL.Path, prefix))
		
		if cleanPath == "." || cleanPath == "/" || strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		fullPath := filepath.Join(dirPath, cleanPath)

		info, err := os.Stat(fullPath)
		if err != nil || info.IsDir() {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, fullPath)
	})
}
