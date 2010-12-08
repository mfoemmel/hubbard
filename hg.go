package hubbard

import "bufio"
import "io"
import "os"
import "path"
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
	return capture(self.dir, []string{ "hg", "log", "--template", `{desc}`, "--rev", sha1})
}

func (self *hgRepo) readFile(sha1 string, filename string) (string, bool) {
	return "error", false
}

func (self *hgRepo) update(sha1 string) {
	run(os.Stdout, nil, self.dir, []string { findExe("hg"), "update", "-C", sha1 })
}

func (self *hgRepo) resolve(ref string) (sha1 string) {
	file, err := os.Open(path.Join(self.dir, ".hgtags"), os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	in := bufio.NewReader(file)
	for {
		line, err := in.ReadString('\n')
		if err != nil {
			panic(err)
		}
		fields := strings.Split(line, " ", 2)
		if strings.TrimSpace(fields[1]) == ref {
			return fields[0]
		}
	}
	panic("unreachable")
}
