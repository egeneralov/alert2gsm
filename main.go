package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"time"
)

const (
	envAccountSID = "TWILIO_ACCOUNT_SID"
	envAuthToken  = "TWILIO_AUTH_TOKEN"
)

var (
	accountSID string
	authToken  string

	urlPost        string
	urlBoilerplate = "https://api.twilio.com/2010-04-01/Accounts/%v/Calls.json"

	twilio  Twilio
	storage = make(map[string]string)

	twilioRun = false

	regexpWebhookTwilio *regexp.Regexp

	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func RandStringRunes(n int) string {
	log.Debugf("RandStringRunes(%d)", n)
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	var (
		err error
	)

	accountSID = os.Getenv(envAccountSID)
	authToken = os.Getenv(envAuthToken)

	if accountSID == "" || authToken == "" {
		log.Fatalf("Environments %v or %v is empty\n", envAccountSID, envAuthToken)
	}
	urlPost = fmt.Sprintf(urlBoilerplate, accountSID)

	regexpWebhookTwilio, err = regexp.Compile(`\/webhook\/twilio\/(.*)\.xml`)

	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())

	server := HTTPServer{
		Configuration: HTTPServerConfiguration{},
		Twilio: Twilio{
			AccountSID: accountSID,
			AuthToken:  authToken,
			// From:       ,
		},
	}

	jsonFile, err := os.Open("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(byteValue, &server.Configuration)
	if err != nil {
		log.Fatal(err)
	}

	server.Twilio.From = server.Configuration.Webhooks.Twilio.From

	server.Start()

	select {}
}
