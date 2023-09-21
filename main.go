package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"testjun/controllers"
	"testjun/controllers/web"
	fioKafka2 "testjun/utils/kafka"
	util "testjun/utils/postgres"
	"time"
)

func main() {
	err := godotenv.Load()

	db := util.GetDbInstance()
	util.AutoMigrate(db)
	initWebService()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	initKafka()

	time.Sleep(time.Second * 100)
}

func initWebService() {
	r := gin.Default()
	r.PUT("/update", web.Update)
	r.DELETE("/delete", web.Delete)
	r.GET("/get", web.Get)

	r.POST("/create", web.Create)
	r.Run(":8080")
}

func initKafka() {

	// pretending it to be live data

	FIOjobCh := make(chan string)
	FIOFailedJobCh := make(chan *controllers.FIOErrorData)
	KafkaSendableChannel := make(chan controllers.KafkaSendable)
	go getFIOJsonStringData(KafkaSendableChannel)
	go controllers.WorkerFIO(FIOjobCh, FIOFailedJobCh)

	go fioKafka2.Produce(KafkaSendableChannel)

	go fioKafka2.Consume(FIOjobCh)
	go controllers.WorkerFailedFIO(
		FIOFailedJobCh,
		KafkaSendableChannel,
	)

}

func getFIOJsonStringData(kafkachan chan controllers.KafkaSendable) []controllers.FIOdata {

	var u = []controllers.FIOdata{
		{
			Name:       "Dmitriy",
			Surname:    "Ushakov",
			Patronymic: "Vasilevich",
		},
		{
			Name:       "Testing",
			Surname:    "SomeSurname",
			Patronymic: "Sergeevich",
		},
		{
			Name:       "",
			Surname:    "testWithoutname",
			Patronymic: "Sergeevich",
		},
		{
			Name:       "testWithoutSurname",
			Patronymic: "Sergeevich",
		},
	}
	for _, job := range u {
		kafkachan <- &job
	}

	return u
}

func initWorkers() {

}
