package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// Utilisation d'un WaitGroup pour attendre que les goroutines terminent leur travail
	var wg sync.WaitGroup

	// Channel pour communiquer l'ID de progression
	idChannel := make(chan int)

	// Goroutine pour gérer les demandes du thread principal
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Goroutine 1: En attente d'une demande du thread principal...")

		// Attente de la demande d'ID
		id := <-idChannel
		fmt.Printf("Goroutine 1: Reçu demande d'ID: %d\n", id)

		// Lancement d'une nouvelle goroutine pour informer de la progression
		wg.Add(1)
		go func(progressID int) {
			defer wg.Done()
			for i := 0; i <= 100; i += 20 {
				fmt.Printf("Progression (ID %d): %d%%\n", progressID, i)
				time.Sleep(1 * time.Second) // Simulation de la progression
			}
		}(id)
	}()

	// Thread principal attend 2 secondes avant d'envoyer une demande d'ID
	time.Sleep(2 * time.Second)
	fmt.Println("Thread principal: Envoi de la demande d'ID...")
	idChannel <- 1

	// Attendre que toutes les goroutines terminent
	wg.Wait()
	fmt.Println("Tous les threads ont terminé.")

}
