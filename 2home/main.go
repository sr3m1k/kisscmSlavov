package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Config структура для парсинга XML
type Config struct {
	GraphVisualizerPath string `xml:"graphVisualizerPath"`
	RepositoryPath      string `xml:"repositoryPath"`
	OutputFile          string `xml:"outputFile"`
	BranchName          string `xml:"branchName"`
}

// Commit структура для представления коммита
type Commit struct {
	Hash   string
	Parent string
}

// Чтение XML конфигурации
func readConfig(configPath string) (Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}
	var config Config
	if err := xml.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

// Получение списка коммитов с помощью git log
func getCommits(repoPath, branchName string) ([]Commit, error) {
	cmd := exec.Command("git", "-C", repoPath, "log", branchName, "--pretty=format:%H %P")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var commits []Commit
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) > 0 {
			hash := parts[0]
			parent := ""
			if len(parts) > 1 {
				parent = parts[1]
			}
			commits = append(commits, Commit{Hash: hash, Parent: parent})
		}
	}
	return commits, nil
}

// Построение Mermaid графа
func buildMermaidGraph(commits []Commit) string {
	var builder strings.Builder
	builder.WriteString("graph TD\n")

	for _, commit := range commits {
		if commit.Parent != "" {
			builder.WriteString(fmt.Sprintf("    %s --> %s\n", commit.Parent, commit.Hash))
		} else {
			builder.WriteString(fmt.Sprintf("    %s\n", commit.Hash))
		}
	}

	return builder.String()
}

// Запись графа в файл
func writeOutput(filePath, content string) error {
	return ioutil.WriteFile(filePath, []byte(content), 0644)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <config.xml>")
	}

	configPath := os.Args[1]
	config, err := readConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	commits, err := getCommits(config.RepositoryPath, config.BranchName)
	if err != nil {
		log.Fatalf("Failed to get commits: %v", err)
	}

	mermaidGraph := buildMermaidGraph(commits)
	fmt.Println("Mermaid Graph:\n", mermaidGraph)

	if err := writeOutput(config.OutputFile, mermaidGraph); err != nil {
		log.Fatalf("Failed to write output: %v", err)
	}
}
