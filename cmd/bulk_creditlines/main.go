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
	response *models.CreditLine
	err      error
}

const filepath = "D:\\Proyectos\\Personales\\job-throtling-go-poc\\fixtures\\create_line_example.csv"

func callService(creditLine models.CreditLine, service *services.CreditLineService, results chan Result) queue.Runnable {
	return queue.NewRunnable(func() {
		resp, err := service.Post(creditLine)
		results <- Result{
			response: resp,
			err:      err,
		}
	})
}

func main() {
	jobQueue := queue.NewQueue(5, 5)
	jobQueue.Start()

	reader := csv.NewReaderChannel[models.CreditLine]()

	service := services.NewCreditLineService()

	creditLines := make(chan models.CreditLine, 10)
	results := make(chan Result, 10)
	totalCreditLines := make(chan int)

	fileReader, _ := os.Open(filepath)

	// Read file to the channel creditLines
	go reader.Read(fileReader, creditLines, totalCreditLines)

	// When receiving records in the channel, add api calls to the queue
	go func() {
		for creditLine := range creditLines {
			job := callService(creditLine, service, results)
			jobQueue.AddJob(job)
		}
	}()

	var processed = 0
	var total = 0
	var receivedTotal = false
	for {
		select {
		case total = <-totalCreditLines:
			receivedTotal = true
		case <-results:
			processed += 1
		}

		if processed%100 == 0 {
			log.Printf("Processed = %d\n", processed)
		}

		if receivedTotal && processed == total {
			log.Println("Done!")
			break
		}
	}
}
