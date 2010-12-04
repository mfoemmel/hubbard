package hubbard

import "bufio"
import "io"
import "os"
import "strings"

type gitRepo struct {
	dir string
}

func (self *gitRepo) log() <-chan *commit {
	c := make(chan *commit, 100)
	pipeIn, pipeOut := io.Pipe()
	go func() {
		argv := []string{ findExe("git"), "log", `--format=%ci|%H|%an|%d|%s`}
		run(pipeOut, nil, self.dir, argv)
		pipeOut.Close()
	}()
	go func() {
		defer close(c)
		in := bufio.NewReader(pipeIn)
		for {
			line, err := in.ReadBytes('\n')
			if err == os.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			
			fields := strings.Split(string(line), "|", -1)
			tags := []string{}
			// EXAMPLE
			// (HEAD, v0.0.18, origin/master, origin/HEAD, master)
			if fields[3] != "" {
				// Strip out the leading space and parentheses.
				tags = strings.Split(fields[3][2:len(fields[3])-1], ", ", -1)
			}

			c <- &commit{fields[0], fields[1], fields[2], fields[4], tags}
		}
	}()
	return c
}

func (self *gitRepo) logComment(sha1 string) (string, bool) {
	return capture(self.dir, []string{ "git", "show", "-s", `--format=%B`, sha1})
}

func (self *gitRepo) readFile(sha1 string, filename string) (string, bool) {
	return capture(self.dir, []string{ "git", "show", sha1 + ":" + filename})
}

func (self *gitRepo) update(sha1 string) {
	run(os.Stdout, nil, self.dir, []string { findExe("git"), "checkout", sha1 })
}
