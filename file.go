package gs

import (
	"os"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func SaveFile(file, s string) {
	f, err := os.Create(file)
	c(err)
	defer f.Close()
	f.Write([]byte(s))
}
