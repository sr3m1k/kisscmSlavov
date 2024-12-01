package main

import (
	"bufio"
	"fmt"
	"github.com/go-yaml/yaml"
	"os"
	"strconv"
	"strings"
)

type Instruction struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// Ассемблер
func assemble(inputFile, outputFile, logFile string) error {
	// Открыть файл исходной программы
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Прочитать исходный файл
	var instructions []Instruction
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		// Преобразуем строку в команду
		cmd, err := processLine(parts)
		if err != nil {
			return err
		}
		instructions = append(instructions, cmd)
	}

	// Создать бинарный файл
	binFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer binFile.Close()

	// Записать бинарные данные
	for _, inst := range instructions {
		binFile.Write([]byte(inst.Value))
	}

	// Логирование
	logFile, err := os.Create(logFile)
	if err != nil {
		return err
	}
	defer logFile.Close()

	yamlEncoder := yaml.NewEncoder(logFile)
	if err := yamlEncoder.Encode(instructions); err != nil {
		return err
	}

	return nil
}

// Функция для обработки каждой строки программы
func processLine(parts []string) (Instruction, error) {
	// Пример для константы
	if parts[0] == "LOAD_CONST" {
		a, err := strconv.Atoi(parts[1])
		if err != nil {
			return Instruction{}, err
		}
		b, err := strconv.Atoi(parts[2])
		if err != nil {
			return Instruction{}, err
		}
		return Instruction{
			Key:   fmt.Sprintf("LOAD_CONST_%d", a),
			Value: fmt.Sprintf("%02X%02X%02X", a, b, 0x00),
		}, nil
	}
	// Аналогично для других команд
	return Instruction{}, fmt.Errorf("неизвестная команда: %s", parts[0])
}

func main() {
	// Получаем параметры из командной строки
	if len(os.Args) != 4 {
		fmt.Println("Использование: program <inputFile> <outputFile> <logFile>")
		return
	}
	inputFile := os.Args[1]
	outputFile := os.Args[2]
	logFile := os.Args[3]

	// Ассемблируем программу
	err := assemble(inputFile, outputFile, logFile)
	if err != nil {
		fmt.Println("Ошибка ассемблера:", err)
	}
}
