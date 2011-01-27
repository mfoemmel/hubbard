package hubbard

import (
	"os"
	"path"
)

type project struct {
	name string
}

func readProjects() []*project {
	fd, err := os.Open(getReposDir(), os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	names, err := fd.Readdirnames(-1)
	if err != nil {
		panic(err)
	}
	projects := make([]*project, len(names))
	for i, name := range names {
		projects[i] = &project{name}
	}
	return projects
}

func (self *project) getLogFile(sha1 string) string {
	return path.Join("data", "build", self.name, sha1+".log")
}

