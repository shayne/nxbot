package nxbot

import (
	"bytes"
	"log"
	"net/http"
	"strings"
)

// OnMotionEvent sets the handler for motion events received over HTTP
func (b *NxBot) OnMotionEvent(handler func(string)) {
	b.motionEventHandler = handler
}

func (b *NxBot) handleHTTP(w http.ResponseWriter, r *http.Request) {
	if b.motionEventHandler == nil {
		return
	}
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	r.Body.Close()

	id := buf.String()
	id = strings.TrimSpace(id)

	b.motionEventHandler(id)
}

func (b *NxBot) startMotionInBackground() {
	http.HandleFunc("/", b.handleHTTP)
	go func() {
		http.ListenAndServe(b.settings.HTTPIPPort, nil)
	}()
}
