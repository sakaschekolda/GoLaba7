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
