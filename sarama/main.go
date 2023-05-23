package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type Message struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MessageHandler interface {
	HandleMessage(msg *sarama.ConsumerMessage)
}

type ErrorHandler interface {
	HandleError(err error)
}

type SimpleErrorHandler struct{}

func (seh *SimpleErrorHandler) HandleError(err error) {
	log.Printf("Error consuming message: %v", err)
}

type KafkaConsumer struct {
	consumer sarama.ConsumerGroup
	topics   []string
}

func NewKafkaConsumer(brokers []string, topics []string, groupID string, config *sarama.Config) (*KafkaConsumer, error) {
	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{consumer: consumer, topics: topics}, nil
}

func (kc *KafkaConsumer) Consume(ctx context.Context, messageHandler MessageHandler, errorHandler ErrorHandler) error {
	handler := consumerGroupHandler{messageHandler: messageHandler, errorHandler: errorHandler}

	for {
		select {
		case <-ctx.Done():
			return kc.consumer.Close()
		default:
			err := kc.consumer.Consume(ctx, kc.topics, handler)
			if err != nil {
				errorHandler.HandleError(err)
			}
		}
	}
}

type consumerGroupHandler struct {
	messageHandler MessageHandler
	errorHandler   ErrorHandler
}

func (cgh consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (cgh consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (cgh consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		cgh.messageHandler.HandleMessage(msg)
		sess.MarkMessage(msg, "")
	}

	return nil
}

type ElasticsearchMessageHandler struct {
	client *elasticsearch.Client
	index  string
}

func NewElasticsearchMessageHandler(client *elasticsearch.Client, index string) *ElasticsearchMessageHandler {
	return &ElasticsearchMessageHandler{client: client, index: index}
}

func (emh *ElasticsearchMessageHandler) findDocumentIDByValue(ctx context.Context, value string) (string, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"value.keyword": value,
			},
		},
		"_source": false,
	}

	var b strings.Builder
	if err := json.NewEncoder(&b).Encode(query); err != nil {
		return "", fmt.Errorf("error encoding query: %v", err)
	}

	res, err := emh.client.Search(
		emh.client.Search.WithContext(ctx),
		emh.client.Search.WithIndex(emh.index),
		// emh.client.Search.WithBody(&b),
		emh.client.Search.WithSize(1),
		emh.client.Search.WithPretty(),
	)
	if err != nil {
		return "", fmt.Errorf("error executing search: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return "", fmt.Errorf("search response contains an error: %s", res.String())
	}

	var searchResults map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResults); err != nil {
		return "", fmt.Errorf("error decoding search results: %v", err)
	}

	hits := searchResults["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(hits) == 0 {
		return "", nil
	}

	return hits[0].(map[string]interface{})["_id"].(string), nil

}

func (emh *ElasticsearchMessageHandler) upsertMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	message := Message{Key: string(msg.Key), Value: string(msg.Value)}

	documentID, err := emh.findDocumentIDByValue(ctx, message.Value)
	if err != nil {
		return fmt.Errorf("error finding document ID by value: %v", err)
	}

	if documentID == "" {
		// Create new document if it does not exist
		document := map[string]interface{}{
			"key":   message.Key,
			"value": message.Value,
		}

		indexReq := esapi.IndexRequest{
			Index:   emh.index,
			Body:    esutil.NewJSONReader(document),
			Refresh: "true",
		}

		res, err := indexReq.Do(ctx, emh.client)
		if err != nil {
			return fmt.Errorf("error indexing document: %v", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			return fmt.Errorf("index response contains an error: %s", res.String())
		}
	} else {
		// Update existing document
		updateReq := esapi.UpdateRequest{
			Index:      emh.index,
			DocumentID: documentID,
			Body:       strings.NewReader(fmt.Sprintf(`{"doc": {"key": %q}}`, message.Key)),
			Refresh:    "true",
		}

		res, err := updateReq.Do(ctx, emh.client)
		if err != nil {
			return fmt.Errorf("error updating document: %v", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			return fmt.Errorf("update response contains an error: %s", res.String())
		}
	}

	log.Printf("Upserted message: %v\n", message)
	return nil
}

func (emh *ElasticsearchMessageHandler) HandleMessage(msg *sarama.ConsumerMessage) {
	if err := emh.upsertMessage(context.Background(), msg); err != nil {
		log.Printf("Error upserting message: %v", err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigterm
		cancel()
	}()

	brokers := []string{"b-4.dagmarkafkamskbeta.6zb6fu.c14.kafka.us-east-1.amazonaws.com:9092", "b-1.dagmarkafkamskbeta.6zb6fu.c14.kafka.us-east-1.amazonaws.com:9092", "b-6.dagmarkafkamskbeta.6zb6fu.c14.kafka.us-east-1.amazonaws.com:9092"}
	groupID := "example_group"
	topics := []string{"acom-offer-all-incr-mongo"}
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := NewKafkaConsumer(brokers, topics, groupID, config)
	if err != nil {
		log.Panicf("Error creating Kafka consumer: %v", err)
	}

	elasticsearchConfig := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	esClient, err := elasticsearch.NewClient(elasticsearchConfig)
	if err != nil {
		log.Panicf("Error creating Elasticsearch client: %v", err)
	}

	elasticsearchIndex := "your_index"

	messageHandler := NewElasticsearchMessageHandler(esClient, elasticsearchIndex)
	errorHandler := &SimpleErrorHandler{}

	if err := consumer.Consume(ctx, messageHandler, errorHandler); err != nil {
		log.Panicf("Error consuming messages: %v", err)
	}
}
