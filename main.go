package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"netrunner/controllers"
	"netrunner/database"
)

func main() {
	database.Connect()
	r := gin.Default()
	r.Use(cors.Default())

	// Host endpoints
	r.GET("/api/v1/host", controllers.GetAllHost)        // Получить все хосты
	r.POST("/api/v1/host", controllers.CreateHost)       // Создать новый хост
	r.PUT("/api/v1/host/:id", controllers.UpdateHost)    // Изменить хост по ID
	r.DELETE("/api/v1/host/:id", controllers.DeleteHost) // Удалить хост по ID

	// Group endpoints

	r.GET("/api/v1/group", controllers.GetAllGroup)        // Получить все группы
	r.POST("/api/v1/group", controllers.CreateGroup)       // Создать новую группу
	r.PUT("/api/v1/group/:id", controllers.UpdateGroup)    // Изменить группу по ID
	r.DELETE("/api/v1/group/:id", controllers.DeleteGroup) // Удалить группу по ID

	// Host-Group endpoints
	r.POST("/api/v1/host-add-group", controllers.AddHostToGroupHandler) // Добавить хосты к группам

	// other endpoints

	r.POST("/api/v1/upload-script", controllers.UploadScript) // Загрузить скрипт nmap на сервер
	r.POST("/api/v1/nmap", controllers.ProcessNmapRequest)    // Запустить nmap на указанных хостах с указанным скриптом

	// Получение результатов Nmap

	r.GET("/api/v1/last-nmap", controllers.GetLastNmap)

	r.GET("/api/v1/all-nmap", controllers.GetAllNmap)
	//r.GET("/api/v1/name-nmap/:filename", controllers.GetReportByName)

	r.Run(":3000")

}
