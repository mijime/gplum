package main

import (
	"golang.org/x/tools/go/vcs"
	"os"
	"path/filepath"
)

type Plum struct {
	repos string
	plugs map[string]*Plug
}

type Plug struct {
	dir  string
	repo string
	at   string
	do   string
	on   []string
	in   []string
	of   []string
}

func (p *Plug) Sync() error {
	if p.isExists() {
		return p.update()
	}

	return p.install()
}

func (p *Plug) Check() {
}

func (p *Plug) Remove() {
}

func (p *Plug) getPath(path string) string {
	dir, _ := filepath.Abs(filepath.Join(p.dir, path))
	return dir
}

func (p *Plug) fetchRoot() (*vcs.RepoRoot, error) {
	return vcs.RepoRootForImportPath(p.repo, false)
}

func (p *Plug) isExists() bool {
	dir := p.getPath("")
	_, err := os.Stat(dir)

	if err != nil {
		return false
	}

	return true
}

func (p *Plug) install() error {
	r, err := p.fetchRoot()

	if err != nil {
		return err
	}

	dir := p.getPath("")
	parentDir := p.getPath("..")

	err = os.Mkdir(parentDir, 0755)

	if err != nil {
		return err
	}

	return r.VCS.CreateAtRev(dir, r.Repo, p.at)
}

func (p *Plug) update() error {
	r, err := p.fetchRoot()

	if err != nil {
		return err
	}

	dir := p.getPath("")

	if err := r.VCS.TagSync(dir, p.at); err != nil {
		return err
	}

	return r.VCS.Download(dir)
}
