package hubbard

import "log"
import "http"
import "strings"
import "os"

func projectList(w http.ResponseWriter, req *http.Request) {
	fd, err := os.Open("data/repos", os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	projects, err := fd.Readdirnames(-1)
	if err != nil {
		panic(err)
	}
	out := newHtmlWriter(w)
	for _, project := range projects {
		out.raw("<a href=\"" + project + "\">")
		out.text(project)
		out.raw("</a>")
		out.raw("<br/>")
	}
}

func projectHandler(w http.ResponseWriter, req *http.Request) {
	projectName := req.URL.Path[1:]
	if projectName == "favicon.ico" {
		return
	}

	if projectName == "" {
		projectList(w, req)
		return
	}

	if req.Method == "POST" {
		// Build the project
		req.ParseForm()
		sha1 := getParameter(req, "sha1")
		builder := newBuilder("go", sha1)
		builder <- buildCmd{}
		c := make(chan []byte, 100)
		builder <- viewCmd{c}
		for line := range c {
			os.Stdout.Write(line)
		}
		return
	}

	repo := findRepo(projectName)
	out := newHtmlWriter(w)
	out.table()
	for c := range repo.log() {
		out.tr()
		{
			out.td()
			{
				out.raw(`<form action="build?sha1=` + c. sha1 + `" method="post"><input type="submit" value="build"/></form>`)
			}
			out.end()

			out.td().text(c.timestamp).end()
			out.td().text(c.sha1[0:6]).end()
			out.td().text(c.author).end()
			out.td().text(strings.Join(c.tags, ",")).end()
			out.td().text(c.timestamp).end()
			
			pkg := "/packages/" + projectName + "/" + c.sha1 + ".tar.gz"
			if exists("data" + pkg) {
				w.Write([]byte(`<td><a href="` + pkg + `">download</a></td>`))
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

func Run() {
	http.HandleFunc("/", projectHandler)
	http.Handle("/packages/", http.FileServer("data/packages", "/packages/"))
	log.Println("Listening on port 4788")
	err := http.ListenAndServe(":4788", nil)
	if err != nil {
		panic(err)
	}
}
