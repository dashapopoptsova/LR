package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres password=postgres dbname=posts sslmode=disable"
	}

	repo, err := newRepository(dsn)
	if err != nil {
		log.Fatal(err)
	}

	svc := newService(repo)
	h := newHandler(svc)

	srv := &http.Server{
		Addr:    addr,
		Handler: h.routes(),
	}

	go func() {
		log.Println("server started on", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("shutting down...")
	srv.Shutdown(ctx)
}
