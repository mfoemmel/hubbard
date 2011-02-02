package hubbard

import (
	"http"
	"log"
	"strings"
	"os"
	"path"
	"template"
)

func getDataDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dataDir := path.Join(cwd, "data")
	return dataDir
}

func getLogDir() string {
	return path.Join(getDataDir(), "build", "logs")
}

func getHtmlDir() string {
	dir, _ := path.Split(os.Args[0])
	return path.Clean(path.Join(dir, "www", "html"))
}

func getCssDir() string {
	dir, _ := path.Split(os.Args[0])
	return path.Clean(path.Join(dir, "www", "css"))
}

func getPackageDir() string {
	return path.Join(getDataDir(), "packages")
}

func getReposDir() string {
	return path.Join(getDataDir(), "repos")
}

func getWorkDir() string {
	return path.Join(getDataDir(), "working")
}

func getDepsDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	depsDir := path.Join(cwd, "deps")
	return depsDir
}

func getDepsDirFor(project string) string {
	return path.Join(getDepsDir(), project)
}

func runTemplate(name string, data interface{}, w http.ResponseWriter) {
	fmap := template.FormatterMap{"": template.HTMLFormatter}
	t, err := template.ParseFile(path.Join(getHtmlDir(), name+".html"), fmap)
	if err != nil {
		log.Println(err)
		return
	}
	err = t.Execute(data, w)
	if err != nil {
		log.Println(err)
		return
	}
}

// HttpHandler
func projectList(w http.ResponseWriter, req *http.Request) {
	runTemplate("projects", readProjects(), w)
}

func newProjectHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()
		name := getParameter(req, "project-name")
		repo := getParameter(req, "project-repo")
		println(name, repo)
		msg := http.URLEscape("not implemented")
		http.Redirect(w, req, "/new-project?err="+msg, http.StatusSeeOther)
	} else {
		req.ParseForm()
		var err string
		if len(req.Form["err"]) > 0 {
			err = req.Form["err"][0]
		}
		runTemplate("new-project", err, w)
	}
}

func projectSummary(w http.ResponseWriter, req *http.Request, projectName string) {
	project := newProject(projectName)
	repo := findRepo(path.Join("data", "repos", projectName))
	out := newHtmlWriter(w)
	out.table()
	for c := range repo.log() {
		out.tr()
		{
			out.td()
			{
				out.raw(`<form action="?sha1=` + c.sha1 + `" method="post"><input type="submit" value="build"/></form>`)
			}
			out.end()

			out.td().text(c.timestamp).end()
			out.td().raw(`<a href="` + projectName + "/" + c.sha1 + `">` + c.sha1[0:6] + `</a>`).end()
			out.td().text(c.author).end()
			out.td().text(strings.Join(c.tags, ",")).end()
			out.td().text(c.timestamp).end()

			pkg := "/packages/" + projectName + "/" + c.sha1 + ".tar.gz"
			if exists("data" + pkg) {
				w.Write([]byte(`<td><a href="` + pkg + `">download</a></td>`))
			} else {
				out.td().end()
			}

			if exists(project.getLogFile(c.sha1)) {
				w.Write([]byte(`<td><a href="` + "/logs/" + projectName + "/" + c.sha1 + ".log" + `">log</a></td>`))
			} else {
				out.td().end()
			}
		}
		out.end()

		out.tr()
		{
			out.td().with("colspan", "6").text(c.comment).end()
		}
		out.end()
	}
	out.end()
}

func revisionSummary(w http.ResponseWriter, req *http.Request, projectName string, sha1 string) {
	repo := findRepo(path.Join("data", "repos", projectName))
	out := newHtmlWriter(w)

	comment, ok := repo.logComment(sha1)
	if !ok {
		out.text("error retrieving commit comment")
		return
	}
	out.pre(comment)
	out.raw("<br/>")
	info, ok := repo.readFile(sha1, "package.hub")
	if !ok {
		out.text("error retrieving info file")
		return
	}
	out.pre(info)
}

func projectHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR", r)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	if req.URL.Path == "/favicon.ico" {
		return
	}

	segments := strings.Split(req.URL.Path[1:], "/", -1)
	if segments[0] == "" {
		segments = segments[1:]
	}
	projectName := req.URL.Path[1:]

	if req.Method == "POST" {
		// Build the project
		req.ParseForm()
		sha1 := getParameter(req, "sha1")
		builder := newBuilder(projectName, sha1)
		builder <- buildCmd{}
		c := make(chan []byte, 100)
		builder <- viewCmd{c}
		for line := range c {
			os.Stdout.Write(line)
		}
		return
	}

	switch len(segments) {
	case 0:
		projectList(w, req)
		return
	case 1:
		projectSummary(w, req, segments[0])
		return
	case 2:
		revisionSummary(w, req, segments[0], segments[1])
		return
	}
}

func Run() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "resolve":
			cmdResolve()
		case "retrieve":
			cmdRetrieve()
		}
		os.Exit(0)
	}

	http.HandleFunc("/", projectHandler)
	http.HandleFunc("/new-project", newProjectHandler)
	http.HandleFunc("/resolve", resolveHandler)
	http.Handle("/packages/", http.FileServer(getPackageDir(), "/packages/"))
	http.Handle("/logs/", http.FileServer(getLogDir(), "/logs/"))
	http.Handle("/css/", http.FileServer(getCssDir(), "/css/"))
	println(getCssDir())
	log.Println("Listening on port 4788")
	err := http.ListenAndServe(":4788", nil)
	if err != nil {
		panic(err)
	}
}
