package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
)

type VirtualMachine struct {
	Memory []byte
}

func (vm *VirtualMachine) Execute(binFile string) error {
	// Открываем бинарный файл
	data, err := ioutil.ReadFile(binFile)
	if err != nil {
		return err
	}

	// Выполняем команды
	for i := 0; i < len(data); i += 3 {
		// Анализируем команду
		cmd := data[i : i+3]
		fmt.Printf("Executing command: %v\n", cmd)
		// Пример выполнения команды
		if cmd[0] == 0xEC { // Загрузка константы
			// Получаем константу
			value := int(cmd[1])<<8 + int(cmd[2])
			vm.Memory = append(vm.Memory, byte(value))
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Использование: interpreter <binaryFile> <resultFile>")
		return
	}
	binFile := os.Args[1]
	resultFile := os.Args[2]

	vm := &VirtualMachine{}
	err := vm.Execute(binFile)
	if err != nil {
		fmt.Println("Ошибка выполнения:", err)
		return
	}

	// Сохранить результат в YAML
	result := map[string]interface{}{"memory": vm.Memory}
	resultData, err := yaml.Marshal(result)
	if err != nil {
		fmt.Println("Ошибка сохранения результата:", err)
		return
	}

	err = ioutil.WriteFile(resultFile, resultData, 0644)
	if err != nil {
		fmt.Println("Ошибка записи файла:", err)
	}
}
