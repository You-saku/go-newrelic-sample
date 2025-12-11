package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func sampleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello golang !!")
}

func newRelicSampleHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello NewRelic !!")
}

func newRelicSampleHandler2(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello NewRelic2 !!")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	licenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY")

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("go-newrelic-sample"),
		newrelic.ConfigLicense(licenseKey),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		log.Println("New Relic Application creation failed:", err)
		return
	}

	txn := app.StartTransaction("Transaction Test")
	defer txn.End()

	http.HandleFunc("/", sampleHandler)
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/newrelic", newRelicSampleHandler))
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/newrelic2", newRelicSampleHandler2))
	http.ListenAndServe(":8080", nil)
}
