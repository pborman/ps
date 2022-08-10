//go:build linux

package ps

import (
	"os"
	"testing"
)

var mypid = os.Getpid()

func TestSomething(t *testing.T) {
	p, err := ProcessByPid(32586)
	if err != nil {
		t.Fatal(err)
	}
	t.Error(p.Path())
	t.Error(p.Command())
	t.Error(p.Argv())
	t.Error(p.Value("HOME"))
	t.Error(p.Environ())
	t.Error("Trigger")
}
