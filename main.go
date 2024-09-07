package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/ITzJustMe13/Chirpy/internal/database"
	"github.com/joho/godotenv"
)


func main(){

	godotenv.Load("secret.env")
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == ""{
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY environment variable is not set")
	}

	db, err := database.NewDB("database.json")
	if err != nil{
		log.Fatal(err)
	}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if dbg != nil && *dbg {
		err := db.ResetDB()
		if err != nil {
			log.Fatal(err)
		}
	}

	
	apicfg := apiConfig{
		fileserverHits: 0,
		DB: db,
		jwtSecret: jwtSecret,
		polkaKey:       polkaKey,
	}


	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", apicfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	serverMux.HandleFunc("GET /api/healthz", handlerfunc)
	serverMux.HandleFunc("GET /admin/metrics", apicfg.handlerMetrics)
	serverMux.HandleFunc("GET /api/reset", apicfg.handlerReset)
	serverMux.HandleFunc("POST /api/polka/webhooks", apicfg.handlerWebhook)

	serverMux.HandleFunc("DELETE /api/chirps/{chirpID}", apicfg.handlerChirpsDelete)
	serverMux.HandleFunc("POST /api/chirps", apicfg.handlerChirpsCreate)
	serverMux.HandleFunc("GET /api/chirps", apicfg.handlerChirpsRetrieve)
	serverMux.HandleFunc("GET /api/chirps/{chirpID}", apicfg.handlerChirpsGet)

	serverMux.HandleFunc("POST /api/users", apicfg.handlerUsersCreate)
	serverMux.HandleFunc("PUT /api/users", apicfg.handlerUsersUpdate)

	serverMux.HandleFunc("POST /api/login", apicfg.handlerLogin)
	serverMux.HandleFunc("POST /api/revoke", apicfg.handlerRevoke)
	serverMux.HandleFunc("POST /api/refresh", apicfg.handlerRefresh)


	server := &http.Server{
		Addr: ":8080",
		Handler: serverMux,
	}

	server.ListenAndServe()
}

func handlerfunc(writer http.ResponseWriter, request *http.Request){
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(http.StatusText(http.StatusOK)))
}


