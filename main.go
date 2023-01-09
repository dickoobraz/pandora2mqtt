package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"

	"github.com/eclipse/paho.mqtt.golang"
)

const (
	pandoraAPIURL      = "https://api.pandora-security.com"
	pandoraEventsPath  = "/v1/events/subscribe"
	pandoraTokenPath   = "/oauth2/token"
	pandoraAuthHeader  = "Authorization"
	pandoraBearer      = "Bearer "
	pandoraContentType = "Content-Type"
	pandoraJSON        = "application/json"

	mqttServer   = "tcp://localhost:1883"
	mqttClientID = "pandora-mqtt-client"

	envMqttUsername = "MQTT_USERNAME"
	envMqttPassword = "MQTT_PASSWORD"
	envPandoraEmail = "PANDORA_EMAIL"
	envPandoraPass  = "PANDORA_PASSWORD"
)

type pandoraEvent struct {
	EventType string          `json:"eventType"`
	Timestamp int64           `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

func main() {
	// Set up OAuth 2.0 client
	oauthConfig := &oauth2.Config{
		// TODO: Fill in OAuth 2.0 configuration
	}
	// TODO: Use environment variables for Pandora login and password
	token, err := getPandoraToken(oauthConfig)
	if err != nil {
		log.Fatal(err)
	}
	oauthClient := oauthConfig.Client(context.Background(), token)

	// Set up MQTT client
	opts := mqtt.NewClientOptions().AddBroker(mqttServer).SetClientID(mqttClientID)
	// Use environment variables for MQTT login and password if set
	if username, ok := os.LookupEnv(envMqttUsername); ok {
		opts.SetUsername(username)
	}
	if password, ok := os.LookupEnv(envMqttPassword); ok {
		opts.SetPassword(password)
	}
	mqttClient := mqtt.NewClient(opts)

	// Connect to MQTT broker
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// Subscribe to Pandora events
	eventsURL := pandoraAPIURL + pandoraEventsPath
	for {
		// Set up HTTP request
		req, err := http.NewRequest("POST", eventsURL, nil)
		if err != nil {
			log.Println(err)
