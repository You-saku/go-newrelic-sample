package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus" // output structured logs

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrlogrus" // newrelic logrus formatter
	"github.com/newrelic/go-agent/v3/newrelic"
)

func sampleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello golang !!")
}

func newRelicSampleHandler(w http.ResponseWriter, r *http.Request) {
	logrus.WithFields(logrus.Fields{
		"attribute": "sample1",
	}).Info("Hello NewRelic Info !!")
}

func newRelicSampleHandler2(w http.ResponseWriter, r *http.Request) {
	logrus.WithFields(logrus.Fields{
		"attribute": "sample2",
	}).Error("Hello NewRelic Error !!")
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

	nrlogrusFormatter := nrlogrus.NewFormatter(app, &logrus.JSONFormatter{})
	// logrusでログを出力する箇所がhandlerの関数なのでlogrus.New()は使わない
	logrus.SetLevel(logrus.DebugLevel)     // Set log level to Debug
	logrus.SetFormatter(nrlogrusFormatter) // Set New Relic logrus formatter

	http.HandleFunc("/", sampleHandler)
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/newrelic", newRelicSampleHandler))
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/newrelic2", newRelicSampleHandler2))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("HTTP server failed:", err)
	}
}
