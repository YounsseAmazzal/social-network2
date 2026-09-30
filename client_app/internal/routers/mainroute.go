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

func Colleectrouters(handler *handlers.Handlers, middleware *middleware.Miiiddeware) http.Handler {
	servermux := http.DefaultServeMux
	serverstactic(servermux)
	Auhentication(servermux, handler, middleware)
	ProfileAndFollowers(servermux,handler,middleware)
	PooostsandComments(servermux,handler,middleware)
	Groupsandevents(servermux,handler,middleware)
	newratelimeter := utils.NewRateLimiter(time.Minute)
	return newratelimeter.ApplyMiddleware("normalapi",servermux)
}
func serverstactic(muxserve *http.ServeMux) {
	uploadDir := utils.Soccialuplodfile
	prefix := "/socialup/"

	muxserve.HandleFunc(prefix, func(ress http.ResponseWriter, reeq *http.Request) {
		//avoid p t
		cleeanPth := filepath.Clean(strings.TrimPrefix(reeq.URL.Path, prefix))
		if cleeanPth == "." || cleeanPth == "/" || strings.HasSuffix(reeq.URL.Path, "/") {
			http.NotFound(ress, reeq)
			return
		}

		fullPath := filepath.Join(uploadDir, cleeanPth)

		info, err := os.Stat(fullPath)
		if err != nil || info.IsDir() {
			http.NotFound(ress, reeq)
			return
		}
		http.ServeFile(ress, reeq, fullPath)
	})
}
