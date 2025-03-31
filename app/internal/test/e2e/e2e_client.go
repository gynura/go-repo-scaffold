package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
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

type Event struct {
	EventType string          `json:"eventType"`
	Payload   json.RawMessage `json:"payload"`
}

// TODO: should refactor this e2e approach, for now its coupled to the local deploy
const (
	userServiceExchangeName = "user_service"
	rabbitMQURL             = "amqp://guest:guest@localhost:5672/"
	createUserURL           = "http://localhost:8080/users"
	queueName               = "e2e_tests_queue"
)

func main() {
	eventChan := make(chan Event, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go subscribeToEvents(eventChan, ctx)

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
	case event := <-eventChan:
		fmt.Println("Received event:", event)
		// Check if event contains the expected user ID
		if !bytes.Contains(event.Payload, []byte(user.ID)) {
			fmt.Println("Test failed: Expected user ID not found in event payload")
			os.Exit(1)
		}
		fmt.Println("Test passed: User creation event received successfully")
	case <-ctx.Done():
		fmt.Println("Test failed: No event received within timeout")
		os.Exit(1)
	}
}

// subscribeToEvents connects to RabbitMQ and listens for events
func subscribeToEvents(eventChan chan<- Event, ctx context.Context) {
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

	err = ch.QueueBind(q.Name, "#", userServiceExchangeName, false, nil)
	if err != nil {
		fmt.Println("Error binding RabbitMQ queue:", err)
		return
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		fmt.Println("Failed to consume messages:", err)
		return
	}

	for {
		select {
		case msg := <-msgs:
			var event Event
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				fmt.Println("Failed to unmarshal event:", err)
				continue
			}
			eventChan <- event
			return
		case <-ctx.Done():
			return
		}
	}
}

// createUser sends an HTTP request to create a user
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
