# Prompt Sebastien

faire un programme en GO qui :

1- démarre un thread qui va attendre une demande d'un thread principal
2- attends 2 secondes que le thread soit bien démarré
3- envoi au thread (et là il faudra me dire comment ils communiquent, mais typiquement ça serait un channel) une requête au thread et attend un id en réponse
4- quand le thread reçoit la notification, il créé un id, initialise une sorte de structure (à déterminer et expliciter) et lance un thread qui va faire une opération longue (et aléatoire), mais qui indiquera au thread 1 (le "scheduler" en quelques sortes, ça aussi c'est une notion pour laquelle je veux que vous questionnez ChatGPT) ... sa progression ... disons que c'est un pourcentage de progression qui va de 0 à 100. Et qui s'arrête quand il atteint 100.
5- Une fois qu'il a reçu son ID ... le thread principal va demander toutes les secondes où se trouve la progression de la demande ... et quand c'est fini, alors il dira "merci" au thread "scheduler"

il est interdit d'utiliser un crate ou une librairie extérieure

# Reponse Github Copilot

Votre code actuel semble déjà accomplir la plupart des tâches que vous avez décrites. Cependant, il semble y avoir une erreur dans la façon dont vous renvoyez la réponse à la requête initiale. Vous devriez utiliser un canal de réponse différent pour éviter les conflits avec les requêtes de progression. Voici comment vous pourriez modifier votre code pour accomplir cela :

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
                    progressResponseChan <- *tasks[id]
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
```

Dans ce code, j'ai ajouté un canal responseChan pour renvoyer la réponse à la requête initiale. Cela permet d'éviter les conflits avec les requêtes de progression qui sont envoyées sur le canal requestChan.

### Erreur fatal error: all goroutines are asleep - deadlock!

# Reponse github copilot

L'erreur de deadlock se produit lorsque toutes les goroutines sont bloquées ou en attente d'une ressource qui ne peut être libérée. Dans votre cas, cela se produit probablement parce que vous essayez d'envoyer une valeur à progressResponseChan mais il n'y a pas de goroutine en attente pour recevoir cette valeur.

Pour résoudre ce problème, vous pouvez utiliser une déclaration select avec une clause default pour rendre l'opération d'envoi non bloquante. Si le canal est plein ou s'il n'y a pas de récepteur, le default sera exécuté et aucune valeur ne sera envoyée. Voici comment vous pouvez modifier votre code à la ligne 67 :

```go
select {
case progressResponseChan <- *tasks[id]:
default:
    fmt.Println("Aucun récepteur pour progressResponseChan")
}
```
