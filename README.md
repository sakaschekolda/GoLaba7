1)	Создание TCP-сервера:
•	Реализуйте простой TCP-сервер, который слушает указанный порт и принимает входящие соединения.
•	Сервер должен считывать сообщения от клиента и выводить их на экран.
•	По завершении работы клиенту отправляется ответ с подтверждением получения сообщения.
```
  package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Server представляет собой TCP-сервер, который управляет соединениями и их завершением.
type Server struct {
	listener net.Listener   // Слушатель для принятия входящих соединений
	wg       sync.WaitGroup // Группа ожидания для отслеживания горутин
}

// NewServer создает новый сервер, который слушает указанный порт.
func NewServer(port string) (*Server, error) {
	// Создает слушатель для указанного порта
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}
	return &Server{listener: listener}, nil
}

// Start запускает сервер и ожидает входящих соединений.
func (s *Server) Start() {
	fmt.Println("Server is running, waiting for connection...")
	for {
		// Принимает входящее соединение
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Error occured while getting connetcion:", err)
			continue // Пропускаем ошибку и ожидаем следующее соединение
		}
		s.wg.Add(1)                 // Увеличиваем счетчик горутин
		go s.handleConnection(conn) // Обрабатываем соединение в отдельной горутине
	}
}

// handleConnection обрабатывает сообщение от клиента.
func (s *Server) handleConnection(conn net.Conn) {
	defer s.wg.Done()  // Уменьшаем счетчик горутин по завершении
	defer conn.Close() // Закрываем соединение по завершении обработки

	buffer := make([]byte, 1024) // Буфер для чтения сообщения
	for {
		// Читаем сообщение от клиента
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error occured while reading:", err)
			return // Завершаем горутину, если произошла ошибка чтения
		}
		message := string(buffer[:n]) // Преобразуем данные в строку
		fmt.Printf("Got a message: %s", message)

		// Отправка ответа клиенту
		_, err = conn.Write([]byte("Got a message\n"))
		if err != nil {
			fmt.Println("Error occured while getting a message:", err)
			return // Завершаем горутину при ошибке отправки
		}
	}
}

// Shutdown корректно завершает работу сервера и закрывает все соединения.
func (s *Server) Shutdown() {
	fmt.Println("Завершение работы сервера...")
	s.listener.Close() // Закрываем слушатель, чтобы больше не принимать соединения
	s.wg.Wait()        // Ждем завершения всех активных горутин
	fmt.Println("Все соединения закрыты.")
}

// main - точка входа в программу.
func main() {
	server, err := NewServer("8080") // Создаем новый сервер на порту 8080
	if err != nil {
		fmt.Println("Не удалось создать сервер:", err)
		return
	}

	// Обработка сигнала завершения
	signalChan := make(chan os.Signal, 1)                      // Канал для получения сигналов
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM) // Уведомляем о SIGINT и SIGTERM
	go func() {
		<-signalChan      // Ожидаем получения сигнала
		server.Shutdown() // Вызываем завершение работы сервера
		os.Exit(0)        // Завершаем программу
	}()

	server.Start() // Запускаем сервер
}

```
2)	Реализация TCP-клиента:
•	Разработайте TCP-клиента, который подключается к вашему серверу.
•	Клиент должен отправлять сообщение, введённое пользователем, и ожидать ответа.
•	После получения ответа от сервера клиент завершает соединение.
```
  package main

import (
	"fmt"
	"net"
)

// main - точка входа в программу клиента.
func main() {
	// Подключаемся к серверу на порту 8080
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error occured while connecting:", err)
		return
	}
	defer conn.Close() // Закрываем соединение по завершении работы клиента

	for {
		var message string
		fmt.Print("Enter message (or 'exit' to exit): ")
		_, err := fmt.Scanln(&message)       // Читаем ввод от пользователя
		if err != nil || message == "exit" { // Выход при ошибке или вводе "exit"
			break
		}

		_, err = conn.Write([]byte(message)) // Отправляем сообщение на сервер
		if err != nil {
			fmt.Println("Error occured while sending a message:", err)
			return
		}

		buffer := make([]byte, 1024) // Буфер для чтения ответа
		n, err := conn.Read(buffer)  // Читаем ответ от сервера
		if err != nil {
			fmt.Println("Error occured while getting a message:", err)
			return
		}
		// Печатаем ответ от сервера
		fmt.Printf("Server's response: %s", string(buffer[:n]))
	}
}

```
3-4)	Асинхронная обработка клиентских соединений:
•	Добавьте в сервер многопоточную обработку нескольких клиентских соединений.
•	Используйте горутины для обработки каждого нового соединения.
•	Реализуйте механизм graceful shutdown: сервер должен корректно завершать все активные соединения при остановке.
```
  package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Приветственное сообщение
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hi!"))
}

// Обработчик POST запросов с данными в формате JSON
func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Вывод данных в консоль
	fmt.Println("Got a data:", data)

	// Ответ клиенту
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Got a data!"))
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/data", dataHandler)

	fmt.Println("Running a server in port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error occured while running server:", err)
	}
}

// Invoke-RestMethod -Uri http://localhost:8080/data -Method POST -Headers @{ "Content-Type" = "application/json" } -Body '{"name": "John", "age": 30}'
```
5)	Создание HTTP-сервера:
•	Реализуйте базовый HTTP-сервер с обработкой простейших GET и POST запросов.
•	Сервер должен поддерживать два пути:
•	GET /hello — возвращает приветственное сообщение.
•	POST /data — принимает данные в формате JSON и выводит их содержимое в консоль.
```
  package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Middleware для логирования запросов
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()           // Запоминаем время начала обработки запроса
		next(w, r)                    // Позволяем следующему обработчику выполнить свою работу
		duration := time.Since(start) // Вычисляем продолжительность
		log.Printf(
			"Method: %s, URL: %s, Duration: %s\n",
			r.Method,
			r.URL.Path,
			duration,
		)
	}
}

// Обработчик для приветственного сообщения
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hi!"))
}

// Обработчик для получения данных
func dataHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var data map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		fmt.Println("Data that we got:", data)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Got a data!"))

	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Get request /data"))

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	// Регистрация маршрутов с использованием middleware
	http.HandleFunc("/hello", loggingMiddleware(helloHandler))
	http.HandleFunc("/data", loggingMiddleware(dataHandler))

	// Запуск сервера
	fmt.Println("Running server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error occured while running server:", err)
	}
}

// Invoke-RestMethod -Uri http://localhost:8080/data -Method POST -Headers @{ "Content-Type" = "application/json" } -Body '{"name": "John", "age": 30}'
```
7)	Добавление маршрутизации и middleware:
•	Реализуйте обработку нескольких маршрутов и добавьте middleware для логирования входящих запросов.
•	Middleware должен логировать метод, URL, и время выполнения каждого запроса.
```
  package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Middleware для логирования запросов
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()           // Запоминаем время начала обработки запроса
		next(w, r)                    // Позволяем следующему обработчику выполнить свою работу
		duration := time.Since(start) // Вычисляем продолжительность
		log.Printf(
			"Method: %s, URL: %s, Duration: %s\n",
			r.Method,
			r.URL.Path,
			duration,
		)
	}
}

// Обработчик для приветственного сообщения
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hi!"))
}

// Обработчик для получения данных
func dataHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var data map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		fmt.Println("Data that we got:", data)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Got a data!"))

	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Get request /data"))

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	// Регистрация маршрутов с использованием middleware
	http.HandleFunc("/hello", loggingMiddleware(helloHandler))
	http.HandleFunc("/data", loggingMiddleware(dataHandler))

	// Запуск сервера
	fmt.Println("Running server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error occured while running server:", err)
	}
}

// Invoke-RestMethod -Uri http://localhost:8080/data -Method POST -Headers @{ "Content-Type" = "application/json" } -Body '{"name": "John", "age": 30}'

```
9)	Веб-сокеты:
•	Реализуйте сервер на основе веб-сокетов для чата.
•	Клиенты должны подключаться к серверу, отправлять и получать сообщения.
•	Сервер должен поддерживать несколько клиентов и рассылать им сообщения, отправленные любым подключённым клиентом.
Go:
```
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

```
HTML:
```
  <!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>WebSocket chat</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            background-color: #f4f4f9;
        }
        h1 {
            color: #333;
        }
        #chat {
            border: 1px solid #ccc;
            padding: 10px;
            height: 400px;
            overflow-y: scroll;
            margin-bottom: 10px;
        }
        input[type="text"] {
            width: calc(100% - 22px);
            padding: 10px;
            margin-bottom: 10px;
        }
    </style>
</head>
<body>
    <h1>Чат</h1>
    <div id="chat"></div>
    <input type="text" id="message" placeholder="Enter message..." />
    
    <script>
        const chat = document.getElementById('chat');
        const messageInput = document.getElementById('message');

        // Создаем веб-сокет-соединение
        const ws = new WebSocket('ws://localhost:8080/ws');

        // Обрабатываем входящие сообщения
        ws.onmessage = function(event) {
            const message = document.createElement('div');
            message.textContent = event.data;
            chat.appendChild(message);
            chat.scrollTop = chat.scrollHeight; // Прокручиваем вниз, чтобы видеть последние сообщения
        };

        // Обрабатываем отправку сообщения
        messageInput.addEventListener('keyup', function(event) {
            if (event.key === 'Enter') {
                const message = messageInput.value;
                if (message) {
                    ws.send(message);
                    messageInput.value = ''; // Очищаем поле ввода
                }
            }
        });
    </script>
</body>
</html>
```
