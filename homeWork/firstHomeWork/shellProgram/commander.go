package shellProgram

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Функция для выполнения команды `ls`
func (s *Shell) ls() {
	files, err := ioutil.ReadDir(s.CurrentPath)
	if err != nil {
		fmt.Println("Error:", err)
		s.logAction("ls", err.Error())
		return
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}

	s.logAction("ls", "success")
}

// Функция для выполнения команды `cd`
func (s *Shell) cd(dir string) {
	newPath := filepath.Join(s.currentPath, dir)
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		fmt.Println("Directory does not exist.")
		s.logAction(fmt.Sprintf("cd %s", dir), "Directory does not exist")
		return
	}

	s.currentPath = newPath
	s.logAction(fmt.Sprintf("cd %s", dir), "success")
}

// Функция для выполнения команды `exit`
func (s *Shell) exit() {
	s.logAction("exit", "exit shellRunner")
	fmt.Println("Goodbye!")
}

// Функция для выполнения команды `cp`
func (s *Shell) cp(src, dst string) {
	srcPath := filepath.Join(s.currentPath, src)
	dstPath := filepath.Join(s.currentPath, dst)

	data, err := ioutil.ReadFile(srcPath)
	if err != nil {
		fmt.Println("Error:", err)
		s.logAction(fmt.Sprintf("cp %s %s", src, dst), err.Error())
		return
	}

	err = ioutil.WriteFile(dstPath, data, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		s.logAction(fmt.Sprintf("cp %s %s", src, dst), err.Error())
		return
	}

	s.logAction(fmt.Sprintf("cp %s %s", src, dst), "success")
}

// Функция для выполнения команды `whoami`
func (s *Shell) whoami() {
	fmt.Println(s.config.User)
	s.logAction("whoami", s.config.User)
}

// Функция для выполнения команды `touch`
func (s *Shell) touch(filename string) {
	newFile := filepath.Join(s.currentPath, filename)
	file, err := os.Create(newFile)
	if err != nil {
		fmt.Println("Error:", err)
		s.logAction(fmt.Sprintf("touch %s", filename), err.Error())
		return
	}
	file.Close()
	s.logAction(fmt.Sprintf("touch %s", filename), "success")
}
