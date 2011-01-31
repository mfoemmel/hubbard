package hubbard

import "path"

type commit struct {
	timestamp string
	sha1      string
	author    string
	comment   string
	tags      []string
}

type repo interface {
	log() <-chan *commit
	logComment(sha1 string) (string, bool)
	readFile(sha1 string, filename string) (string, bool)
	update(sha1 string)
	resolve(ref string) (sha1 string)
}

func findRepo(dir string) repo {
	if exists(path.Join(dir, ".hg")) {
		return &hgRepo{dir}
	} else if exists(path.Join(dir, ".git")) {
		return &gitRepo{dir}
	}

	panic("unknown repository type for directory: " + dir)
}
