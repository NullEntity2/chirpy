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
	format := `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
	fmt.Fprintf(resp, format, cfg.fileserverHits.Load())
	resp.Header().Add("Content-Type", "text/html")
}

func (cfg *apiConfig) resetFunc(resp http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	fmt.Fprintf(resp, "Hits reset")
}

func main() {
	apiCfg := apiConfig{}
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	serveMux.HandleFunc("GET /api/healthz", healthzFunc)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.hitsFunc)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetFunc)
	serveMux.HandleFunc("POST /api/validate_chirp", validateChirp)
	server := http.Server{Handler: serveMux, Addr: ":8081"}
	log.Fatal(server.ListenAndServe())
}

func healthzFunc(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(200)
	resp.Write([]byte("OK"))
}
