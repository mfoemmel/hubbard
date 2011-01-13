package hubbard

import "bufio"
import "exec"
import "os"
import "strings"

var hg string

func init() {
	var err os.Error
	hg, err = exec.LookPath("hg")
	if err != nil {
		panic("hg not found on PATH")
	}
}

type hgRepo struct {
	dir string
}

func (self *hgRepo) log() <-chan *commit {
	c := make(chan *commit, 100)
	go func() {
		defer close(c)
		argv := []string{hg, "log", "--template", `{date|isodate}\t{node}\t{author}\t{tags}\t{desc|firstline}\n`, "-b", "default"}
		cmd, err := exec.Run(hg, argv, nil, self.dir, exec.DevNull, exec.Pipe, exec.PassThrough)
		if err != nil {
			panic(err)
		}
		r := bufio.NewReader(cmd.Stdout)
		for {
			line, err := r.ReadString('\n')
			if line != "" {
				fields := strings.Split(string(line), "\t", -1)
				tags := strings.Split(fields[3], " ", -1)
				c <- &commit{fields[0], fields[1], fields[2], fields[4], tags}
			}
			if err == os.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
		}
		msg, err := cmd.Wait(0)
		if err != nil {
			panic(err)
		}
		if msg.WaitStatus != 0 {
			panic("process failed")
		}
	}()
	return c
}

func (self *hgRepo) logComment(sha1 string) (string, bool) {
	return capture(self.dir, []string{"hg", "log", "--template", `{desc}`, "--rev", sha1})
}

func (self *hgRepo) readFile(sha1 string, filename string) (string, bool) {
	return "error", false
}

func (self *hgRepo) update(sha1 string) {
	run(os.Stdout, nil, self.dir, []string{findExe("hg"), "update", "-C", sha1})
}

func (self *hgRepo) resolve(ref string) (sha1 string) {
	branches, err := captureLines(self.dir, []string{"hg", "branches", "--debug"})
	if err != nil {
		panic(err)
	}

	tags, err := captureLines(self.dir, []string{"hg", "tags", "--debug"})
	if err != nil {
		panic(err)
	}

	for _, line := range append(branches, tags...) {
		fields := strings.Split(line, " ", 2)
		if strings.TrimSpace(fields[0]) == ref {
			return strings.Split(strings.TrimSpace(fields[1]), ":", 2)[1]
		}
	}
	panic("ref not found: " + ref)
}
