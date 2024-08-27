package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pangami/gateway-service/middleware"
	"github.com/pangami/gateway-service/route"
	"github.com/pangami/gateway-service/util"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	e := route.Init()
	e.Use(middleware.ErrorHandlerMiddleware)

	data, err := util.Json.MarshalIndent(e.Routes(), "", "  ")
	if err != nil {
		panic(fmt.Sprint(err))
	}

	fmt.Println(string(data))
	e.Logger.Fatal(e.Start(":" + os.Getenv("APP_PORT")))
}
