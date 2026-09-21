package app

import (
	"fmt"
	"social/internal/config"
	"social/internal/database"
)

func Server() {
	configg := config.Takeconfig("config.json")
	fmt.Println("that is the config", configg.DBPath)
	fmt.Println("that is the Port", configg.Port)
	fmt.Println("that is the WebOrigin", configg.WebOrigin)
	// // run the serveeer
	daaatabase := database.NeewDb(configg.DBPath)
	println("hiiii",daaatabase)
	defer daaatabase.Database.Close()
	// log.Println("the seerver will run on http:localhost:8080")
	// log.Fatal(http.ListenAndServe(":"+configg.Port, nil))
}
