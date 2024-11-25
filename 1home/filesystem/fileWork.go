package filesystem

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type VirtualFileSystem struct {
	rootPath string
	current  string
}

func NewVFS(zipPath string) (*VirtualFileSystem, error) {
	tempDir, err := os.MkdirTemp("", "vfs-*")
	if err != nil {
		return nil, err
	}

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	for _, file := range reader.File {
		destPath := filepath.Join(tempDir, file.Name)

		if file.FileInfo().IsDir() {
			err := os.MkdirAll(destPath, os.ModePerm)
			if err != nil {
				return nil, err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return nil, err
		}

		outFile, err := os.Create(destPath)
		if err != nil {
			return nil, err
		}

		rc, err := file.Open()
		if err != nil {
			err := outFile.Close()
			if err != nil {
				return nil, err
			}
			return nil, err
		}

		_, err = io.Copy(outFile, rc)
		err = outFile.Close()
		if err != nil {
			return nil, err
		}
		err = rc.Close()
		if err != nil {
			return nil, err
		}

	}

	return &VirtualFileSystem{
		rootPath: tempDir,
		current:  tempDir,
	}, nil
}

func (vfs *VirtualFileSystem) ChangeDirectory(path string) error {
	var targetPath string
	if filepath.IsAbs(path) {
		targetPath = filepath.Clean(path)
	} else {
		targetPath = filepath.Clean(filepath.Join(vfs.current, path))
	}

	// Prevent directory traversal outside root
	if !strings.HasPrefix(targetPath, vfs.rootPath) {
		return fmt.Errorf("access denied")
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return errors.New("not a directory")
	}

	vfs.current = targetPath
	return nil
}

func (vfs *VirtualFileSystem) ListDirectory() ([]os.FileInfo, error) {
	files, err := os.ReadDir(vfs.current)
	if err != nil {
		return nil, err
	}

	var fileInfos []os.FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		fileInfos = append(fileInfos, info)
	}

	return fileInfos, nil
}

func (vfs *VirtualFileSystem) CopyFile(src, dst string) error {
	srcPath := filepath.Join(vfs.current, src)
	dstPath := filepath.Join(vfs.current, dst)

	input, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}

	return os.WriteFile(dstPath, input, 0644)
}

func (vfs *VirtualFileSystem) Touch(filename string) error {
	path := filepath.Join(vfs.current, filename)
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}
