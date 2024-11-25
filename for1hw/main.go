package main

import (
	"archive/zip"
	"fmt"
	"os"
	"time"
)

func createTestVFSArchive(outputPath string) error {
	// Создаем ZIP-файл
	zipFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Создаем writer для ZIP
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Структура тестовой файловой системы
	directories := []string{
		"home",
		"home/user",
		"etc",
		"var",
	}

	files := []struct {
		path     string
		contents string
	}{
		{"home/user/readme.txt", "Тестовый файл в домашней директории"},
		{"etc/config.txt", "Тестовая конфигурация"},
		{"var/log.txt", "Тестовый лог"},
	}

	// Создаем директории
	for _, dir := range directories {
		header := &zip.FileHeader{
			Name:     dir + "/",
			Method:   zip.Store,
			Modified: time.Now(),
		}
		_, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
	}

	// Создаем файлы
	for _, file := range files {
		header := &zip.FileHeader{
			Name:     file.path,
			Method:   zip.Deflate,
			Modified: time.Now(),
		}
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = writer.Write([]byte(file.contents))
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	outputPath := "test_vfs.zip"
	err := createTestVFSArchive(outputPath)
	if err != nil {
		fmt.Printf("Ошибка создания архива: %v\n", err)
		return
	}
	fmt.Printf("Архив %s создан успешно\n", outputPath)
}
