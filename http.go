package rboot

import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"net"
	"sync"
	"encoding/json"
)

type httpCall struct {
	mux *http.ServeMux

	memoryRead func(key string) []byte
	memorySave func(key string, value []byte)

	listener net.Listener
	inMessage    chan Message
	outMessage chan Message

	mu    sync.Mutex
}

func NewHttpCall(listener net.Listener) *httpCall {
	return &httpCall{
		mux:      http.NewServeMux(),
		listener: listener,
	}
}

func (hc *httpCall) Boot(bot *Robot) {
	hc.memoryRead = bot.MemoRead
	hc.memorySave = bot.MemoSave
	hc.inMessage = bot.Incoming()

	hc.mux.HandleFunc("/pop", hc.httpPop)
	hc.mux.HandleFunc("/send", hc.httpSend)
	hc.mux.HandleFunc("/memoryRead", hc.httpMemoryRead)
	hc.mux.HandleFunc("/memorySave", hc.httpMemorySave)
	srv := &http.Server{Handler: hc.mux}
	srv.SetKeepAlivesEnabled(false)
	go srv.Serve(hc.listener)
}

func (hc *httpCall) httpPop(w http.ResponseWriter, req *http.Request) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	defer req.Body.Close()

	for out := range hc.outMessage {
		if err := json.NewEncoder(w).Encode(&out); err != nil {
			log.Fatal(err)
		}
	}



}

func (hc *httpCall) httpSend(w http.ResponseWriter, req *http.Request) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	var msg Message
	if err := json.NewDecoder(req.Body).Decode(&msg); err != nil {
		panic(err)
	}
	defer req.Body.Close()

	go func(m Message) {
		hc.inMessage <- m
	}(msg)

	fmt.Fprintln(w, "OK")
}

func (hc *httpCall) httpMemoryRead(w http.ResponseWriter, req *http.Request) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	defer req.Body.Close()

	key := req.URL.Query().Get("key")

	fmt.Fprintf(w, "%s", hc.memoryRead(key))
}

func (hc *httpCall) httpMemorySave(w http.ResponseWriter, req *http.Request) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	defer req.Body.Close()

	key := req.URL.Query().Get("key")

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	hc.memorySave(key, b)
	fmt.Fprintln(w, "OK")
}