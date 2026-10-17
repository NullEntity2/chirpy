package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(resp, req)
	})
}

func (cfg *apiConfig) hitsFunc(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "Hits: %d", cfg.fileserverHits.Load())
}

func (cfg *apiConfig) resetFunc(resp http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	fmt.Fprintf(resp, "Hits reset")
}

func main() {
	apiCfg := apiConfig{}
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	serveMux.HandleFunc("GET /healthz", healthzFunc)
	serveMux.HandleFunc("GET /metrics", apiCfg.hitsFunc)
	serveMux.HandleFunc("POST /reset", apiCfg.resetFunc)
	server := http.Server{Handler: serveMux, Addr: ":8081"}
	log.Fatal(server.ListenAndServe())
}

func healthzFunc(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(200)
	resp.Write([]byte("OK"))
}
