package main

import (
	"encoding/json"
	"errors"
	"golang.org/x/tools/go/vcs"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	ProgressSyncing int = 0
	ProgressDoing   int = 1
	ProgressSynced  int = 2
	ProgressFailed  int = 3
)

const (
	StatusSynced int = 0
	StatusChange int = 1
	StatusNoSync int = 2
	StatusCached int = 3
	StatusFailed int = 4
)

type PlugManager struct {
	Root   string
	Plugs  map[string]*Plug
	Status map[string]int
}

func NewPlugManager(rootpath string) *PlugManager {
	root, _ := filepath.Abs(rootpath)

	return &PlugManager{
		Root:   root,
		Plugs:  map[string]*Plug{},
		Status: map[string]int{},
	}
}

func NewPlugManagerFromJSON(ctx io.Reader) (*PlugManager, error) {
	var plum PlugManager

	dec := json.NewDecoder(ctx)

	err := dec.Decode(&plum)

	if err != nil {
		return nil, err
	}

	return &plum, nil
}

func (p *PlugManager) ToJSON(ctx io.Writer) error {
	enc := json.NewEncoder(ctx)

	return enc.Encode(p)
}

func (p *PlugManager) Register(plug *Plug) error {
	if plug.Repo == "" {
		return errors.New("Require repo: need repository path, ex) github.com/mijime/gplum")
	}

	if plug.Name == "" {
		plug.Name = plug.Repo
	}

	if plug.Dir == "" {
		plug.Dir, _ = filepath.Abs(filepath.Join(p.Root, plug.Repo))
	}

	p.Plugs[plug.Name] = plug
	p.Status[plug.Name] = StatusNoSync

	return nil
}

type PlugState struct {
	Name     string
	Progress int
	Err      error
	Time     time.Time
}

func (p *PlugManager) Sync(state chan *PlugState) error {
	var wg sync.WaitGroup

	for _, plug := range p.Plugs {
		wg.Add(1)

		go func(plug *Plug) {
			defer wg.Done()

			state <- &PlugState{
				Name:     plug.Name,
				Progress: ProgressSyncing,
				Err:      nil,
				Time:     time.Now(),
			}

			err := plug.Sync()
			if err == nil {
				state <- &PlugState{
					Name:     plug.Name,
					Progress: ProgressDoing,
					Err:      nil,
					Time:     time.Now(),
				}
				// err = plug.Do()
			}

			if err != nil {
				state <- &PlugState{
					Name:     plug.Name,
					Progress: ProgressFailed,
					Err:      err,
					Time:     time.Now(),
				}
				p.Status[plug.Name] = StatusFailed
				return
			}

			state <- &PlugState{
				Name:     plug.Name,
				Progress: ProgressSynced,
				Err:      nil,
				Time:     time.Now(),
			}
			p.Status[plug.Name] = StatusSynced
		}(plug)
	}

	wg.Wait()
	return nil
}

type Plug struct {
	Name string
	Dir  string
	Repo string
	At   string
	Do   string
	On   []string
	In   []string
	Of   []string
}

func (p *Plug) ToJSON(ctx io.Writer) error {
	enc := json.NewEncoder(ctx)

	return enc.Encode(p)
}

func NewPlugFromJSON(ctx io.Reader) (*Plug, error) {
	var plug Plug

	dec := json.NewDecoder(ctx)

	err := dec.Decode(&plug)

	if err != nil {
		return nil, err
	}

	return &plug, nil
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
	dir, _ := filepath.Abs(filepath.Join(p.Dir, path))
	return dir
}

func (p *Plug) fetchRoot() (*vcs.RepoRoot, error) {
	return vcs.RepoRootForImportPath(p.Repo, false)
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

	err = os.MkdirAll(parentDir, 0755)

	if err != nil {
		return err
	}

	return r.VCS.CreateAtRev(dir, r.Repo, p.At)
}

func (p *Plug) update() error {
	r, err := p.fetchRoot()

	if err != nil {
		return err
	}

	dir := p.getPath("")

	err = r.VCS.TagSync(dir, p.At)

	if err != nil {
		return err
	}

	return r.VCS.Download(dir)

}
