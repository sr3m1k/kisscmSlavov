package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func replaceSpacesWithTab(inputPath, outputPath string) error {

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("error opening input file: %v", err)
	}
	defer func() {
		err := inputFile.Close()
		if err != nil {

		}
	}()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer func() {
		err := outputFile.Close()
		if err != nil {

		}
	}()

	reader := bufio.NewReader(inputFile)
	writer := bufio.NewWriter(outputFile)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				return fmt.Errorf("error reading line: %v", err)
			}
			// Последняя строка без символа новой строки
			line = strings.TrimRight(line, "\n")
			line = strings.ReplaceAll(line, "    ", "\t")
			_, err = writer.WriteString(line)
			if err != nil {
				return fmt.Errorf("error writing to output file: %v", err)
			}
			break
		}

		line = strings.ReplaceAll(line, "~", "\t")
		_, err = writer.WriteString(line)
		if err != nil {
			return fmt.Errorf("error writing to output file: %v", err)
		}
	}

	// Сбрасываем буфер и закрываем файл
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing buffer: %v", err)
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <input file> <output file>")
		return
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	err := replaceSpacesWithTab(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Replacement complete!")
	}
}
