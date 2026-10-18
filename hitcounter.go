package main

import (
	"fmt"
	"net/http"
)

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
	if cfg.platform != "dev" {
		resp.WriteHeader(http.StatusForbidden)
		return
	}

	if err := cfg.dbQueries.DeleteAllUsers(req.Context()); err != nil {
		fmt.Printf("error truncating table: %s\n", err)
	}

	cfg.fileserverHits.Store(0)
	fmt.Fprintf(resp, "Hits and user db reset")
}
