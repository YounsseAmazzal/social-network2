package app

import (
	"fmt"
	"log"
	"net/http"
	"social/internal/config"
	"social/internal/database"
	"social/internal/handlers"
	"social/internal/middleware"
	"social/internal/repositories"
	"social/internal/routers"
	"social/internal/services"
)

func Server() {
	configg := config.Takeconfig("config.json")
	fmt.Println("that is the config", configg.DBPath)
	fmt.Println("that is the Port", configg.Port)
	fmt.Println("that is the WebOrigin", configg.WebOrigin)
	// // run the serveeer
	daaatabase := database.NeewDb(configg.DBPath)
	println("hiiii", daaatabase)
	// defer daaatabase.Database.Close()
	newrepo := repositories.Neewwrepo(daaatabase)
	fmt.Println("that is the newrepo", daaatabase)
	newservice := services.NewService(newrepo)
	fmt.Println("that is the services", newservice)
	wsmaanager := handlers.Neeewsmanger(newservice)
	fmt.Println("that is the wssocket", wsmaanager)
	handler := handlers.NeewHandler(newservice, wsmaanager)
	fmt.Println("that is the handler", handler)
	miidl := middleware.Neewmidd(newservice)
	fmt.Println("that is the midleware", miidl)
	routtes := routers.Colleectrouters(handler, miidl)
	mainHandler:=miidl.Corsmiddwaare(configg.WebOrigin)(routtes)
	log.Println("the seerver will run on http:localhost:8080")
	if err := http.ListenAndServe(":"+configg.Port, mainHandler); err != nil {
		log.Fatal("error on server")
	}
}
