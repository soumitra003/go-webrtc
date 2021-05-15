package sse

import (
	"errors"
	"net/http"

	"github.com/soumitra003/go-webrtc/internal/render"
)

func (m *ModuleSSE) subscribeToEvents(writer http.ResponseWriter, request *http.Request) {
	f, ok := writer.(http.Flusher)
	if !ok {
		render.RenderError(writer, errors.New("streaming unsupported"))
		return
	}
	ch := make(chan []byte)
	m.broker.newClients <- ch

	// Listen to the closing of the http connection via the CloseNotifier
	notify := writer.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		// Remove this client from the map of attached clients
		// when `EventHandler` exits.
		m.broker.defunctClients <- ch
	}()

	// Set the headers related to event streaming.
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	// writer.Header().Set("Transfer-Encoding", "chunked")

	// Don't close the connection, instead loop endlessly.
	for {
		msg, open := <-ch
		if !open {
			break
		}

		// Write to the ResponseWriter, `w`.
		render.RenderStreamMessage(writer, msg)

		f.Flush()
	}
}
