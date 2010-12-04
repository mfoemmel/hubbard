package hubbard

import "bufio"
import "io"
import "os"
import "strings"

type hgRepo struct {
	dir string
}

func (self *hgRepo) log() <-chan *commit {
	c := make(chan *commit, 100)
	pipeIn, pipeOut := io.Pipe()
	go func() {
		argv := []string{ findExe("hg"), "log", "--template", `{date|isodate}\t{node}\t{author}\t{tags}\t{desc|firstline}\n`, "-b", "default"} 
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
			
			fields := strings.Split(string(line), "\t", -1)
			tags := strings.Split(fields[3], " ", -1)
			c <- &commit{fields[0], fields[1], fields[2], fields[4], tags}		
		}
	}()
	return c
}

func (self *hgRepo) logComment(sha1 string) (string, bool) {
	panic("not implemented")
}

func (self *hgRepo) readFile(sha1 string, filename string) (string, bool) {
	panic("not implemented")
}

func (self *hgRepo) update(sha1 string) {
	run(os.Stdout, nil, self.dir, []string { findExe("hg"), "update", "-C", sha1 })
}
