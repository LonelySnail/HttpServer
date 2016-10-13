package server

import (
	"log"
	"net/http"
	"time"

	"../config"
)

var logger *log.Logger

func createServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Handle)
	var server = http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return &server
}

func StartServer(logger *log.Logger) {
	cfg := config.GlobalConfig()
	defer RecoverPanic()
	logger.Println("Start HTTP Server:", cfg.Listen)
	server := createServer(cfg.Listen)
	err := server.ListenAndServe()
	if err != nil {
		logger.Println("Cannot Start HTTP Server")
		log.Fatal(err)
	}
}

func RecoverPanic() {
	if err := recover(); err != nil {
		logger.Println(err)
	}
}
