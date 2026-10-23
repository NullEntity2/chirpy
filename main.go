package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/NullEntity2/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	JWTSecret      string
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error opening db connection: %s", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		dbQueries: dbQueries,
		platform:  platform,
		JWTSecret: jwtSecret,
	}
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	serveMux.HandleFunc("GET /api/healthz", healthzFunc)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.hitsFunc)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetFunc)
	serveMux.HandleFunc("POST /api/users", apiCfg.postUsersFunc)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.postChirpsFunc)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.getChirpsFunc)
	serveMux.HandleFunc("GET /api/chirps/{id}", apiCfg.getChirpFunc)
	serveMux.HandleFunc("POST /api/login", apiCfg.postLoginFunc)
	serveMux.HandleFunc("POST /api/refresh", apiCfg.postRefreshFunc)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.postRevokeFunc)
	serveMux.HandleFunc("PUT /api/users", apiCfg.putUsersFunc)
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirpFunc)
	server := http.Server{Handler: serveMux, Addr: ":8081"}
	log.Fatal(server.ListenAndServe())
}
