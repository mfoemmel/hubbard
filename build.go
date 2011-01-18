package hubbard

import "bufio"
import "os"
import "path"
import "strings"

type buildCmd struct{}

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
	dir := path.Join(getWorkDir(), project)

	logDir := path.Join(getLogDir(), project)
	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		panic(err)
	}

	//_, err := os.Open(path.Join(logDir, sha1 + ".log"), os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0666)
	//if err != nil {
	//	panic(err)
	//}

	// Update to the given rev for the project.
	findRepo(dir).update(sha1)

	// Only retrieve dependencies or build the project if there's
	// actually a 'package.hub' file in the project.
	hubFile := path.Join(dir, "package.hub")
	if exists(hubFile) {
		println("DEBUG *** Build processing hubFile: ", hubFile)

		// Retrieve any project dependencies.
		println("DEBUG *** Build running srvRetrieve() for ", project, sha1)
		srvRetrieve(dir)

		// Parse the [build] section of the package.hub file.
		println("DEBUG *** Build parsing package.hub ...")
		workDir, buildCmd := parseBuild(hubFile)
		if err != nil {
			panic(err)
		}
		baseDir := getBaseDir(hubFile)
		workDir = path.Join(baseDir, workDir)
		println("DEBUG *** Building project: ", project)
		println("DEBUG *** Build BaseDir: ", baseDir)
		println("DEBUG *** Build WorkDir: ", workDir)
		println("DEBUG *** Build Cmd: ", buildCmd)
		// Build the project.
		success := buildExec(workDir, buildCmd, builder)
		// TODO: Don't panic the server if a build fails.
		if success != true {
			panic("Build command failed!")
		}
	}

	packageDir := path.Join(cwd, "data", "packages", project)
	err = os.MkdirAll(packageDir, 0777)
	if err != nil {
		panic(err)
	}

	println("packaging " + dir + " as " + path.Join(packageDir, sha1+".tar.gz"))
	err = archive(dir, path.Join(packageDir, sha1+".tar.gz"))
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

func buildExec(dir string, argv []string, builder chan<- interface{}) bool {
	pipeIn, pipeOut := pipe()

	pid, err := os.ForkExec(argv[0], argv, os.Environ(), dir, []*os.File{os.Stdin, pipeOut, pipeOut})
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

// Parses the [build] section for a given hubFile. (usually 'package.hub')
// hubFile is a path to the file.
func parseBuild(hubFile string) (workDir string, buildCmd []string) {
	println("DEBUG *** parseBuild processing: ", hubFile)
	// Defaults for building a project.
	workDir = "."
	buildCmd = []string{"make", "test"}

	lines, err := readFileLines(hubFile)
	if err != nil {
		panic(err)
	}

	// Flag to indicate whether we're in the [build] section or not
	inBuild := false

	// Loop through lines in package.hub
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") {
			inBuild = line == "[build]"
		} else if inBuild {
			components := strings.Split(line, "=", 2)
			switch components[0] {
			case "workdir":
				workDir = components[1]
			case "command":
				// TODO: How many command args should we support?  More than 10?
				buildCmd = strings.Split(components[1], " ", 10)
			}
		}
	}
	return workDir, buildCmd
}

// A project's base directory is the directory containing the 'package.hub' file.
func getBaseDir(hubFile string) string {
	baseDir, hubFile := path.Split(hubFile)
	return baseDir
}
