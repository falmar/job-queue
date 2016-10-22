// Copyright 2016 David Lavieri.  All rights reserved.
// Use of this source code is governed by a MIT License
// License that can be found in the LICENSE file.

package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	stopWorkers := make(chan bool)
	stopListen := make(chan bool)
	closeListen := make(chan bool)

	// pool.go
	pool := &WorkerPool{MaxWorkers: 5, Work: true}
	pool.Start()

	go func() {
		<-quit

		go pool.Stop(stopWorkers)

		stopListen <- true
	}()

	// server.go
	go listen(pool, stopListen, closeListen)

	log.Printf("Workers Done: %t\n", <-stopWorkers)

	closeListen <- true
}

func random(min, max int64) int64 {
	return (rand.Int63n(max-min) + min) + 1
}
