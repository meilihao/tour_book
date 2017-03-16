package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type PayloadCollection struct {
	WindowsVersion string    `json:"version"`
	Token          string    `json:"token"`
	Payloads       []Payload `json:"data"`
}

type Payload struct {
	Num int `json:"num"`
}

// Job represents the job to be run
type Job struct {
	Payload Payload
}

// Worker represents the worker that executes the job
type Worker struct {
	id          int
	workerPool  chan chan Job
	jobChannel  chan Job
	quitChannel chan bool
}

func NewWorker(id int, workerPool chan chan Job) Worker {
	return Worker{
		id:          id,
		workerPool:  workerPool,
		jobChannel:  make(chan Job),
		quitChannel: make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.workerPool <- w.jobChannel

			select {
			case job := <-w.jobChannel:
				// do job
				fmt.Println(w.id, job)

			case <-w.quitChannel:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quitChannel <- true
	}()
}

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	workerPool chan chan Job
	maxWorkers int
	jobQueue   chan Job
}

func NewDispatcher(maxWorkers int, jobQueue chan Job) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		workerPool: pool,
		maxWorkers: maxWorkers,
		jobQueue:   jobQueue,
	}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(i+1, d.workerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				workJobChannel := <-d.workerPool

				// dispatch the job to the worker job channel
				workJobChannel <- job
			}(job)
		}
	}
}

var (
	// A buffered channel that we can send work requests on.
	MaxQueue          = 100
	JobQueue chan Job = make(chan Job, MaxQueue)
)

// http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html
// http://studygolang.com/articles/5846
// refactored : https://gist.github.com/harlow/dbcd639cf8d396a2ab73
//time for i in {1..4096}; do curl -X POST -H "Content-Type: application/json" localhost:8080/work -d '{"vsersion":"a","token":"b","data":[{"num":1},{"num":2}]}'; done
// 不知为什么,修改MaxWorker(8|100),发现测试用例的耗时都是很接近
func main() {
	MaxWorker := 100
	dispatcher := NewDispatcher(MaxWorker, JobQueue)
	dispatcher.Run()

	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		requestHandler(w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	var content = &PayloadCollection{}
	err := json.NewDecoder(io.LimitReader(r.Body, 128)).Decode(&content)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	for _, payload := range content.Payloads {

		// let's create a job with the payload
		job := Job{Payload: payload}

		// Push the job onto the queue.
		JobQueue <- job
	}

	w.WriteHeader(http.StatusOK)
}
