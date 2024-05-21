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
	// Démarre le thread qui attend les requêtes
	requestChan := make(chan Request)
	go requestHandler(requestChan)

	// Attend 2 secondes que le thread soit bien démarré
	time.Sleep(2 * time.Second)

	// Envoie une requête au thread et attend la réponse
	request := Request{}
	requestChan <- request
	response := <-requestChan

	// Demande la progression de la requête toutes les secondes jusqu'à ce qu'elle soit terminée
	for !response.Completed {
		fmt.Printf("Progression de la requête %d : %d%%\n", response.ID, response.Progress)
		time.Sleep(1 * time.Second)
		response = <-requestChan
	}

	fmt.Println("Merci au thread scheduler !")
}

func requestHandler(requestChan chan Request) {
	var requestID int = 1

	for {
		// Attend une requête du thread principal
		request := <-requestChan

		// Crée un ID pour la requête
		request.ID = requestID
		requestID++

		// Crée un canal pour communiquer avec le thread de traitement
		progressChan := make(chan int)

		// Démarre le thread de traitement
		go processRequest(request.ID, progressChan)

		// Envoie les mises à jour de progression au thread principal
		for progress := range progressChan {
			request.Progress = progress
			if progress == 100 {
				request.Completed = true
				break
			}
			requestChan <- request
		}

		// Renvoie la requête terminée au thread principal
		requestChan <- request
	}
}

func processRequest(id int, progressChan chan int) {
	// Simule une opération longue et aléatoire
	progress := 0
	for progress < 100 {
		progress += rand.Intn(5) + 1
		progressChan <- progress
		time.Sleep(time.Duration(rand.Intn(500)+500) * time.Millisecond)
	}
	progressChan <- 100
}
