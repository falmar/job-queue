// Copyright 2016 David Lavieri.  All rights reserved.
// Use of this source code is governed by a MIT License
// License that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
)

func listen(pool *WorkerPool, stopListen <-chan bool, closeListen <-chan bool) {
	ln, err := net.Listen("tcp", ":8080")
	handle := true

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Listening at:", ln.Addr())

	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()

			if err != nil {
				log.Println(err)
				continue
			}

			// check if can handle new connections
			if handle {
				go handleConn(conn, pool)
				continue
			}

			// close connetions
			conn.Close()
		}

	}()

	for {
		select {
		case <-stopListen:
			handle = false
			log.Println("Stop handling new connections")
		case <-closeListen:
			log.Println("Stop Listening")
			return
		}
	}
}

func handleConn(conn net.Conn, pool *WorkerPool) {
	writer := bufio.NewWriter(conn)
	scanner := bufio.NewScanner(conn)

	writer.WriteString("Give a integer to double: ")
	writer.Flush()

	for scanner.Scan() {
		// conver string to int
		number, err := strconv.ParseInt(scanner.Text(), 10, 64)

		if err != nil {
			writer.WriteString(fmt.Sprintf("%s\n", err.Error()))
			writer.WriteString("Give a integer to double: ")
			writer.Flush()
			continue
		}

		// create result channel to receive from worker
		resultChan := make(chan int64)
		pool.Do(number, resultChan)

		writer.WriteString(fmt.Sprintf("Result (%d*%d): %d \n", number, number, <-resultChan))
		writer.Flush()

		conn.Close()
	}
}
