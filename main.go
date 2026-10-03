package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"clasp/internal/aggregator"
	"clasp/internal/api"
	"clasp/internal/bangs"
	"clasp/internal/engine"
	"clasp/internal/engine/helpers"
	"clasp/internal/engine/sources/brave"
	"clasp/internal/engine/sources/ddg"
	"clasp/internal/engine/sources/startpage"
	"clasp/internal/limiter"
	"clasp/internal/suggest"
	"clasp/internal/wikipedia"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	client, err := helpers.NewBrowserClient()
	if err != nil {
		log.Fatalf("failed to create browser client: %v", err)
	}
	engines := []engine.Engine{
		brave.New(client),
		ddg.New(client),
		startpage.New(client),
	}

	agg := aggregator.New(engines)
	bangStore := bangs.New(context.Background())
	wikipediaClient := wikipedia.NewClient(client)
	suggestClient := suggest.NewClient(client)

	apiHandler := api.New(agg, bangStore, wikipediaClient, suggestClient)
	rl := limiter.New(limiter.LoadConfig())
	defer rl.Close()

	mux := http.NewServeMux()
	mux.Handle("/api/", rl.Middleware(apiHandler))
	mux.HandleFunc("/", serveStatic)

	addr := ":" + port
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("clasp listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("server stopped")
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	staticDir := "frontend/build"

	if r.URL.Path == "/" {
		http.ServeFile(w, r, staticDir+"/200.html")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/")
	filePath := filepath.Join(staticDir, filepath.Clean("/"+path))

	absStatic, _ := filepath.Abs(staticDir)
	absFile, _ := filepath.Abs(filePath)
	if !strings.HasPrefix(absFile, absStatic) {
		http.NotFound(w, r)
		return
	}

	if _, err := os.Stat(filePath); err == nil {
		http.ServeFile(w, r, filePath)
		return
	}

	http.ServeFile(w, r, staticDir+"/200.html")
}
