package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// User is decoupled from the User defined in the app code so the e2e test is agnostic to the service's structs
type User struct {
	ID           string                 `json:"id"`
	SpecificData map[string]interface{} `json:"specificData"`
}

// ExpectedMessagePayload represents the actual payload we receive which is just {"id": "user-id"}
type ExpectedMessagePayload struct {
	ID string `json:"id"`
}

// TODO: should refactor this e2e approach, for now its coupled to the local default deploy
const (
	userServiceExchangeName = "user_service"
	rabbitMQURL             = "amqp://guest:guest@localhost:5672/"
	createUserURL           = "http://localhost:8080/users"
	queueName               = "e2e_tests_queue"
	expectedRoutingKey      = "user.created"
)

func main() {
	payloadChan := make(chan ExpectedMessagePayload, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go subscribeToEvents(payloadChan, ctx)

	user := User{
		ID: uuid.New().String(),
		SpecificData: map[string]interface{}{
			"firstName": "Sasuke",
			"lastName":  "Uchiha",
			"nickname":  uuid.New().String(),
			"password":  "IhateNaruto2",
			"email":     uuid.New().String() + "@gmail.com",
			"country":   "LeafVillage",
		},
	}

	if err := createUser(user); err != nil {
		fmt.Println("Failed to create user:", err)
		os.Exit(1)
	}

	select {
	case payload := <-payloadChan:
		fmt.Printf("Received user.created event with payload: %+v\n", payload)

		// Check if event contains the expected user ID
		if payload.ID != user.ID {
			fmt.Printf("Test failed: Expected user ID %s but got %s\n", user.ID, payload.ID)
			os.Exit(1)
		}

		fmt.Println("Test passed: User creation event received successfully with correct ID")

	case <-ctx.Done():
		fmt.Println("Test failed: No event received within timeout")
		os.Exit(1)
	}
}

// subscribeToEvents connects to RabbitMQ and listens for events
func subscribeToEvents(payloadChan chan<- ExpectedMessagePayload, ctx context.Context) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		fmt.Println("Error connecting to RabbitMQ:", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Error opening RabbitMQ channel:", err)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		fmt.Println("Error declaring RabbitMQ queue:", err)
		return
	}

	// Bind only to the user.created routing key
	err = ch.QueueBind(queueName, expectedRoutingKey, userServiceExchangeName, false, nil)
	if err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
	}

	retrievedMessages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		fmt.Println("Failed to consume messages:", err)
		return
	}

	for {
		select {
		case msg := <-retrievedMessages:
			fmt.Printf("Received message with routing key: %s\n", msg.RoutingKey)
			fmt.Printf("Raw message body: %s\n", string(msg.Body))

			// We check if this is the expected event based on the routing key
			if msg.RoutingKey == expectedRoutingKey {
				var payload ExpectedMessagePayload
				if err := json.Unmarshal(msg.Body, &payload); err != nil {
					fmt.Printf("Failed to unmarshal payload: %v\n", err)
					continue
				}

				payloadChan <- payload
				return // Stop listening after getting the expected event
			} else {
				fmt.Printf("Ignoring message with unexpected routing key: %s\n", msg.RoutingKey)
			}

		case <-ctx.Done():
			return
		}
	}
}

func createUser(user User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	resp, err := http.Post(createUserURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response %d: %s", resp.StatusCode, string(body))
	}

	fmt.Println("User created successfully.")
	return nil
}
