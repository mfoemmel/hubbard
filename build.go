package hubbard

import "bufio"
import "os"

type buildCmd struct {
}

type viewCmd struct {
	out chan<- []byte
}

type logCmd struct {
	line []byte
}

func newBuilder(project string, sha1 string) chan<- interface{} {
	viewers := [](chan<- []byte){}
	c := make(chan interface{}, 100)
	go func() {
		for msg := range c {
			switch cmd := msg.(type) {
			case buildCmd:
				go func() {
					build(project, sha1, c)
				}()
			case viewCmd:
				println("here")
				viewers = append(viewers, cmd.out)
			case logCmd:
				for _, viewer := range viewers {
					viewer <- cmd.line
				}
			}
		}
	}()
	return c
}

func build(project string, sha1 string, builder chan<- interface{}) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir := cwd + "/data/working/" + project
	
	run(os.Stdout, nil, dir, []string { findExe("hg"), "update", "-C", sha1 })
	if buildExec(dir + "/src", []string { "/bin/bash", dir + "/src/all.bash" }, builder) {
		run(os.Stdout, nil, dir, []string { findExe("tar"), "-cvf", cwd + "/data/packages/" + project + "/" + sha1 + ".tar.gz", "--exclude", ".hg", "." })
	}
}

func buildExec(dir string, argv []string, builder chan <- interface{}) bool {
	pipeIn, pipeOut := pipe()

	pid, err := os.ForkExec(argv[0], argv, os.Environ(), dir, []*os.File{ os.Stdin, pipeOut, pipeOut })
	if err != nil {
		panic(err)
	}
	err = pipeOut.Close()

	in := bufio.NewReader(pipeIn)
	for {
		line, err := in.ReadBytes('\n')
		if err == os.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		builder <- logCmd{line}
	}

	msg, err := os.Wait(pid, 0)
	if err != nil {
		panic(err)
	}

	return msg.WaitStatus == 0 
}