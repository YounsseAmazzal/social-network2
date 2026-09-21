package config

import (
	"encoding/json"
	"log"
	"os"
	"social/internal/models"
)

func Takeconfig(paath string)*models.Config{
	readdata,err:=os.ReadFile(paath)
	if err!=nil{
		log.Fatal("Cann not read the config file ",err)
	}

	var configg models.Config
	if err:=json.Unmarshal(readdata,&configg);err!=nil{
		log.Fatal("Caan not paarse the config file ",err)
	}
	if configg.Port==""{
		configg.Port="8080"
	}
	if configg.WebOrigin==""{
		configg.WebOrigin="http://localhost:3000"
	}
	if configg.DBPath==""{
		log.Fatal("the data base path it's very importent ",err)
	}
	return &configg
}