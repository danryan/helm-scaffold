package os

import (
	"errors"
	"fmt"

	"github.com/spf13/afero"
)

func ReadFile(fs afero.Fs, filename string) (string, error) {
	if filename == "" {
		return "", errors.New("readFile needs a filename")
	}

	if info, err := fs.Stat(filename); err == nil {
		if info.Size() > 1000000 {
			return "", fmt.Errorf("File %q is too big", filename)
		}
	} else {
		return "", err
	}

	b, err := afero.ReadFile(fs, filename)

	if err != nil {
		return "", err
	}

	return string(b), nil
}
