package main

import _ "golang.org/x/tools/go/vcs"

type Plum struct {
	Path    string
	Plugins map[string]*Plug
}

type Plug struct {
	Name    string
	Repo    string
	Dir     string
	At      string
	On      []*Plug
	Options map[string][]string
}

func (p *Plum) Register()   {}
func (p *Plum) Deregister() {}
func (p *Plum) Prune()      {}
func (p *Plug) Sync()       {}
func (p *Plug) Check()      {}
func (p *Plug) Remove()     {}
