package shellProgram

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Shell struct {
	Config      Config
	CurrentPath string
	LogFile     *os.File
}

// Run Функция для запуска эмулятора shellRunner
func (s *Shell) Run() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s@%s:%s$ ", s.Config.User, s.Config.Hostname, s.CurrentPath)
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)
		parts := strings.Split(command, " ")

		switch parts[0] {
		case "ls":
			s.ls()
		case "cd":
			if len(parts) > 1 {
				s.cd(parts[1])
			} else {
				fmt.Println("Missing argument")
			}
		case "cp":
			if len(parts) == 3 {
				s.cp(parts[1], parts[2])
			} else {
				fmt.Println("Usage: cp [src] [dst]")
			}
		case "whoami":
			s.whoami()
		case "touch":
			if len(parts) > 1 {
				s.touch(parts[1])
			} else {
				fmt.Println("Missing argument")
			}
		case "exit":
			s.exit()
			return
		default:
			fmt.Println("Unknown command")
		}
	}
}
