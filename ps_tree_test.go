package ps

import (
	"fmt"
	"testing"
)

func TestProcessMap(t *testing.T) {
	pm, err := GetProcessMap()
	if err != nil {
		t.Fatal(err)
	}
	if pm == nil {
		t.Fatal("GetProcessMap returned nil map")
	}
	me := pm.Pids[mypid]
	if me == nil {
		t.Fatal("self not in process map")
	}
	ppid, err := me.Ppid()
	if err != nil {
		t.Fatal(err)
	}
	if ppid == 0 {
		t.Fatal("self has no parent")
	}
	found := false
	for _, pid := range pm.GetChildren(ppid) {
		if pid == me.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("self not in parent's children")
	}
	found = false
	for _, pid := range pm.GetDescendants(ppid) {
		if pid == me.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("self not in parent's descendants")
	}
	if got, want := len(pm.GetDecendents(ppid)), len(pm.GetDescendants(ppid)); got != want {
		t.Errorf("GetDecendents returned %d pids, GetDescendants returned %d", got, want)
	}
	if GetChildren(-1) != nil && len(GetChildren(-1)) != 0 {
		t.Errorf("GetChildren of missing pid should be empty")
	}
}

func TestProcessMapCycle(t *testing.T) {
	pm := &ProcessMap{
		Pids: map[int]*Process{
			1: {ID: 1},
			2: {ID: 2},
		},
		Children: map[int][]int{
			1: {2},
			2: {1},
		},
	}
	d := pm.GetDescendants(1)
	if len(d) != 1 || d[0] != 2 {
		t.Errorf("GetDescendants with cycle got %v, want [2]", d)
	}
}

func PrintProcess(p *Process, prefix string) {
	if p == nil {
		return
	}
	command, _ := p.Command()
	argv, _ := p.Argv()
	fmt.Printf("%s%d %s %q\n", prefix, p.ID, command, argv)
	for _, child := range p.Children {
		PrintProcess(child, prefix+"  ")
	}
}
