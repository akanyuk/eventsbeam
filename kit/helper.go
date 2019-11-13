package kit

import (
	"os"
	"path/filepath"
)

func ExecutablePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}
