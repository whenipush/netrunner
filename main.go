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
	// govno
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
		// Создание хоста
		// POST /host
		hostRoutes.POST("/", controllers.CreateHost)

		// Получение всех хостов
		// GET /host
		hostRoutes.GET("/", controllers.GetAllHost)

		// Получение хоста по IP (query-параметр ip=...)
		// GET /host/search?ip=1.2.3.4
		hostRoutes.GET("/search", controllers.GetHostByID)

		// Обновление хоста по ID
		// PUT /host/:id
		hostRoutes.PUT("/:id", controllers.UpdateHost)

		// Удаление хоста по ID
		// DELETE /host/:id
		hostRoutes.DELETE("/:id", controllers.DeleteHost)
	}

	// Эндпоинт для добавления хостов в группы
	// POST /add-hosts-to-groups
	r.POST("/add-hosts-to-groups", controllers.AddHostToGroupHandler)
	// Group endpoints

	groupRoutes := r.Group("/api/v1/group")
	{
		// Создание группы
		groupRoutes.POST("/", controllers.CreateGroup)

		// Получение всех групп
		groupRoutes.GET("/", controllers.GetAllGroup)

		// Получение группы по имени (через query-параметр "group=...")
		// Пример: GET /group/search?group=GroupName
		groupRoutes.GET("/search", controllers.GetGroupByName)

		// Обновление группы по ID (прокидываем id как /group/:id)
		groupRoutes.PUT("/:id", controllers.UpdateGroup)

		// Удаление группы по ID (прокидываем id как /group/:id)
		groupRoutes.DELETE("/:id", controllers.DeleteGroup)
	}

	// other endpoints nmap

	r.GET("/api/v1/pentest/:number_task", controllers.GetPentestJsonByNumberTask)
	r.POST("/api/v1/upload-script", controllers.UploadScript) // Загрузить скрипт nmap на сервер
	r.POST("/api/v1/nmap", controllers.ProcessNmapRequest)
	r.GET("/api/v1/task-status/:number_task", controllers.GetTaskStatus)
	r.DELETE("/api/v1/task/:number_task", controllers.DeleteTask)
	r.GET("/api/v1/task-all", controllers.GetTaskAll)

	// Получение результатов Nmap
	r.GET("/api/v1/last-nmap", controllers.GetLastNmap)
	r.GET("/api/v1/all-nmap", controllers.GetAllNmap)
	//r.GET("/api/v1/name-nmap/:filename", controllers.GetReportByName)

	// Загружаем 2 файла БДУ ФСТЭКА, сохраняем и обрабатываем
	// Не трогать, тут как часики работает
	r.POST("/api/v1/update-database-bdu", controllers.UploadDatabaseBDU)
	r.GET("/api/v1/ws", controllers.HandleWebSocket)
	r.Run(":3001")

}
