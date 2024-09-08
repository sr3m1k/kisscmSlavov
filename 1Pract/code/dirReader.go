package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func listEmptyFiles(directory string) (int, error) {
	// Открываем директорию
	entries, err := os.ReadDir(directory)
	if err != nil {
		return 0, fmt.Errorf("error reading directory: %v", err)
	}

	// Проходим по всем файлам в директории
	for _, entry := range entries {
		if entry.IsDir() {
			continue // Пропускаем поддиректории
		}

		// Проверяем расширение файла
		if filepath.Ext(entry.Name()) == ".txt" {
			filePath := filepath.Join(directory, entry.Name())
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				return 0, fmt.Errorf("error stating file %s: %v", filePath, err)
			}

			// Проверяем размер файла
			if fileInfo.Size() == 0 {
				fmt.Println(entry.Name())
			} else {
				fmt.Println(len(entry.Name()), "no such files")
			}
		}
		return len(entry.Name()), nil
	}

	return 0, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <directory>")
		return
	}

	directory := os.Args[1]

	count, err := listEmptyFiles(directory)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Found %d files in %s\n", count, directory)
}
