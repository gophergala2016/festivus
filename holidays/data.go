package holidays

import (
	"bufio"
	"bytes"
	"io/ioutil"
)

// parseFile returns all holidays from file.
func parseFile(path string) ([]Hday, error) {
	// load file content
	c, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var data []Hday
	s := bufio.NewScanner(bytes.NewReader(c))
	for s.Scan() {
		h, err := NewHday(s.Text())
		if err != nil {
			return nil, err
		}
		data = append(data, h)
	}
	return data, s.Err()
}
