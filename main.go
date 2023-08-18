package main

import (
	"gin-helloworld/app"
	"log"
)

func main() {
	app, err := app.Default()
	if err != nil {
		log.Fatalln(err)
	}
	app.Run()
}
