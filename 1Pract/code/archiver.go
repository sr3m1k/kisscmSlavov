package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func validateArgs() (string, string, error) {
	if len(os.Args) < 3 {
		return "", "", fmt.Errorf("usage: go run main.go <directory> <extension>")
	}
	directory := os.Args[1]
	extension := os.Args[2]
	return directory, extension, nil
}

func createZipArchive(zipPath string) (*os.File, *zip.Writer, error) {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating zip file: %v", err)
	}
	zipWriter := zip.NewWriter(zipFile)
	return zipFile, zipWriter, nil
}

func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", filePath, err)
	}
	defer func() {
		err := file.Close()
		if err != nil {

		}
	}()

	zipEntry, err := zipWriter.Create(filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("error creating zip entry for file %s: %v", filePath, err)
	}

	_, err = io.Copy(zipEntry, file)
	if err != nil {
		return fmt.Errorf("error writing file %s to zip: %v", filePath, err)
	}

	return nil
}

// Поиск файлов и их добавление в архив
func archiveFiles(directory, extension, zipPath string) error {
	zipFile, zipWriter, err := createZipArchive(zipPath)
	if err != nil {
		return err
	}
	defer func() {
		err := zipWriter.Close()
		if err != nil {

		}
	}()
	defer func() {
		err := zipFile.Close()
		if err != nil {

		}
	}()

	err = filepath.WalkDir(directory, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), extension) {
			if err := addFileToZip(zipWriter, path); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking the path %v: %v", directory, err)
	}

	return nil
}

func main() {
	directory, extension, err := validateArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	zipPath := "tar.zip"
	if err := archiveFiles(directory, extension, zipPath); err != nil {
		fmt.Printf("Error archiving files: %v\n", err)
	} else {
		fmt.Println("Archiving completed successfully!")
	}
}
