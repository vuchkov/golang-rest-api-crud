package main

import (
	"gitlab.com/devskiller-tasks/rest-api-blog-golang/bootstrap"
	"log"
)

func main() {
	defaultPort := 8080
	if err := bootstrap.Init(defaultPort); err != nil {
		log.Fatalf("Service will be shutdown because error ocurred:  %+v", err.Error())
	}
}
