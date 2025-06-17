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

	// * Validasi semua environment variable yang diperlukan
	requiredEnvVars := map[string]string{
		"FIREBASE_TYPE":                        os.Getenv("FIREBASE_TYPE"),
		"FIREBASE_PROJECT_ID":                  os.Getenv("FIREBASE_PROJECT_ID"),
		"FIREBASE_PRIVATE_KEY_ID":              os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
		"FIREBASE_PRIVATE_KEY":                 os.Getenv("FIREBASE_PRIVATE_KEY"),
		"FIREBASE_CLIENT_EMAIL":                os.Getenv("FIREBASE_CLIENT_EMAIL"),
		"FIREBASE_CLIENT_ID":                   os.Getenv("FIREBASE_CLIENT_ID"),
		"FIREBASE_AUTH_URI":                    os.Getenv("FIREBASE_AUTH_URI"),
		"FIREBASE_TOKEN_URI":                   os.Getenv("FIREBASE_TOKEN_URI"),
		"FIREBASE_AUTH_PROVIDER_X509_CERT_URL": os.Getenv("FIREBASE_AUTH_PROVIDER_X509_CERT_URL"),
		"FIREBASE_CLIENT_X509_CERT_URL":        os.Getenv("FIREBASE_CLIENT_X509_CERT_URL"),
	}

	// * Periksa apakah ada environment variable yang kosong
	for key, value := range requiredEnvVars {
		if value == "" {
			log.Fatalf("%s is not found in env", key)
		}
	}

	projectID := requiredEnvVars["FIREBASE_PROJECT_ID"]

	// * Proses private key
	privateKey := requiredEnvVars["FIREBASE_PRIVATE_KEY"]
	// * Ganti kembali \\n menjadi \n
	privateKey = strings.ReplaceAll(privateKey, "\\n", "\n")

	// * Buat struktur kredensial dalam bentuk map
	creds := map[string]string{
		"type":                        requiredEnvVars["FIREBASE_TYPE"],
		"project_id":                  projectID,
		"private_key_id":              requiredEnvVars["FIREBASE_PRIVATE_KEY_ID"],
		"private_key":                 privateKey,
		"client_email":                requiredEnvVars["FIREBASE_CLIENT_EMAIL"],
		"client_id":                   requiredEnvVars["FIREBASE_CLIENT_ID"],
		"auth_uri":                    requiredEnvVars["FIREBASE_AUTH_URI"],
		"token_uri":                   requiredEnvVars["FIREBASE_TOKEN_URI"],
		"auth_provider_x509_cert_url": requiredEnvVars["FIREBASE_AUTH_PROVIDER_X509_CERT_URL"],
		"client_x509_cert_url":        requiredEnvVars["FIREBASE_CLIENT_X509_CERT_URL"],
	}

	// * Tambahkan universe_domain jika ada
	if universeDomain := os.Getenv("FIREBASE_UNIVERSE_DOMAIN"); universeDomain != "" {
		creds["universe_domain"] = universeDomain
	}

	// * Ubah map menjadi JSON dalam bentuk byte slice
	credsJSON, err := json.Marshal(creds)
	if err != nil {
		log.Fatalf("failed to marshal credentials to JSON: %v", err)
	}

	// * Buat credential option menggunakan JSON yang sudah kita buat
	opt := option.WithCredentialsJSON(credsJSON)

	// * Buat konfigurasi Firebase dengan project ID yang eksplisit
	config := &firebase.Config{
		ProjectID: projectID,
	}

	// * Inisialisasi aplikasi Firebase dengan config dan option
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("error initializing app with manual credentials: %v\n", err)
	}

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
