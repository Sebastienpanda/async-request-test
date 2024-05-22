package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	id          uint32
	progress    *sync.Mutex
	progressVal uint8
	progressTx  chan SchedulerMessage
}

func NewTask(id uint32, progressTx chan SchedulerMessage) *Task {
	return &Task{
		id:         id,
		progress:   &sync.Mutex{},
		progressTx: progressTx,
	}
}

func (t *Task) start() {
	progress := t.progress
	taskID := t.id
	progressTx := t.progressTx

	go func() {
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := 0; i <= 100; i++ {
			time.Sleep(time.Duration(rng.Intn(90)+10) * time.Millisecond)
			progress.Lock()
			t.progressVal = uint8(i)
			progress.Unlock()
			progressTx <- SchedulerMessage{Type: ProgressUpdate, TaskID: taskID, Progress: uint8(i)}
		}
		progressTx <- SchedulerMessage{Type: TaskFinished, TaskID: taskID}
	}()
}

type SchedulerMessageType int

const (
	NewTasks SchedulerMessageType = iota
	TaskFinished
	GetProgress
	ProgressUpdate
	Thanks
)

type SchedulerMessage struct {
	Type     SchedulerMessageType
	TaskID   uint32
	Progress uint8
	Response chan interface{}
}

func startTask(id uint32, tasksIDs *[]uint32, schedulerTx chan SchedulerMessage) {
	*tasksIDs = append(*tasksIDs, id)
	task := NewTask(id, schedulerTx)
	task.start()
}

func main() {
	schedulerTx := make(chan SchedulerMessage)
	schedulerRx := schedulerTx

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		actualID := uint32(0)
		waitingTasks := []uint32{}
		tasksIDs := []uint32{}
		progressMap := make(map[uint32]uint8)

		for {
			msg := <-schedulerRx
			switch msg.Type {
			case NewTasks:
				if len(tasksIDs) < 10 {
					startTask(actualID, &tasksIDs, schedulerTx)
					msg.Response <- actualID
				} else {
					fmt.Printf("--- Maximum number of threads reached. Task %d will be delayed.\n", actualID)
					waitingTasks = append(waitingTasks, actualID)
					msg.Response <- actualID
				}
				actualID++
			case TaskFinished:
				for i, id := range tasksIDs {
					if id == msg.TaskID {
						tasksIDs = append(tasksIDs[:i], tasksIDs[i+1:]...)
						break
					}
				}
				fmt.Printf("--- Task %d finished.\n", msg.TaskID)
				if len(waitingTasks) > 0 && len(tasksIDs) < 10 {
					waitingID := waitingTasks[0]
					waitingTasks = waitingTasks[1:]
					startTask(waitingID, &tasksIDs, schedulerTx)
					fmt.Printf("--- Waiting task %d started.\n", waitingID)
				}
			case GetProgress:
				progress, exists := progressMap[msg.TaskID]
				if exists {
					if progress == 100 {
						delete(progressMap, msg.TaskID)
					}
					msg.Response <- int8(progress)
				} else if contains(waitingTasks, msg.TaskID) {
					msg.Response <- int8(-1)
				} else {
					msg.Response <- int8(-2)
				}
			case ProgressUpdate:
				progressMap[msg.TaskID] = msg.Progress
			case Thanks:
				fmt.Println("--- All tasks finished.")
				fmt.Println("--- Merci au scheduler !")
				return
			}
		}
	}()

	ids := []uint32{}

	for i := 1; i <= 15; i++ {
		NewTasksTx := make(chan interface{})
		schedulerTx <- SchedulerMessage{Type: NewTasks, Response: NewTasksTx}
		NewTasksID := (<-NewTasksTx).(uint32)
		ids = append(ids, NewTasksID)
		fmt.Printf("New task created with ID: %d\n", NewTasksID)
	}

	fmt.Println(ids)

	for {
		time.Sleep(1 * time.Second)

		completedIDs := []uint32{}

		for _, id := range ids {
			progressTx := make(chan interface{})
			schedulerTx <- SchedulerMessage{Type: GetProgress, TaskID: id, Response: progressTx}
			progress := (<-progressTx).(int8)
			if progress >= 0 {
				fmt.Printf("Task %d progress: %d%% / ", id, progress)
			}

			if progress == -2 {
				completedIDs = append(completedIDs, id)
			}
		}
		fmt.Println()

		ids = filter(ids, func(x uint32) bool {
			return !contains(completedIDs, x)
		})

		if len(ids) == 0 {
			break
		}
	}

	fmt.Println("End of code !")

	schedulerTx <- SchedulerMessage{Type: Thanks}
	wg.Wait()
}

func contains(slice []uint32, item uint32) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func filter(slice []uint32, test func(uint32) bool) (ret []uint32) {
	for _, v := range slice {
		if test(v) {
			ret = append(ret, v)
		}
	}
	return
}
