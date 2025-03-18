package socket

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

// Message represents the structure of a chat message
type Message struct {
	ID        string `json:"id"`        // Unique ID for the message
	Sender    string `json:"sender"`    // "pharmacist" or "doctor"
	Text      string `json:"text"`      // Message content
	Timestamp string `json:"timestamp"` // Timestamp of the message
}

var clients = make(map[*websocket.Conn]bool) // Track connected clients
var broadcast = make(chan Message)           // Channel for broadcasting messages

// SetWebSocketRoutes registers WebSocket handlers
func SetWebSocketRoutes(app *fiber.App) {
	// WebSocket endpoint
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket handler
	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		log.Println("New WebSocket connection established")
		defer func() {
			conn.Close()
			delete(clients, conn)
			log.Println("WebSocket connection closed")
		}()

		// Register client
		clients[conn] = true

		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					log.Println("WebSocket connection closed unexpectedly:", err)
				} else {
					log.Println("Error reading message:", err)
				}
				break
			}

			// Assign a unique ID to the message
			msg.ID = uuid.New().String()

			log.Printf("Received message: %+v\n", msg)

			// Broadcast the message to all clients
			broadcast <- msg
		}
	}))

	// Start the message handler
	go handleMessages()
}

// handleMessages broadcasts messages to all connected clients
func handleMessages() {
	for {
		msg := <-broadcast

		// Send the message to all connected clients
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Error broadcasting message:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
