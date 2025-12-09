package main

// http://localhost:8080/upload - пример адреса

import (
	"bytes"
	"log"
	"net/http"

	server "github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	// создаем логгер
	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "logger: ", log.Lshortfile)
	)

	//создаем сервер

	serverConfig := server.NewServer(logger)

	httpServer := &http.Server{
		Addr:         serverConfig.Addr.String(),
		Handler:      serverConfig.Handler,
		ErrorLog:     log.New(serverConfig.ErrorLog, "HTTP_SERVER_ERROR: ", log.LstdFlags|log.Lshortfile),
		ReadTimeout:  serverConfig.ReadTimeout,
		WriteTimeout: serverConfig.WriteTimeout,
		IdleTimeout:  serverConfig.IdleTimeout,
	}

	// запускаем сервер
	if err := httpServer.ListenAndServe(); err != nil {
		logger.Printf("FATAL: Ошибка при запуске HTTP-сервера: %v\n", err)
	}

}
