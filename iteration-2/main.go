package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Structure représentant une tâche avec un ID et une progression.
type Task struct {
	ID       int
	Progress int
}

// Fonction de l'opération longue et aléatoire.
func longOperation(task *Task, progressChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for task.Progress < 100 {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond) // Simulation de temps de traitement aléatoire
		task.Progress += rand.Intn(20)
		if task.Progress > 100 {
			task.Progress = 100
		}
		progressChan <- task.Progress
	}
}

// Fonction du thread scheduler.
func scheduler(requestChan <-chan struct{}, idChan chan<- int, progressChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	taskID := 0
	for range requestChan {
		taskID++
		task := &Task{ID: taskID, Progress: 0}
		idChan <- task.ID
		go longOperation(task, progressChan, wg)
	}
}

func main() {
	requestChan := make(chan struct{})
	idChan := make(chan int)
	progressChan := make(chan int)
	var wg sync.WaitGroup

	// Démarrage du scheduler
	wg.Add(1)
	go scheduler(requestChan, idChan, progressChan, &wg)

	// Attendre 2 secondes pour que le scheduler soit bien démarré
	time.Sleep(2 * time.Second)

	// Envoi de la demande au scheduler
	requestChan <- struct{}{}

	// Récupération de l'ID de la tâche
	taskID := <-idChan
	fmt.Printf("ID de la tâche reçue: %d\n", taskID)

	// Suivi de la progression de la tâche toutes les secondes
	progress := 0
	for progress < 100 {
		time.Sleep(1 * time.Second)
		progress = <-progressChan
		fmt.Printf("Progression de la tâche %d: %d%%\n", taskID, progress)
	}

	// Dire merci une fois la tâche terminée
	fmt.Println("Merci au scheduler")

	// Fermeture des canaux et attente de la fin des goroutines
	close(requestChan)
	wg.Wait()
}
