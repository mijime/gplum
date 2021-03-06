package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPlugSync(t *testing.T) {
	plug := &Plug{
		Dir:  "_test/sham",
		Repo: "github.com/mijime/sham",
	}

	err := plug.Sync()

	if err != nil {
		t.Error(err)
		return
	}
}

func TestNewPlugFromJSON(t *testing.T) {
	stdin := bytes.NewBufferString(`{"repo":"github.com/mijime/sham"}`)
	stdout := new(bytes.Buffer)

	plug, err := NewPlugFromJSON(stdin)

	if err != nil {
		t.Error(err)
		return
	}

	plug.ToJSON(stdout)
}

func TestPlugManagerRegister(t *testing.T) {
	stdout := new(bytes.Buffer)

	plum := NewPlugManager("_test")
	plug := &Plug{Repo: "github.com/mijime/sham"}

	plum.Register(plug)

	plum.ToJSON(stdout)
	fmt.Println(stdout)
}

func TestPlugManagerSync(t *testing.T) {
	stdout := new(bytes.Buffer)

	s := make(chan *PlugState)
	closed := make(chan bool)

	go func() {
		for {
			select {
			case state := <-s:
				if state != nil {
					fmt.Println(state)
				}

			case <-closed:
				return
			}
		}
	}()

	plum := NewPlugManager("_test")

	plum.Register(&Plug{Repo: "github.com/mijime/sham"})
	plum.Register(&Plug{Repo: "github.com/mijime/merje"})

	plum.Sync(s)
	closed <- true

	close(s)
	close(closed)

	plum.ToJSON(stdout)
	fmt.Println(stdout)
}
