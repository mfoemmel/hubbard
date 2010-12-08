package hubbard

import "fmt"
import "http"
import "io"
import "log"
import "path"

// This handler returns the SHA1 associated with a particular tag or branch name for a project
func resolveHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR:", r)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "ERROR: %v\n", r)
		}
	}()

	w.SetHeader("Content-Type", "text/plain")
	req.ParseForm()
	projectName := getParameter(req, "project")
	ref := getParameter(req, "ref")
	repo := findRepo(path.Join("data", "repos", projectName))
	io.WriteString(w, repo.resolve(ref) + "\n")
}
