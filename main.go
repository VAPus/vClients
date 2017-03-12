package main

import (
	"log"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"fmt"
)

type applicationHandler struct {
	Config *config
}

func main() {
	// Create main application handler
	appHandler := applicationHandler{}

	cfg, err := loadConfigurationFile("config.toml")

	if err != nil {
		log.Fatal(err)
	}

	appHandler.Config = cfg

	// Create router
	router := httprouter.New()

	// Listen
	if appHandler.Config.SSL.Enabled {

		if err := http.ListenAndServeTLS(
			fmt.Sprintf("%v:%v", appHandler.Config.Host, appHandler.Config.Port),
			appHandler.Config.SSL.Cert,
			appHandler.Config.SSL.Key,
			router,
		); err != nil {
			log.Fatal(err)
		}
	}

	if err := http.ListenAndServe(fmt.Sprintf("%v:%v", appHandler.Config.Host, appHandler.Config.Port), router); err != nil {
		log.Fatal(err)
	}
}

func wrapper() {

}