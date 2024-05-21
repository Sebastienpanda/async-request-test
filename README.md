# Prompt Sebastien

Mon maitre de stage me demande de faire cet exercice en GO

Dans ce code de démo je ne veux pas de gestion réseau, http etc ...

je veux juste un main qui lance 1 thread qui va attendre une demande de la part du thread principal ... celui ci va attendre 2s puis demander au thread 1 ID qui lui servira à savoir où en est la demande et le thread devra lancer un autre thread qui informera le 1er de sa progression.

# Reponse Chagpt

Pour réaliser cet exercice en Go, nous allons créer un programme avec deux goroutines (qui sont l'équivalent des threads en Go). La première goroutine (thread) attend une demande du thread principal. Après 2 secondes, le thread principal envoie une demande à la première goroutine pour obtenir un ID. Ensuite, cette première goroutine lance une autre goroutine pour suivre sa progression et informe le thread principal.

Voici le code Go qui réalise cette tâche :

```go
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
```

Explication du code :

Main Function: La fonction principale lance une goroutine qui attend une demande d'ID et ensuite une autre goroutine pour suivre la progression.

WaitGroup: Utilisation de sync.WaitGroup pour attendre que toutes les goroutines terminent avant de quitter le programme.

Channel: Utilisation d'un canal (channel) pour la communication entre le thread principal et la première goroutine.

Goroutine 1: La première goroutine attend sur le canal pour recevoir un ID. Une fois qu'elle reçoit l'ID, elle lance une autre goroutine pour suivre la progression.

Progress Goroutine: La goroutine de progression affiche l'état de la progression à chaque seconde.

Lorsque vous exécutez ce programme, vous devriez voir que le thread principal envoie une demande après 2 secondes, puis la progression est affichée en plusieurs étapes.

Pour exécuter ce code, assurez-vous d'avoir Go installé sur votre machine et enregistrez le code dans un fichier nommé main.go. Ensuite, exécutez les commandes suivantes dans votre terminal :

```sh
go run main.go
```
