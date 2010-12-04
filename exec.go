package hubbard

import "bytes"
import "io"
import "fmt"
import "os"
import "path"
import "strings"

func capture(dir string, argv []string) (string, bool) {
	argv[0] = findExe(argv[0])
	buf := bytes.NewBuffer(nil)
	if run(buf, nil, dir, argv) {
		return buf.String(), true
	}
	return "", false
}

func run(w io.Writer, env []string, dir string, argv []string) bool {
	if !exists(dir) {
		panic("directory doesn't exist: " + dir)
	}

	pipeIn, pipeOut := pipe()

	if env == nil {
		env = os.Environ()
	}

	fmt.Printf("%s -- %v\n", dir, argv)

	pid, err := os.ForkExec(argv[0], argv, os.Environ(), dir, []*os.File{ os.Stdin, pipeOut, os.Stderr })
	if err != nil {
		panic(err)
	}
	err = pipeOut.Close()

	if err != nil {
		panic(err)
	}

	_, err = io.Copy(w, pipeIn)
	if err != nil {
		panic(err)
	}
	
	msg, err := os.Wait(pid, 0)
	if err != nil {
		panic(err)
	}

	return msg.WaitStatus == 0 
}

func pipe() (*os.File, *os.File) {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	return r,w
}

func findExe(exe string) string {
	for _, dir := range strings.Split(os.Getenv("PATH"), ":", -1) {
		file := path.Join(dir,exe)
		if exists(file) {
			return file
		}
	}
	panic("The executable '" + exe + "' not found in PATH")
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if pathErr, ok := err.(*os.PathError); ok && pathErr.Error == os.ENOENT {
		return false
	}
	panic(err)
}

