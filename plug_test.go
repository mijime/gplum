package main

import "testing"

func TestPlugSync(t *testing.T) {
	plug := &Plug{
		dir:  "_test/sham",
		repo: "github.com/mijime/sham",
	}

	err := plug.Sync()

	if err != nil {
		t.Error(err)
		return
	}
}
