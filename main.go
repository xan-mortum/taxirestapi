package main

import (
	"github.com/go-openapi/loads"
	"github.com/op/go-logging"
	"github.com/xan-mortum/taxirestapi/components"
	"github.com/xan-mortum/taxirestapi/gen/restapi"
	"github.com/xan-mortum/taxirestapi/gen/restapi/operations"
	"github.com/xan-mortum/taxirestapi/handlers"
	"os"
)

const Port = 8082

var log = logging.MustGetLogger("taxirestapi")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func main() {
	file, err := os.Create("log.log")
	backendFile := logging.NewLogBackend(file, "", 0)
	backendStdin := logging.NewLogBackend(os.Stdin, "", 0)

	backend2Formatter := logging.NewBackendFormatter(backendStdin, format)

	backend1Leveled := logging.AddModuleLevel(backendFile)
	backend1Leveled.SetLevel(logging.DEBUG, "")

	logging.SetBackend(backend1Leveled, backend2Formatter)

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatal(err)
	}
	api := operations.NewTaxirestapiAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer func() {
		err := server.Shutdown()
		if err != nil {
			log.Fatal(err)
		}
	}()

	db := components.NewDB()
	db.Generate()
	db.Start()

	taxiHandlers := handlers.NewTaxiHandler(log, db)
	api.RequestHandler = operations.RequestHandlerFunc(taxiHandlers.RequestHandler)
	api.RequestsHandler = operations.RequestsHandlerFunc(taxiHandlers.RequestsHandler)

	api.ServerShutdown = func() {
		db.Stop()
	}

	server.Port = Port
	err = server.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
