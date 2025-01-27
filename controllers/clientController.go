package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ClientConnection struct {
	IP         string
	Connection *websocket.Conn
}

var clientsConnection map[string]ClientConnection = make(map[string]ClientConnection)

func HandleClientSocketConnection(c *gin.Context) {
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
	adrr := conn.RemoteAddr()
	client := ClientConnection{
		IP:         adrr.String(),
		Connection: conn,
	}
	clientsConnection[adrr.String()] = client
	defer func() {
		conn.Close()
		delete(clientsConnection, adrr.String())
	}()
	for {
		_, message, err := conn.ReadMessage()
		log.Printf("Client message: %s", message)
		if err != nil {
			delete(clientsConnection, adrr.String())
			break
		}
	}
}
