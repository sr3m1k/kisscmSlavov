package main

import (
	"log"
	"os"
	"path/filepath"

	"1home/gui"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide the path to the configuration file")
	}

	configPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	emulator, err := gui.NewShellEmulator(configPath)
	if err != nil {
		log.Fatal(err)
	}

	emulator.Run()
}
