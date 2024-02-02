package main

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

// Constants for different job statuses.
const (
	queued    = "queued"
	running   = "running"
	completed = "completed"
	failed    = "failed"
)

var (
	jobQueue     = make(chan *Job, 100)
	jobStatusMap = make(map[string]*jobStatus)
	mutex        = &sync.Mutex{}
	wg           sync.WaitGroup
)

// JobType represents different types of jobs that can be handled.
type JobType int

// Enumeration of different JobTypes.
const (
	WriteTimeToFile JobType = iota
	AnotherJobType
	// Future job types should be added here.
)

// Job defines the structure of a job including its type, payload, and callback.
type Job struct {
	ID       string
	JobType  JobType
	Payload  interface{}
	Callback func(result interface{}, err error)
}

// jobStatus represents the current status of a job along with its result or error.
type jobStatus struct {
	Status string
	Result interface{}
	Error  error
}

// worker is a goroutine that processes jobs from the jobQueue.
func worker() {
	for job := range jobQueue {
		processJob(job)
	}
}

// processJob handles the execution and updating of the job status.
func processJob(job *Job) {
	status := &jobStatus{Status: running}
	mutex.Lock()
	jobStatusMap[job.ID] = status
	mutex.Unlock()

	result, err := executeJob(job)
	updateJobStatus(job.ID, result, err)

	job.Callback(result, err)
	wg.Done()
}

// executeJob executes the given job based on its JobType.
func executeJob(job *Job) (interface{}, error) {
	switch job.JobType {
	case WriteTimeToFile:
		return writeTimeToFile(job.Payload)
	case AnotherJobType:
		return anotherJobFunction(job.Payload)
	default:
		return nil, errors.New("unknown job type")
	}
}

// SubmitJob adds a job to the jobQueue and tracks its status.
func SubmitJob(job *Job) {
	wg.Add(1)
	mutex.Lock()
	jobStatusMap[job.ID] = &jobStatus{Status: queued}
	mutex.Unlock()
	jobQueue <- job
}

// GetJobStatus returns the status of a job by its ID.
func GetJobStatus(jobID string) *jobStatus {
	mutex.Lock()
	defer mutex.Unlock()
	return jobStatusMap[jobID]
}

// GetAllJobsStatus returns the status of all jobs.
func GetAllJobsStatus() map[string]*jobStatus {
	mutex.Lock()
	defer mutex.Unlock()
	return jobStatusMap
}

// InitWorkers initializes a specified number of worker goroutines.
func InitWorkers(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go worker()
	}
}

// updateJobStatus updates the status of a job in the jobStatusMap.
func updateJobStatus(jobID string, result interface{}, err error) {
	mutex.Lock()
	defer mutex.Unlock()

	status := jobStatusMap[jobID]
	if err != nil {
		status.Status = failed
		status.Error = err
	} else {
		status.Status = completed
		status.Result = result
	}
}

// writeTimeToFile handles the specific logic for writing the time to a file.
// This is a function to test the jobs system and not meant for production use.
func writeTimeToFile(payload interface{}) (interface{}, error) {
	waitTime, ok := payload.(time.Duration)
	if !ok {
		return nil, errors.New("invalid payload")
	}

	// Wait for the specified time
	fmt.Println("Waiting for", waitTime)
	time.Sleep(waitTime)

	// Get the current time
	currentTime := time.Now().Format(time.RFC3339)

	// Create a file in /tmp with the current time
	filePath := fmt.Sprintf("/tmp/job_%s.txt", currentTime)
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Write the current time to the file
	_, err = file.WriteString(currentTime)
	if err != nil {
		return nil, fmt.Errorf("failed to write to file: %v", err)
	}

	return "File created with time: " + currentTime, nil
}

// anotherJobFunction represents a placeholder for future job types.
func anotherJobFunction(payload interface{}) (interface{}, error) {
	// Code for another type of job
	// ...
	return nil, errors.New("not implemented")
}
