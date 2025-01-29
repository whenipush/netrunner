package controllers

import (
	"fmt"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tatsushid/go-fastping"
)

var Pinger = fastping.NewPinger()

func PingHosts(c *gin.Context) {
	activeHosts := make([]string, 0)

	Pinger.OnRecv = func(i *net.IPAddr, d time.Duration) {
		activeHosts = append(activeHosts, i.String())
	}

	if err := Pinger.Run(); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Failed to ping hosts: %s", err.Error())})
		return
	}
	// Сканируем сеть по указанному диапазону

	// Отправляем список активных хостов в ответ
	c.JSON(200, gin.H{
		"activeHosts": activeHosts,
	})
}
