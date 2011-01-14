package hubbard

import "bytes"
import "fmt"
import "http"
import "io"
import "io/ioutil"
import "log"
import "os"
import "path"
import "strings"

// HttpHandler
// This handler returns the SHA1 associated with a particular tag or branch name for a project
func resolveHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR:", r)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v\n", r)
		}
	}()

	w.SetHeader("Content-Type", "text/plain")
	req.ParseForm()
	projectName := getParameter(req, "project")
	ref := getParameter(req, "ref")
	repo := findRepo(path.Join(getReposDir(), projectName))
	io.WriteString(w, repo.resolve(ref)+"\n")
}

// Client side logic for retrieving SHA1s  
func cmdResolve() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", r)
			os.Exit(1)
		}
	}()

	lines, err := readFileLines("package.hub")
	if err != nil {
		panic(err)
	}

	// Flag to indicate whether we're in the [versions] section or not
	inVersions := false

	buf := bytes.NewBuffer(nil)

	// Loop through lines in package.hub
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			buf.Write([]byte(line + "\n"))
			continue
		}
		if strings.HasPrefix(line, "[") {
			inVersions = line == "[versions]"
			buf.Write([]byte(line + "\n"))
		} else if inVersions {
			fields := strings.Split(strings.Split(line, "=", 2)[0], "/", 2)
			project := fields[0]
			ref := fields[1]
			sha1, err := resolve(project, ref)
			if err != nil {
				panic(err)
			}
			buf.Write([]byte(project + "/" + ref + "=" + sha1 + "\n"))
		} else {
			buf.Write([]byte(line + "\n"))
		}
	}

	err = ioutil.WriteFile("package.hub", buf.Bytes(), 0666)
	if err != nil {
		panic(err)
	}
}

func resolve(project string, ref string) (string, os.Error) {
	// URLUnescape?
	project = http.URLEscape(project)
	ref = http.URLEscape(ref)
	resp, _, err := http.Get("http://localhost:4788/resolve?project=" + project + "&ref=" + ref)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		panic(strings.TrimSpace(string(body)))
	}
	return strings.TrimSpace(string(body)), nil
}
