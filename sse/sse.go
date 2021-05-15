package sse

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/soumitra003/goframework/config"
)

// Broker generic message broker
type Broker struct {
	clients        map[chan []byte]bool
	newClients     chan chan []byte
	defunctClients chan chan []byte
	messages       chan []byte
	stopSignal     chan bool
}

// Start starts the broker
func (b *Broker) Start() {
	go func() {
		for {
			select {
			case s := <-b.newClients:
				b.addNewClient(s)
			case s := <-b.defunctClients:
				b.removeClient(s)
			case msg := <-b.messages:
				b.broadcastMessage(msg)
			case _ = <-b.stopSignal:
				break
			}
		}
	}()
}

// Stop stop the broker
func (b *Broker) Stop() {
	b.stopSignal <- true
}

func (b *Broker) addNewClient(c chan []byte) {
	b.clients[c] = true
}

func (b *Broker) removeClient(c chan []byte) {
	delete(b.clients, c)
	close(c)
}

func (b *Broker) broadcastMessage(msg []byte) {
	for s := range b.clients {
		s <- msg
	}
}

type ModuleSSE struct {
	config *config.Config
	broker *Broker
}

//New creates module instance
func New(config config.Config) *ModuleSSE {
	md := &ModuleSSE{config: &config}
	return md
}

// Init initializes auth module
func (h *ModuleSSE) Init(ctx context.Context, config config.Config) {
	h.broker = &Broker{
		make(map[chan []byte]bool),
		make(chan (chan []byte)),
		make(chan (chan []byte)),
		make(chan []byte),
		make(chan bool),
	}

	h.broker.Start()

	go func() {
		for i := 0; ; i++ {

			// Create a little message to send to clients,
			// including the current time.
			h.broker.messages <- []byte(fmt.Sprintf("%d - the time is %v", i, time.Now()))

			// Print a nice log message and sleep for 5s.
			log.Printf("Sent message %d ", i)
			time.Sleep(5e9)

		}
	}()
}
