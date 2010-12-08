package hubbard

import "bufio"
import "io"
import "os"

// Returns a slice containing each line of text from the file
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
