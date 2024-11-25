package gui

import (
	"fmt"
	"log"
	"strings"

	"1home/config"
	"1home/filesystem"
	"1home/logger"
	"github.com/andlabs/ui"
)

type ShellEmulator struct {
	config     *config.ShellConfig
	vfs        *filesystem.VirtualFileSystem
	logger     *logger.Logger
	mainWindow *ui.Window
	output     *ui.MultilineEntry
	input      *ui.Entry
}

func NewShellEmulator(configPath string) (*ShellEmulator, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	vfs, err := filesystem.NewVFS(cfg.VFSArchive)
	if err != nil {
		return nil, err
	}

	newLogger := logger.NewLogger(cfg.LogFile)

	return &ShellEmulator{
		config: cfg,
		vfs:    vfs,
		logger: newLogger,
	}, nil
}

func (s *ShellEmulator) initGUI() {
	ui.OnShouldQuit(func() bool {
		s.mainWindow.Destroy()
		return true
	})

	s.mainWindow = ui.NewWindow(
		fmt.Sprintf("Shell Emulator - %s", s.config.ComputerName),
		800, 600, false,
	)
	s.mainWindow.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})

	s.output = ui.NewMultilineEntry()
	s.output.SetReadOnly(true)

	s.input = ui.NewEntry()
	s.input.OnChanged(func(*ui.Entry) {
		s.processCommand(s.input.Text())
	})

	box := ui.NewVerticalBox()
	box.Append(s.output, true)
	box.Append(s.input, false)

	s.mainWindow.SetChild(box)
	s.mainWindow.SetMargined(true)
	s.mainWindow.Show()
}

func (s *ShellEmulator) processCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	err := s.logger.Log(s.config.Username, parts[0], parts[1:])
	if err != nil {
		return
	}

	switch parts[0] {
	case "ls":
		s.handleLS()
	case "cd":
		s.handleCD(parts)
	case "cp":
		s.handleCP(parts)
	case "touch":
		s.handleTouch(parts)
	case "whoami":
		s.handleWhoami()
	case "exit":
		ui.Quit()
	default:
		s.output.Append(fmt.Sprintf("Unknown command: %s\n", parts[0]))
	}
}

func (s *ShellEmulator) handleLS() {
	files, err := s.vfs.ListDirectory()
	if err != nil {
		s.output.Append(fmt.Sprintf("Error: %v\n", err))
		return
	}
	for _, file := range files {
		s.output.Append(fmt.Sprintf("%s\n", file.Name()))
	}
}

func (s *ShellEmulator) handleCD(parts []string) {
	if len(parts) < 2 {
		s.output.Append("cd requires a path\n")
		return
	}
	if err := s.vfs.ChangeDirectory(parts[1]); err != nil {
		s.output.Append(fmt.Sprintf("Error: %v\n", err))
	}
}

func (s *ShellEmulator) handleCP(parts []string) {
	if len(parts) < 3 {
		s.output.Append("cp requires source and destination\n")
		return
	}
	if err := s.vfs.CopyFile(parts[1], parts[2]); err != nil {
		s.output.Append(fmt.Sprintf("Error: %v\n", err))
	}
}

func (s *ShellEmulator) handleTouch(parts []string) {
	if len(parts) < 2 {
		s.output.Append("touch requires a filename\n")
		return
	}
	if err := s.vfs.Touch(parts[1]); err != nil {
		s.output.Append(fmt.Sprintf("Error: %v\n", err))
	}
}

func (s *ShellEmulator) handleWhoami() {
	s.output.Append(fmt.Sprintf("%s\n", s.config.Username))
}

func (s *ShellEmulator) Run() {
	err := ui.Main(func() {
		s.initGUI()
	})
	if err != nil {
		log.Fatal(err)
	}
}
