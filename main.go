package main

import (
	"gin-helloworld/app"
	_ "gin-helloworld/docs"
	"log"
)

func main() {
	app, err := app.Default()
	if err != nil {
		log.Fatalln(err)
	}
	app.Run()
}
