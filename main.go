package main

import (
	"github.com/arganaphangquestian/gin-jwt/repository"
	"github.com/arganaphangquestian/gin-jwt/route"
)

func main() {
	repository := repository.New()
	app := route.New(repository)
	app.Run()
}
