//go:build linux

package ps

import (
	"os"
	"testing"
)

var mypid = os.Getpid()

func TestSomething(t *testing.T) {
	p, err := ProcessByPid(28539)
	if err != nil {
		t.Fatal(err)
	}
	t.Error(p.Path())
	t.Error(p.Command())
	t.Error(p.Argv())
	t.Error(p.Value("HOME"))
	t.Error(p.Environ())
	t.Error(p.Groups())
	t.Error(p.Tty())
	t.Error(p.Tty())
	t.Error("Trigger")
	t.Error(p.Footprint())
	t.Error(p.Footprint(false))
	t.Error(p.Footprint(true))
}
