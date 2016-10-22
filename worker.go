// Copyright 2016 David Lavieri.  All rights reserved.
// Use of this source code is governed by a MIT License
// License that can be found in the LICENSE file.

package main

import (
	"log"
	"time"
)

// Worker just self explainaroty
type Worker struct {
	ID         int64
	Result     chan int64
	Job        chan int64
	WorkerPool chan *Worker
	QuitChan   chan bool
}

// Start processing jobs
func (w *Worker) Start() {
	go func() {
		for {
			// add ourselves to the worker pool
			w.WorkerPool <- w

			select {
			// wait for jobs
			case job := <-w.Job:

				// process the job
				result := job * job

				// calculate ETA
				eta := time.Duration(random(10, 30))

				// log out
				log.Println("Worker:", w.ID, "Received:", job, "Result:", result, "ETA:", time.Second*eta)

				// simulate working time
				time.Sleep(time.Second * eta)

				// job done. Send the result back
				w.Result <- result

			// wait for quit signal to stop working
			case <-w.QuitChan:
				// quit working
				return
			}
		}
	}()
}

// Stop working
func (w *Worker) Stop(stop chan<- bool) {
	// notify it to stop
	w.QuitChan <- true

	// close worker channels
	close(w.QuitChan)
	close(w.Job)

	log.Printf("Worker %d stop\n", w.ID)

	// notify that worker have stop
	stop <- true
}
