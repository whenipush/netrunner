package controllers

import (
	"fmt"
	"log"
	"net/http"
	"netrunner/models"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket управление
var (
	clients   = make(map[*websocket.Conn]bool) // Подключённые клиенты
	broadcast = make(chan models.TaskStatus)   // Канал для отправки сообщений
	mu        = &sync.Mutex{}                  // Мьютекс для синхронизации
)

func GeneralRoot(c *gin.Context) {
	c.JSON(200, map[string]interface{}{
		"netrunnerVersion": "0.1.0",
		"netrunnerStatus":  "up",
	})
}

// HandleWebSocket обрабатывает WebSocket-соединения
func HandleWebSocket(c *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Разрешить соединения с любого источника
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Ошибка обновления до WebSocket: %v", err)
		return
	} else {
		log.Printf("Успешное подключение к websocket: %v", conn.RemoteAddr())
	}
	defer conn.Close()

	mu.Lock()
	clients[conn] = true
	mu.Unlock()
	/*var task []models.TaskStatus

	if err := database.DB.Preload("Hosts").Find(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	conn.WriteJSON(task)*/
	// Ожидание сообщений
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			break
		}
	}
}

// broadcastTask отправляет задачу всем подключённым WebSocket-клиентам
func BroadcastTask(task models.TaskStatus) {
	mu.Lock()
	defer mu.Unlock()

	for client := range clients {
		err := client.WriteJSON(task)
		if err != nil {
			log.Printf("Ошибка отправки через WebSocket: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// Не знаю зачем нам кастомные скрипты, мы их нормально обрабатывать не сможем
func UploadScript(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(500, fmt.Sprintf("get form file err: %v", err))
		return
	}

	if file.Filename == "" {
		c.JSON(400, "File name is required")
		return
	}
	filesplit := strings.Split(file.Filename, ".")
	fileExt := filesplit[len(filesplit)-1]

	if fileExt != "lua" {
		c.JSON(400, gin.H{"Error": "Only lua files are allowed"})
		return
	}

	filepath := "./scripts/" + file.Filename
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(500, fmt.Sprintf("upload file err: %v", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "filename": file.Filename})
}
