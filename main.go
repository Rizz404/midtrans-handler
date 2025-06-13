package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/Rizz404/midtrans-handler/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"google.golang.org/api/option"
)

type apiConfig struct {
	Firestore         *firestore.Client
	MidtransCore      *coreapi.Client
	MidtransServerKey string
}

func main() {
	godotenv.Load()

	addr := os.Getenv("ADDR")
	if addr == "" {
		log.Fatal("Address is not found in env")
	}

	// * Database
	ctx := context.Background()

	// ---- MULAI PERUBAHAN ----
	// Muat kredensial Firebase dari environment variables
	privateKey := os.Getenv("FIREBASE_PRIVATE_KEY")
	if privateKey == "" {
		log.Fatal("FIREBASE_PRIVATE_KEY is not found in env")
	}
	// Ganti kembali \\n menjadi \n
	privateKey = strings.ReplaceAll(privateKey, "\\n", "\n")

	// Buat struktur kredensial dalam bentuk map
	creds := map[string]string{
		"type":                        os.Getenv("FIREBASE_TYPE"),
		"project_id":                  os.Getenv("FIREBASE_PROJECT_ID"),
		"private_key_id":              os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
		"private_key":                 privateKey,
		"client_email":                os.Getenv("FIREBASE_CLIENT_EMAIL"),
		"client_id":                   os.Getenv("FIREBASE_CLIENT_ID"),
		"auth_uri":                    os.Getenv("FIREBASE_AUTH_URI"),
		"token_uri":                   os.Getenv("FIREBASE_TOKEN_URI"),
		"auth_provider_x509_cert_url": os.Getenv("FIREBASE_AUTH_PROVIDER_X509_CERT_URL"),
		"client_x509_cert_url":        os.Getenv("FIREBASE_CLIENT_X509_CERT_URL"),
	}

	// Ubah map menjadi JSON dalam bentuk byte slice
	credsJSON, err := json.Marshal(creds)
	if err != nil {
		log.Fatalf("failed to marshal credentials to JSON: %v", err)
	}

	// Buat credential option menggunakan JSON yang sudah kita buat
	opt := option.WithCredentialsJSON(credsJSON)

	// Inisialisasi aplikasi Firebase dengan option tersebut
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app with manual credentials: %v\n", err)
	}
	// ---- SELESAI PERUBAHAN ----

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error getting Firestore client: %v\n", err)
	}
	defer firestoreClient.Close()

	// * MidtransClient
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		log.Fatal("MIDTRANS_SERVER_KEY is not found in env")
	}

	midtransClient := coreapi.Client{}
	midtransClient.New(serverKey, midtrans.Sandbox)

	apiCfg := apiConfig{
		Firestore:         firestoreClient,
		MidtransCore:      &midtransClient,
		MidtransServerKey: serverKey,
	}

	router := chi.NewRouter()

	// * Middleware
	router.Use(middleware.RequestLoggerMiddleware)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	// * Routes
	v1Router.Get("/health", handlerHealth)
	v1Router.Mount("/webhooks", webhookRoutes(&apiCfg))
	v1Router.Mount("/payment-methods", paymentMethodRoutes(&apiCfg))
	v1Router.Mount("/orders", OrderRoutes(&apiCfg))

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Printf("Server running on http://localhost%s", addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
