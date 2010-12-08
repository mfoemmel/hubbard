package hubbard

import "bufio"
import "io"
import "os"

// Returns a slice containing each line of text from the specified file
func readFileLines(filename string) ([]string, os.Error) {
	file, err := os.Open(filename, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return readLines(file)
}

// Returns a slice containing each line of text from the input
func readLines(r io.Reader) ([]string, os.Error) {
	lines := make([]string, 0)
	in := bufio.NewReader(r)
	for {
		line, err := in.ReadString('\n')
		if line != "" {
			lines = append(lines, line)
		}
		if err == os.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return lines, nil
}
