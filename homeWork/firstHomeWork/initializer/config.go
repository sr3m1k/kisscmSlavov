package initializer

import (
	"errors"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
)

func FillConfig(filePath string, conf any) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer func() {
		if fileErr := file.Close(); fileErr != nil {
			fmt.Print("Error closing the file")
			if err != nil {
				err = errors.Join(err, fileErr)
			} else {
				err = fileErr
			}
		}
	}()
	decoder := toml.NewDecoder(file)
	if err := decoder.Decode(conf); err != nil {
		return fmt.Errorf("failed to decode file: %v", err)
	}

	return err
}
