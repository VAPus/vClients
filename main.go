package main

import (
	"log"
	"github.com/VAPus/vClients/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"fmt"
	"html/template"
)

var appHandler = applicationHandler{}

type applicationHandler struct {
	Config *util.Config
	Response http.ResponseWriter
	Request *http.Request
	Template *template.Template
	Clients []util.Client
}

type controller func(handler applicationHandler) (int, error)

func main() {
	cfg, err := util.LoadConfigurationFile("config.toml")

	if err != nil {
		log.Fatal(err)
	}

	appHandler.Config = cfg

	tpl, err := util.LoadTemplates("views")

	if err != nil {
		log.Fatal(err)
	}

	appHandler.Template = tpl

	list, err := util.GetClientList("clients")

	if err != nil {
		log.Fatal(err)
	}

	appHandler.Clients = list

	log.Println(list)

	// Create router
	router := httprouter.New()
	router.ServeFiles("/public/*filepath", http.Dir("public"))
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
	app.Template.ExecuteTemplate(app.Response, "main.html", nil)

	return 200, nil
}

func wrapper(c controller) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		// Reload templates on development mode
		if appHandler.Config.Mode == "dev" {

			tpl, err := util.LoadTemplates("views")

			if err != nil {
				log.Fatal(err)
			}

			appHandler.Template = tpl
		}

		// Create controller handler
		handler := applicationHandler{
			Config: appHandler.Config,
			Request: r,
			Response: w,
			Template: appHandler.Template,
		}

		// Call controller
		status, err := c(handler)

		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}
	}
}