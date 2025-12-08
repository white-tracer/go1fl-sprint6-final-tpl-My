package server

import (
	"io"
	"log"
	"net"
	"net/http"
	"time"

	handlers "github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
)

type Server struct {
	Addr         net.Addr
	Handler      http.Handler
	ErrorLog     io.Writer
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func NewServer(logger *log.Logger) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.MainHandler)
	mux.HandleFunc("/upload", handlers.UploadHandler)

	defaultAddrString := ":8080"

	defaultAddNet, err := net.ResolveTCPAddr("tcp", defaultAddrString)
	if err != nil {
		logger.Printf("FATAL: Не удалось разрешить адрес '%s' для сервера: %v", defaultAddrString, err)
	}

	server := &Server{
		Addr:         defaultAddNet,
		Handler:      mux,
		ErrorLog:     logger.Writer(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return server

}
