package hubbard

import "path"

type project struct {
	name string
}

func (self *project) getLogFile(sha1 string) string {
	return path.Join("data", "build", self.name, sha1+".log")

}
