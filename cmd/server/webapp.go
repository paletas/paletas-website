package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/paletas/paletas_website/internal/server/webapp"
)

func main() {
	portArg := flag.Int("port", 8080, "port to listen on")
	basePathArg := flag.String("basepath", "", "optional base path of the application")

	flag.Parse()

	basePath := *basePathArg
	if basePath == "" {
		var err error
		basePath, err = os.Getwd()
		if err != nil {
			fmt.Println("Error getting application base path:", err)
			return
		}
	} else {
		os.Chdir(basePath)
	}

	fmt.Printf("Application base path is %v\n", basePath)

	webapp := webapp.NewWebApp()

	listenPort := fmt.Sprintf(":%d", *portArg)
	fmt.Printf("Listening on port %s\n", listenPort)

	webapp.Start(listenPort)
}
