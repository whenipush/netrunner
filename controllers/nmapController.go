package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type SQLmapRequest struct {
	Target  string            `json:"target" binding:"required"` // URL для сканирования
	Options map[string]string `json:"options"`                   // Дополнительные параметры для SQLmap
}

type NmapRequest struct {
	Hosts  []string `form:"hosts" binding:"required"`
	Script string   `form:"script" binding:"required"`
}

func startNmap(host string, script string) {
	// Заглушка для запуска nmap

	// nmap -sV --script vuln -oX scan_results.xml 192.168.20.218
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02")

	// Формирование имени файла отчёта
	report := fmt.Sprintf("report/report_%s.xml", currentDate)

	// Формирование команды NMap
	nmapCmd := exec.Command("nmap", "-sV", "--script", "vuln", "-oX", report, host)

	// Выполнение команды
	output, err := nmapCmd.CombinedOutput()
	if err != nil {
		log.Println("Ошибка запуска NMap:", err)
		return
	}

	// Вывод результата
	fmt.Println("Результаты NMap сохранены в файл:", report)
	fmt.Println(string(output))
	fmt.Println(script)

}
func ProcessNmapRequest(c *gin.Context) {

	var input NmapRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": "Запрос в обработке"})

	go func() {
		for _, host := range input.Hosts {
			startNmap(host, input.Script)
		}
		log.Println("Сделано!")
	}()

}

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

func GetLastNmap(c *gin.Context) {
	// Путь к папке с файлами
	dir := "reports"

	files, err := os.ReadDir(dir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read directory"})
		return
	}

	if len(files) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No files found in directory"})
		return
	}

	// Найти самый последний файл
	var latestFile os.DirEntry
	var latestModTime time.Time

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		if info.ModTime().After(latestModTime) {
			latestFile = file
			latestModTime = info.ModTime()
		}
	}

	if latestFile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No files found in directory"})
		return
	}

	filePath := filepath.Join(dir, latestFile.Name())
	c.File(filePath)
}

func GetAllNmap(c *gin.Context) {
	// Путь к папке с файлами
	dir := "reports"

	files, err := os.ReadDir(dir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read directory"})
		return
	}

	if len(files) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No files found in directory"})
		return
	}

	// Сортируем файлы по дате модификации
	sort.Slice(files, func(i, j int) bool {
		infoI, _ := files[i].Info()
		infoJ, _ := files[j].Info()
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	// Создаем массив с путями к файлам
	var filePaths []string
	for _, file := range files {
		filePaths = append(filePaths, filepath.Join(dir, file.Name()))
	}

	c.JSON(http.StatusOK, gin.H{"files": filePaths})
}
