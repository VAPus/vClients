package main

import (
	"log"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"fmt"
)

var appHandler = applicationHandler{}

type applicationHandler struct {
	Config *config
	Response http.ResponseWriter
	Request *http.Request
}

type controller func(handler applicationHandler) (int, error)

func main() {
	cfg, err := loadConfigurationFile("config.toml")

	if err != nil {
		log.Fatal(err)
	}

	appHandler.Config = cfg

	// Create router
	router := httprouter.New()
	router.GET("/", wrapper(clientList))

	log.Printf("Listening on %v:%v", appHandler.Config.Host, appHandler.Config.Port)

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

func clientList(app applicationHandler) (int, error) {
	app.Response.Write([]byte("1212"))

	return 200, nil
}

func wrapper(c controller) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		// Create controller handler
		handler := applicationHandler{
			Config: appHandler.Config,
			Request: r,
			Response: w,
		}

		// Call controller
		status, err := c(handler)

		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		w.WriteHeader(status)
	}
}