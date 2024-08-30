package ps

import (
	"fmt"
	"testing"
)

func ExampleProcessMap(t *testing.T) {
	pm := GetProcessMap()
	PrintProcess(pm.Pids[1], "")
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
