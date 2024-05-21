package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Request struct {
	ID        int
	Progress  int
	Completed bool
}

func main() {
	requestChan := make(chan Request)
	responseChan := make(chan Request)
	progressRequestChan := make(chan int)
	progressResponseChan := make(chan Request)
	doneChan := make(chan bool)

	go requestHandler(requestChan, responseChan, progressRequestChan, progressResponseChan, doneChan)

	// Envoie une requête au thread et récupère l'ID
	request := Request{}
	requestChan <- request
	response := <-responseChan
	fmt.Printf("Requête envoyée, ID : %d\n", response.ID)

	// Demande la progression de la requête toutes les secondes jusqu'à ce qu'elle soit terminée
	go func() {
		for !response.Completed {
			time.Sleep(1 * time.Second)
			progressRequestChan <- response.ID
			response = <-progressResponseChan
			fmt.Printf("Progression de la requête %d : %d%%\n", response.ID, response.Progress)
		}
		doneChan <- true
	}()

	<-doneChan
	fmt.Println("Merci au thread scheduler !")
}

func requestHandler(requestChan chan Request, responseChan chan Request, progressRequestChan chan int, progressResponseChan chan Request, doneChan chan bool) {
	var requestID int = 1
	tasks := make(map[int]*Request)

	for {
		select {
		case request := <-requestChan:
			request.ID = requestID
			requestID++

			progressChan := make(chan int)
			go processRequest(request.ID, progressChan)

			tasks[request.ID] = &request
			responseChan <- request

			go func(id int, progressChan chan int) {
				for progress := range progressChan {
					tasks[id].Progress = progress
					if progress == 100 {
						tasks[id].Completed = true
					}
					select {
					case progressResponseChan <- *tasks[id]:
					default:
						fmt.Println("Aucun récepteur pour progressResponseChan")
					}
				}
			}(request.ID, progressChan)

		case id := <-progressRequestChan:
			if task, exists := tasks[id]; exists {
				progressResponseChan <- *task
			}

		case <-doneChan:
			return
		}
	}
}

func processRequest(id int, progressChan chan int) {
	defer close(progressChan)
	progress := 0
	for progress < 100 {
		progress += rand.Intn(5) + 1
		progressChan <- progress
		time.Sleep(time.Duration(rand.Intn(500)+500) * time.Millisecond)
	}
	progressChan <- 100
}
