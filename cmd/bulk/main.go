package main

import (
	"job-throtling-go-poc/pkg/csv"
	"job-throtling-go-poc/pkg/models"
	"job-throtling-go-poc/pkg/queue"
	"job-throtling-go-poc/pkg/services"
	"log"
	"os"
)

type Result struct {
	response *models.Address
	err      error
}

const filepath = "D:\\Proyectos\\Personales\\job-throtling-go-poc\\fixtures\\addresses.csv"

func callService(address models.Address, service *services.AddressService, results chan Result) queue.Runnable {
	return queue.NewRunnable(func() {
		resp, err := service.Post(address)
		results <- Result{
			response: resp,
			err:      err,
		}
	})
}

func main() {
	jobQueue := queue.NewQueue(2, 2)
	jobQueue.Start()

	reader := csv.NewReaderChannel[models.Address]()

	addressService := services.NewAddressService()

	addresses := make(chan models.Address)
	results := make(chan Result)

	fileReader, _ := os.Open(filepath)
	go reader.Read(fileReader, addresses)

	go func() {
		for address := range addresses {
			job := callService(address, addressService, results)
			jobQueue.AddJob(job)
		}
	}()

	done := make(chan bool)
	go func() {
		numberOfAddresses := 6
		for i := 0; i < numberOfAddresses; i++ {
			result := <-results
			log.Printf("%s\n", result.response.Name)
		}
		done <- true
	}()
	<-done
}
