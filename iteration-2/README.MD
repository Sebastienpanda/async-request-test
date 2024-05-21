# Prompt Sebastien

faire un programme en GO qui :

1- démarre un thread qui va attendre une demande d'un thread principal
2- attends 2 secondes que le thread soit bien démarré
3- envoi au thread (et là il faudra me dire comment ils communiquent, mais typiquement ça serait un channel) une requête au thread et attend un id en réponse
4- quand le thread reçoit la notification, il créé un id, initialise une sorte de structure (à déterminer et expliciter) et lance un thread qui va faire une opération longue (et aléatoire), mais qui indiquera au thread 1 (le "scheduler" en quelques sortes, ça aussi c'est une notion pour laquelle je veux que vous questionnez ChatGPT) ... sa progression ... disons que c'est un pourcentage de progression qui va de 0 à 100. Et qui s'arrête quand il atteint 100.
5- Une fois qu'il a reçu son ID ... le thread principal va demander toutes les secondes où se trouve la progression de la demande ... et quand c'est fini, alors il dira "merci" au thread "scheduler"

il est interdit d'utiliser un crate ou une librairie extérieure

# Reponse Poe

Voici un programme en Go qui répond à vos exigences :

```go
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
```

Explication du programme :

Le programme démarre un thread requestHandler qui attend les requêtes du thread principal.
Après 2 secondes, le thread principal envoie une requête au requestHandler via un canal requestChan.
Le requestHandler crée un ID unique pour la requête, crée un canal progressChan pour communiquer avec le thread de traitement, et démarre le thread de traitement processRequest.
Le processRequest simule une opération longue et aléatoire en mettant à jour la progression de 0 à 100 dans le canal progressChan.
Le requestHandler envoie les mises à jour de progression au thread principal via le canal requestChan.
Une fois que la progression atteint 100%, le thread principal affiche "Merci au thread scheduler !"
La communication entre les threads se fait via des canaux Go. Le requestHandler joue le rôle de "scheduler" en recevant les requêtes du thread principal, en les distribuant aux threads de traitement, et en renvoyant les résultats au thread principal.

```sh
go run main.go
```
