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

// EventType defines the different types of user events, decoupled from the ones defined inside the service code
type EventType string

const (
	EventUserCreated EventType = "user.created"
	EventUserEdited  EventType = "user.edited"
	EventUserRemoved EventType = "user.removed"
)

// TODO: should refactor this e2e approach, for now its coupled to the local default deploy
const (
	userServiceExchangeName = "user_service"
	rabbitMQURL             = "amqp://guest:guest@localhost:5672/"
	userServiceBaseURL      = "http://localhost:8080/users"
	queueName               = "e2e_tests_queue"
	testTimeout             = 10 * time.Second
)

func main() {
	fmt.Println("=== Testing User Creation ===")
	user := createUserAndVerifyEvent()

	fmt.Println("\n=== Testing User Edit ===")
	editUserAndVerifyEvent(user)

	fmt.Println("\n=== Testing User Deletion ===")
	deleteUserAndVerifyEvent(user)

	fmt.Println("\n=== E2E Tests passed! ===")
}

func createUserAndVerifyEvent() User {
	payloadChan := make(chan ExpectedMessagePayload, 1)
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Start listening for the user.created event
	go subscribeToEvents(payloadChan, ctx, EventUserCreated)

	// Create a new user
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

	// Wait for the event
	select {
	case payload := <-payloadChan:
		fmt.Printf("Received user.created event with payload: %+v\n", payload)

		// Check if event contains the expected user ID
		if payload.ID != user.ID {
			fmt.Printf("Test failed: Expected user ID %s but got %s\n", user.ID, payload.ID)
			os.Exit(1)
		}

		fmt.Println("Test passed: User creation event received successfully with correct ID")
		return user

	case <-ctx.Done():
		fmt.Println("Test failed: No user.created event received within timeout")
		os.Exit(1)
		return User{} // Will never reach here, just to satisfy the compiler
	}
}

func editUserAndVerifyEvent(user User) {
	payloadChan := make(chan ExpectedMessagePayload, 1)
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Start listening for the user.edited event
	go subscribeToEvents(payloadChan, ctx, EventUserEdited)

	// Update the user's information
	updatedUser := User{
		ID: user.ID,
		SpecificData: map[string]interface{}{
			"firstName": "Sasuke",
			"lastName":  "Uchiha",
			"nickname":  "edgeLordNew",
			"password":  "IhateNaruto3",
			"email":     user.SpecificData["email"].(string),
			"country":   "LeafVillage",
		},
	}

	if err := editUser(updatedUser); err != nil {
		fmt.Println("Failed to edit user:", err)
		os.Exit(1)
	}

	// Wait for the event
	select {
	case payload := <-payloadChan:
		fmt.Printf("Received user.edited event with payload: %+v\n", payload)

		// Check if event contains the expected user ID
		if payload.ID != user.ID {
			fmt.Printf("Test failed: Expected user ID %s but got %s\n", user.ID, payload.ID)
			os.Exit(1)
		}

		fmt.Println("Test passed: User edit event received successfully with correct ID")

	case <-ctx.Done():
		fmt.Println("Test failed: No user.edited event received within timeout")
		os.Exit(1)
	}
}

func deleteUserAndVerifyEvent(user User) {
	payloadChan := make(chan ExpectedMessagePayload, 1)
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Start listening for the user.removed event
	go subscribeToEvents(payloadChan, ctx, EventUserRemoved)

	if err := deleteUser(user.ID); err != nil {
		fmt.Println("Failed to delete user:", err)
		os.Exit(1)
	}

	// Wait for the event
	select {
	case payload := <-payloadChan:
		fmt.Printf("Received user.removed event with payload: %+v\n", payload)

		// Check if event contains the expected user ID
		if payload.ID != user.ID {
			fmt.Printf("Test failed: Expected user ID %s but got %s\n", user.ID, payload.ID)
			os.Exit(1)
		}

		fmt.Println("Test passed: User removal event received successfully with correct ID")

	case <-ctx.Done():
		fmt.Println("Test failed: No user.removed event received within timeout")
		os.Exit(1)
	}
}

// subscribeToEvents connects to RabbitMQ and listens for a specific event type
func subscribeToEvents(payloadChan chan<- ExpectedMessagePayload, ctx context.Context, eventType EventType) {
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

	// Use a unique queue name for each test to avoid message conflicts
	uniqueQueueName := fmt.Sprintf("%s_%s_%s", queueName, eventType, uuid.New().String())

	q, err := ch.QueueDeclare(uniqueQueueName, false, true, true, false, nil)
	if err != nil {
		fmt.Println("Error declaring RabbitMQ queue:", err)
		return
	}

	// Bind only to the specific event routing key
	err = ch.QueueBind(q.Name, string(eventType), userServiceExchangeName, false, nil)
	if err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
	}

	fmt.Printf("Listening for '%s' events on queue '%s'\n", eventType, q.Name)

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
			if msg.RoutingKey == string(eventType) {
				var payload ExpectedMessagePayload
				if err := json.Unmarshal(msg.Body, &payload); err != nil {
					fmt.Printf("Failed to unmarshal payload: %v\n", err)
					continue
				}

				payloadChan <- payload
				return // Stop listening after getting the expected event
			} else {
				fmt.Printf("Ignoring message with unexpected routing key: %s (expected %s)\n",
					msg.RoutingKey, eventType)
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

	resp, err := http.Post(userServiceBaseURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response %d: %s", resp.StatusCode, string(body))
	}

	fmt.Println("User created successfully")
	return nil
}

func editUser(user User) error {
	data, err := json.Marshal(map[string]interface{}{
		"specificData": user.SpecificData,
	})
	if err != nil {
		return err
	}

	// Create PUT request
	url := fmt.Sprintf("%s/%s", userServiceBaseURL, user.ID)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response %d: %s", resp.StatusCode, string(body))
	}

	fmt.Println("User edited successfully")
	return nil
}

func deleteUser(userID string) error {
	url := fmt.Sprintf("%s/%s", userServiceBaseURL, userID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response %d: %s", resp.StatusCode, string(body))
	}

	fmt.Println("User deleted successfully")
	return nil
}
