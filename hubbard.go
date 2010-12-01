package hubbard

import "log"
import "html"
import "http"
import "strings"
import "os"

func projectHandler(w http.ResponseWriter, req *http.Request) {
	// todo extract project name from URL
	projectName := "go"

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
			
			pkg := "data/packages/" + projectName + "/" + c.sha1 + ".tar.gz"
			if exists(pkg) {
				w.Write([]byte(`<td><a href="` + pkg + `">download</a></td>`))
			} else {
				out.td().end()
			}
		}
		out.end()

		out.tr()
		{
			out.td().with("colspan", "6").text(html.EscapeString(c.comment)).end()
		}
		out.end()
	}
	out.end()
}

func Run() {
	http.HandleFunc("/", projectHandler)
	log.Println("Listening on port 4788")
	err := http.ListenAndServe(":4788", nil)
	if err != nil {
		panic(err)
	}
}