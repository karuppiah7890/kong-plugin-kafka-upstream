package main

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/karuppiah7890/go-pdk"
	"log"
)

type PluginConfig struct {
	Producer sarama.SyncProducer
}

var Version = "0.1.0"

func New() interface{} {
	brokerList := []string{"localhost:9092"}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}
	return &PluginConfig{Producer: producer}
}

func (jsonEncoder JSONEncoder) Encode() ([]byte, error) {
	jsonData, err := json.Marshal(map[string]interface{}{"body": jsonEncoder})
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

type JSONEncoder string

func (jsonEncoder JSONEncoder) Length() int {
	data, err := jsonEncoder.Encode()
	if err != nil {
		return 0
	}
	return len(data)
}

func (conf *PluginConfig) Access(kong *pdk.PDK) {
	body, err := kong.Request.GetRawBody()
	if err != nil {
		sendError(kong, err)
		return
	}

	_, _, err = conf.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: "test",
		Value: JSONEncoder(body),
	})
	if err != nil {
		sendError(kong, err)
		return
	}

	err = kong.Response.Exit(200, map[string]interface{}{"message": "message sent"}, nil)
	if err != nil {
		_ = kong.Log.Err(err.Error())
	}
}

func sendError(kong *pdk.PDK, err error) {
	exitErr := kong.Response.Exit(200, map[string]interface{}{
		"message": "message not sent :/",
		"error":   err.Error(),
	}, nil)
	if exitErr != nil {
		_ = kong.Log.Err(err.Error())
	}
}
