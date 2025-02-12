package handlers

import (
	"fmt"
	"log"
	"os"
)

func LogNmapError(title string) (*log.Logger, error) {
	// Открываем файл для логирования
	path := fmt.Sprintf("logs/nmap/%s.log", title)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Создаем новый логгер
	logger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)

	return logger, nil
}
