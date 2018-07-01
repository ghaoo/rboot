package rboot

import (
	"net"
	"sync"
	"net/http"
	"fmt"
	"encoding/json"
	"log"
)

type httpRule struct {
	mux *http.ServeMux

	memoryRead func(key string) []byte
	memorySave func(key string, value []byte)

	listener net.Listener
	outCh    chan Message

	mu    sync.Mutex
	inbox []Message
}

func (r *httpRule) httpPop(w http.ResponseWriter, req *http.Request) {
	r.mu.Lock()
	defer r.mu.Unlock()
	defer req.Body.Close()

	var msg Message

	if len(r.inbox) > 1 {
		msg, r.inbox = r.inbox[0], r.inbox[1:]
	} else if len(r.inbox) == 1 {
		msg = r.inbox[0]
		r.inbox = []Message{}
	} else if len(r.inbox) == 0 {
		fmt.Fprint(w, "{}")
		return
	}

	if err := json.NewEncoder(w).Encode(&msg); err != nil {
		log.Fatal(err)
	}
}