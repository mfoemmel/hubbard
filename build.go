package hubbard

import "bufio"
import "os"
import "path"

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
	
	logDir := path.Join(cwd, "data", "build", project)
	err = os.MkdirAll(logDir, 0777)
	if err != nil {
		panic(err)
	}

	//_, err := os.Open(path.Join(logDir, sha1 + ".log"), os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0666)
	//if err != nil {
	//	panic(err)
	//}

	findRepo(dir).update(sha1)

	packageDir := path.Join(cwd, "data", "packages", project)
	err = os.MkdirAll(packageDir, 0777)
	if err != nil {
		panic(err)
	}

  println("packaging " + dir + " as " + path.Join(packageDir, sha1 + ".tar.gz"))
  err = archive(dir, path.Join(packageDir, sha1 + ".tar.gz"))
  println("Finished packaging!")
	if err != nil {
		panic(err)
	}

/*
	r := bufio.NewReader(cmd.Stdout)
	for {
		line, err := r.ReadBytes('\n')
		if len(line) != 0 {
			os.Stdout.Write(line)
			log.Write(line)
		}
		if err == os.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
	}
*/
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
