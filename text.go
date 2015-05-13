package gs

import (
	"bufio"
	"os"
)

func Lines(file string) (lines []string) {
	f, err := os.Open(file)
	c(err)
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	c(s.Err())
	return
}
