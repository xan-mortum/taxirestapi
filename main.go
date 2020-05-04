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

//логер. первый какой мне попался из тех которые умеют писать в несколько мест
var log = logging.MustGetLogger("taxirestapi")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func main() {
	//настраиваем логер на запись в файл и в консоль
	file, err := os.Create("log.log")
	backendFile := logging.NewLogBackend(file, "", 0)
	backendStdin := logging.NewLogBackend(os.Stdin, "", 0)

	backend2Formatter := logging.NewBackendFormatter(backendStdin, format)

	backend1Leveled := logging.AddModuleLevel(backendFile)
	backend1Leveled.SetLevel(logging.DEBUG, "")

	logging.SetBackend(backend1Leveled, backend2Formatter)

	//использую swagger server для генерации OpenAPI. в папке gen лежит сгерерированный код
	//в swagger.yml описание эндпоитов
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

	//немного философский вопрос где писать подключаемые структуры
	//лично мне больше нравиться писать не в main пакете а где то в другом месте
	//так сделал и тут. но, что бы не усложнять, засунул все одну папку
	//это тоже не совсем правильно, так структуры имеют доступ к приватным полям друг друга
	//но решил этот момент не прорабатывать потому что тестовое не об этом было

	//тут создаем экземпляр хранилица всех тех буков
	db := components.NewDB()
	//генерируем 50 штук
	db.Generate()
	//запускаем
	db.Start()


	//структура для хандлеров и добавление их к эндпоинтам
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
