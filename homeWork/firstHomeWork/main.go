package main

import (
	"archive/zip"
	"firstHomeWork/initializer"
	"firstHomeWork/shellProgram"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Структура для хранения конфигурации
type Config struct {
	User       string `toml:"user"`
	Hostname   string `toml:"hostname"`
	VFSZipPath string `toml:"vfs_zip_path"`
	LogFile    string `toml:"log_file"`
}

// Структура для логирования

// Функция для загрузки конфигурации из YAML

// Функция для распаковки ZIP архива
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

// Функция для логирования действий

func main() {

	config := Config{
		User:       "root",
		Hostname:   "localhost",
		VFSZipPath: "./vfs_zip",
		LogFile:    "./log",
	}

	// Загружаем конфигурацию
	err := initializer.FillConfig("config.toml", config)
	if err != nil {
		log.Println("fail to fill config")
	}

	// Открываем файл для логов
	logFile, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer func() {
		err := logFile.Close()
		if err != nil {

		}
	}()

	// Создаем временную директорию для работы с виртуальной файловой системой
	vfsPath := "./vfs"
	if _, err := os.Stat(vfsPath); !os.IsNotExist(err) {
		os.RemoveAll(vfsPath)
	}
	os.Mkdir(vfsPath, os.ModePerm)

	// Распаковываем виртуальную файловую систему
	err = unzip(config.VFSZipPath, vfsPath)
	if err != nil {
		fmt.Println("Error unpacking VFS:", err)
		return
	}

	// Инициализируем эмулятор shellRunner
	shell := shellProgram.Shell{
		Config:      config,
		CurrentPath: vfsPath,
		LogFile:     logFile,
	}

	// Запуск эмулятора
	shell.Run()
}
