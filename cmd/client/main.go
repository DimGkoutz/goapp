package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

// Counter represents the data received from the websocket
type Counter struct {
	Iteration int    `json:"iteration"`
	Value     string `json:"value"`
}

func main() {
	// Parse command-line arguments
	var numConnections int
	flag.IntVar(&numConnections, "n", 1, "Number of parallel connections")
	flag.Parse()

	if numConnections <= 0 {
		log.Fatalf("Number of connections must be greater than 0")
	}

	// Set up signal handling
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Track all active connections
	var connections []*websocket.Conn

	// Track when all connections have closed
	var wg sync.WaitGroup
	connectionsClosed := make(chan struct{})

	// Start connections
	for i := 0; i < numConnections; i++ {
		wg.Add(1)
		conn := handleConnection(i, &wg)
		if conn != nil {
			connections = append(connections, conn)
		} else {
			wg.Done() // Connection failed, decrement the wait group
		}
	}

	// Start a goroutine to notify when all connections close
	go func() {
		wg.Wait()
		close(connectionsClosed)
	}()

	// Wait for interrupt signal or all connections to close
	select {
	case <-interrupt:
		fmt.Println("Interrupted, Client Stopped")

		// Send close messages to all connections
		for _, conn := range connections {
			if conn != nil {
				conn.WriteMessage(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
				)
				conn.Close()
			}
		}

		// Small delay to allow close messages to be sent
		time.Sleep(50 * time.Millisecond)

	case <-connectionsClosed:
		fmt.Println("All connections closed, Client Stopped")
	}
}

func handleConnection(id int, wg *sync.WaitGroup) *websocket.Conn {
	defer wg.Done()

	// Connect to the WebSocket server
	url := "ws://localhost:8080/goapp/ws"
	dialer := &websocket.Dialer{
		HandshakeTimeout: 45 * time.Second,
	}

	// Set the Origin header to match the allowed origin
	header := http.Header{}
	header.Add("Origin", "http://localhost:8080")

	c, _, err := dialer.Dial(url, header)
	if err != nil {
		log.Printf("[conn #%d] error connecting: %v", id, err)
		return nil
	}

	// Start a goroutine to read messages
	go func() {
		defer func() {
			c.Close()
			wg.Done()
		}()

		wg.Add(1)

		for {
			var counter Counter
			err := c.ReadJSON(&counter)
			if err != nil {
				return
			}

			// Display the received data
			fmt.Printf("[conn #%d] iteration: %d, value: %s\n",
				id, counter.Iteration, counter.Value)
		}
	}()

	return c
}
