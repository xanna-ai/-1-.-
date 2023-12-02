package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	ID             int   `json:"id"`
	DaysRequired   int   `json:"days_required"`
	WorkersNeeded  int   `json:"workers_needed"`
	PreviousJobs   []int `json:"previous_jobs"`
	CompletedTask  bool  `json:"-"`
	CurrentWorkers int   `json:"-"`
}

type Project struct {
	Tasks []Task `json:"tasks"`
}

func CalculateDuration(project Project) int {
	maxResource := 10
	totalDuration := 0

	for i, task := range project.Tasks {
		taskDuration := 0
		resourceUsed := 0

		if i > 0 && !project.Tasks[i-1].CompletedTask {
			continue
		}

		if task.CurrentWorkers > maxResource {
			task.CurrentWorkers = maxResource
		}

		for day := 0; day < task.DaysRequired; day++ {
			for resourceUsed < maxResource && task.CurrentWorkers > 0 {
				if resourceUsed+1 <= maxResource {
					task.CurrentWorkers--
					resourceUsed++
				}

				taskDuration++
			}
		}

		if taskDuration > task.DaysRequired { 
			task.DaysRequired = taskDuration
		}

		totalDuration += task.DaysRequired 
		project.Tasks[i].CompletedTask = true
	}

	return totalDuration
}

func main() {
	rand.Seed(time.Now().UnixNano())

	taskData, err := ioutil.ReadFile("tasks.json")
	if err != nil {
		fmt.Println("Error reading tasks.json:", err)
		return
	}

	var project Project
	err = json.Unmarshal(taskData, &project)
	if err != nil {
		fmt.Println("Error parsing tasks.json:", err)
		return
	}

	for i := range project.Tasks {
		project.Tasks[i].CurrentWorkers = rand.Intn(6) + 5 
	}

	projectJSON, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling project:", err)
		return
	}

	err = saveToFile("project.json", projectJSON)
	if err != nil {
		fmt.Println("Error saving project to file:", err)
		return
	}

	duration := CalculateDuration(project)
	fmt.Printf("Minimum Duration of the project: %d days\n", duration)

	const numSequences = 1000000
	const parallelism = 100

	var wg sync.WaitGroup
	wg.Add(parallelism)

	startTime := time.Now()
	sequenceChan := make(chan []int, parallelism)

	for i := 0; i < parallelism; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numSequences/parallelism; j++ {
				randomSequence := generateRandomSequence()
				sequenceChan <- randomSequence
			}
		}()
	}

	go func() {
		wg.Wait()
		close(sequenceChan)
	}()

	for sequence := range sequenceChan {
		project.Tasks[2].CurrentWorkers = len(sequence)
		_ = CalculateDuration(project)
	}

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)

	fmt.Printf("Parallel calculation of %d sequences took %s\n", numSequences, elapsed)
}

func saveToFile(filename string, data []byte) error {
	return nil
}

func generateRandomSequence() []int {
	numWorkers := rand.Intn(6) + 5
	var randomSequence []int

	for i := 0; i < numWorkers; i++ {
		randomSequence = append(randomSequence, i+1)
	}

	return randomSequence
}
