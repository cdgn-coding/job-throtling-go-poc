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
	jobQueue := queue.NewQueue(3, 3)
	jobQueue.Start()

	reader := csv.NewReaderChannel[models.Address]()

	addressService := services.NewAddressService()

	addresses := make(chan models.Address)
	totalAddresses := make(chan int)
	results := make(chan Result)

	fileReader, _ := os.Open(filepath)

	// Read file to the channel addresses
	go reader.Read(fileReader, addresses, totalAddresses)

	// When receiving records in the channel, add api calls to the queue
	go func() {
		for address := range addresses {
			job := callService(address, addressService, results)
			jobQueue.AddJob(job)
		}
	}()

	// Read all results and write to done channel in the end
	done := make(chan bool)
	go func() {
		var processed = 0
		var total = 0
		var receivedTotal = false
		for {
			select {
			case total = <-totalAddresses:
				receivedTotal = true
			case <-results:
				processed += 1
			}

			if receivedTotal && processed == total {
				log.Println("Done!")
				break
			}
		}
		done <- true
	}()
	<-done
}
