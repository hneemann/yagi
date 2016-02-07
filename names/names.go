package names

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func createOutNameInt(out, tem string) string {
	if out != "" {
		return out
	}

	name := path.Base(tem)
	if name == tem {
		return "gen-" + name
	}

	return name
}

func CreateOutName(out, tem, message string) (string, error) {
	name := createOutNameInt(out, tem)

	if _, err := os.Stat(name); err == nil {
		// file exists, open it
		f, err := os.Open(name)
		if err != nil {
			return "", fmt.Errorf("can not open file %v, got error: %v", name, err)
		}
		defer f.Close()

		// file read the file header
		header := make([]byte, len(message))
		_, err = f.Read(header)
		if err != nil {
			return "", fmt.Errorf("can not read from file %v, got error: %v", name, err)
		}

		if string(header) != message {
			return "", errors.New("can not overwrite file " + name + ": It seems not to be created by yagi!")
		}
	}

	return name, nil
}

func GetPackageName(pac, out string) (string, error) {
	if pac != "" {
		return pac, nil
	}

	abs, err := filepath.Abs(out)
	if err != nil {
		return "", fmt.Errorf("can not create absolute path of %v, got error: %v", abs, err)
	}

	return path.Base(path.Dir(abs)), nil
}
