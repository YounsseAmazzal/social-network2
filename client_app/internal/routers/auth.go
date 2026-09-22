package routers

import (
	"fmt"
	"net/http"
)

func AllRouttesserve(handler *http.Handler)  {
	muxxx := http.NewServeMux()
	fmt.Println("server mux", muxxx)
}
