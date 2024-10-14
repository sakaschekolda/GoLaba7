package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Client представляет клиента чата
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// ChatServer представляет сервер чата
type ChatServer struct {
	clients   map[*Client]bool
	broadcast chan []byte
	mutex     sync.Mutex
}

// Создаем новый сервер чата
func newChatServer() *ChatServer {
	return &ChatServer{
		clients:   make(map[*Client]bool),
		broadcast: make(chan []byte),
	}
}

// Добавляем нового клиента
func (cs *ChatServer) addClient(c *Client) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	cs.clients[c] = true
}

// Удаляем клиента
func (cs *ChatServer) removeClient(c *Client) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	delete(cs.clients, c)
	close(c.send)
}

// Обрабатываем сообщения
func (cs *ChatServer) handleMessages() {
	for {
		message := <-cs.broadcast
		for client := range cs.clients {
			select {
			case client.send <- message:
			default:
				cs.removeClient(client)
			}
		}
	}
}

// Запускаем сервер
func (cs *ChatServer) run() {
	go cs.handleMessages()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error occured while connecting:", err)
			return
		}
		client := &Client{conn: conn, send: make(chan []byte)}
		cs.addClient(client)

		go cs.handleClient(client)
	})

	fmt.Println("Running server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error occured while running server:", err)
	}
}

// Обработка клиента
func (cs *ChatServer) handleClient(client *Client) {
	defer func() {
		client.conn.Close()
		cs.removeClient(client)
	}()

	go func() {
		for {
			_, message, err := client.conn.ReadMessage()
			if err != nil {
				return
			}
			cs.broadcast <- message
		}
	}()

	for message := range client.send {
		if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

func main() {
	server := newChatServer()
	server.run()
}
