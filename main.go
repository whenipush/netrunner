package main

import (
	"netrunner/controllers"
	"netrunner/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// kek
func main() {
	database.Connect()
	r := gin.Default()
	//r.Use(cors.New(cors.Config{
	//	AllowOrigins:     []string{"http://localhost:5500"},                 // Разрешите фронтенд
	//	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"}, // Разрешенные методы
	//	AllowHeaders:     []string{"Content-Type", "Authorization"},         // Разрешенные заголовки
	//	AllowCredentials: true,                                              // Разрешите использование куки/сессий
	//}))

	r.Use(cors.Default())

	// Группа для работы с хостами (Hosts)
	hostRoutes := r.Group("/api/v1/host")
	{
		// POST /api/v1/host - Создать хост
		hostRoutes.POST("/", controllers.CreateHost)

		// GET /api/v1/host - Получить все хосты
		hostRoutes.GET("/", controllers.GetAllHost)

		// GET /api/v1/host/search?ip=1.2.3.4 - Найти хост по IP
		hostRoutes.GET("/search", controllers.GetHostByID)

		// PUT /api/v1/host/:id - Обновить хост по ID
		hostRoutes.PUT("/:id", controllers.UpdateHost)

		// DELETE /api/v1/host/:id - Удалить хост по ID
		hostRoutes.DELETE("/:id", controllers.DeleteHost)
	}

	// Эндпоинты для работы с группами (Groups)
	groupRoutes := r.Group("/api/v1/group")
	{
		// POST /api/v1/group - Создать группу
		groupRoutes.POST("/", controllers.CreateGroup)

		// GET /api/v1/group - Получить все группы
		groupRoutes.GET("/", controllers.GetAllGroup)

		// GET /api/v1/group/search?group=GroupName - Найти группу по имени
		groupRoutes.GET("/search", controllers.GetGroupByName)

		// PUT /api/v1/group/:id - Обновить группу по ID
		groupRoutes.PUT("/:id", controllers.UpdateGroup)

		// DELETE /api/v1/group/:id - Удалить группу по ID
		groupRoutes.DELETE("/:id", controllers.DeleteGroup)
	}

	// POST /api/v1/add-hosts-to-groups - Добавить хосты в группы
	r.POST("/api/v1/add-hosts-to-groups", controllers.AddHostToGroup)

	// Эндпоинты для работы с задачами (Tasks)

	// POST /api/v1/task - Создать задачу
	r.POST("/api/v1/task", controllers.CreateTask)

	// GET /api/v1/task-status/:number_task - Проверить статус задачи
	r.GET("/api/v1/task-status/:number_task", controllers.GetTaskStatus)

	// DELETE /api/v1/task/:number_task - Удалить задачу
	r.DELETE("/api/v1/task/:number_task", controllers.DeleteTask)

	// GET /api/v1/task-all - Получить все задачи
	r.GET("/api/v1/task", controllers.GetTaskAll)

	// Остальные эндпоинты

	// POST /api/v1/upload-script - Загрузить скрипт Nmap
	r.POST("/api/v1/upload-script", controllers.UploadScript)

	// GET /api/v1/pentest/:number_task - Получить отчет пентеста
	r.GET("/api/v1/pentest/:number_task", controllers.GetPentestJsonByNumberTask)

	// WebSocket подключение
	// GET /api/v1/ws - Подключение WebSocket
	r.GET("/api/v1/ws", controllers.HandleWebSocket)

	r.POST("/api/v1/ping", controllers.PingHosts)

	// Запуск сервера на порту 3001
	r.Run(":3001")

}
