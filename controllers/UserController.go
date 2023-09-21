package controllers

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"os"
)

type KafkaSendable interface {
	GetKafkaCredentials() string
}
type FIOdata struct {
	Name                     string `json:"name" validate:"required"`
	Surname                  string `json:"surname" validate:"required"`
	Patronymic               string `json:"patronymic"`
	Age, Gender, Nationality string
}

type FIOErrorData struct {
	FIOdata
	Error error
}

func (f *FIOdata) GetKafkaCredentials() string {
	return os.Getenv("FIO_TOPIC_NAME")
}

func (f *FIOErrorData) GetKafkaCredentials() string {
	return os.Getenv("FIO_FAILED_TOPIC_NAME")
}

func WorkerFIO(FIOjobCh chan string, FIOFailedJobCh chan *FIOErrorData) {

	fioStruct := &FIOdata{}
	for job := range FIOjobCh {
		json.Unmarshal([]byte(job), fioStruct)

		if err := fioStruct.validateFIO(); err != nil {

			FIOFailedJobCh <- &FIOErrorData{
				*fioStruct,
				err,
			}
		}
	}
}

func WorkerFailedFIO(FIOFailedJobCh chan *FIOErrorData, KafkaSendableChannel chan KafkaSendable) {
	for job := range FIOFailedJobCh {
		KafkaSendableChannel <- job
	}
}

func (c *FIOdata) validateFIO() error {
	v := validator.New()
	err := v.Struct(c)
	return err
}
