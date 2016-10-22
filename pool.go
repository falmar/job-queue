// Copyright 2016 David Lavieri.  All rights reserved.
// Use of this source code is governed by a MIT License
// License that can be found in the LICENSE file.

package main

import "log"

// WorkerPool or scheduler (?)
type WorkerPool struct {
	MaxWorkers int64
	// worker.go
	Workers []*Worker
	Pool    chan *Worker
	Work    bool
}

// Start the workers
func (wp *WorkerPool) Start() {
	var f int64
	// allocate workers slice
	wp.Workers = []*Worker{}

	// create pool channel
	wp.Pool = make(chan *Worker, wp.MaxWorkers)

	for f = 0; f < wp.MaxWorkers; f++ {
		// create a worker
		worker := &Worker{ID: f, Job: make(chan int64), WorkerPool: wp.Pool, QuitChan: make(chan bool)}

		// add worker to slice
		wp.Workers = append(wp.Workers, worker)

		// start worker
		worker.Start()
	}
}

// Stop define if it should stop workers
func (wp *WorkerPool) Stop(stopWorkers chan<- bool) {
	log.Println("Make workers stop")

	var i int64
	var t = wp.MaxWorkers

	// create channel to get notification when worker have stop
	ws := make(chan bool)

	// send signal to stop
	for _, w := range wp.Workers {
		go w.Stop(ws)
	}

	// wait for their response
	for i = 0; i < t; i++ {
		<-ws
	}

	// close channel
	close(ws)

	// notify main.go all workers have stop
	stopWorkers <- true
}

// Do will add a job the workers
func (wp *WorkerPool) Do(job int64, result chan int64) {
	// check if can handle jobs
	if wp.Work {
		// take an avaiable worker from pool
		worker := <-wp.Pool

		// give it result channel to notify who ever await for it
		worker.Result = result

		// give it the job
		worker.Job <- job

		return
	}

	// else close result chanel
	close(result)
}
