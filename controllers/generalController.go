package controllers

import (
	"fmt"
	"log"
	"net/http"
	"netrunner/database"
	"netrunner/models"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func UploadDatabaseBDU(c *gin.Context) {
	// Ожидаемые файлы
	files := []string{"export.xml", "vullist.xlsx"}
	savedFiles := make(map[string]string) // Для хранения путей к загруженным файлам

	// Загрузка файлов
	for _, expectedFile := range files {
		file, err := c.FormFile(expectedFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Missing required file: %s", expectedFile),
			})
			return
		}

		// Формируем путь для сохранения файла
		filepath := "./vulners/" + expectedFile
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to save file '%s': %v", expectedFile, err),
			})
			return
		}

		savedFiles[expectedFile] = filepath // Сохраняем путь к файлу
	}

	// Передача файлов в обработчик
	//if _, err := handlers.UpdateDatabaseBDU(savedFiles["vullist.xlsx"], savedFiles["export.xml"]); err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"error": fmt.Sprintf("Failed to process uploaded files: %v", err),
	//	})
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{"message": "Files uploaded and processed successfully"})
}

// WebSocket управление
var (
	clients   = make(map[*websocket.Conn]bool) // Подключённые клиенты
	broadcast = make(chan models.TaskStatus)   // Канал для отправки сообщений
	mu        = &sync.Mutex{}                  // Мьютекс для синхронизации
)

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
	var task []models.TaskStatus

	if err := database.DB.Preload("Hosts").Find(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	conn.WriteJSON(task)
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
