package main

import (
	"infrapi"
	"log"
)

func main() {
	log.Println("Loading global config")
	err := infrapi.LoadGlobalConfig()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Connecting to Redis")
	err = infrapi.ConnectRedis()
	if err != nil {
		log.Println(err)
		return
	}

	err = infrapi.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
