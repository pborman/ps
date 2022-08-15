//go:build linux

package ps

import (
	"os"
	"syscall"
	"testing"
)

var mypid = os.Getpid()

func TestProcesses(t *testing.T) {
	procs, err := Processes(false)
	if err != nil {
		t.Fatalf("Processes(false): %v", err)
	}
	pids := map[int]bool{}
	for _, p := range procs {
		if pids[p.ID] {
			t.Fatalf("Process %d repeated", p.ID)
		}
		pids[p.ID] = true
		if p.dir == "" {
			t.Fatalf("Process[%d] does not have dir filled", p.ID)
		}
		if p.path != "" {
			t.Fatalf("Process[%d] has path filled", p.ID)
		}
		if p.comm != "" {
			t.Fatalf("Process[%d] has comm filled", p.ID)
		}
		if p.stat != nil {
			t.Fatalf("Process[%d] has stat filled", p.ID)
		}
		if p.sysstat != nil {
			t.Fatalf("Process[%d] has sysstat filled", p.ID)
		}
		if p.groups != nil {
			t.Fatalf("Process[%d] has groups filled", p.ID)
		}

	}
	if !pids[mypid] {
		t.Fatalf("My PID was not found")
	}
	procs, err = Processes(true)
	if err != nil {
		t.Fatalf("Processes(false): %v", err)
	}

	pids = map[int]bool{}
	for _, p := range procs {
		if pids[p.ID] {
			t.Fatalf("Process %d repeated", p.ID)
		}
		pids[p.ID] = true
		if p.dir == "" {
			t.Fatalf("Process[%d] dost not have dir filled", p.ID)
		}
		if p.path != "" {
			t.Fatalf("Process[%d] has path filled", p.ID)
		}
		if p.comm != "" {
			t.Fatalf("Process[%d] has comm filled", p.ID)
		}
		if p.stat != nil {
			t.Fatalf("Process[%d] has stat filled", p.ID)
		}
		if p.sysstat == nil {
			t.Fatalf("Process[%d] does not have sysstat filled", p.ID)
		}
		if p.groups != nil {
			t.Fatalf("Process[%d] has groups filled", p.ID)
		}
	}
	if !pids[mypid] {
		t.Fatalf("My PID was not found")
	}

	for _, p := range procs {
		if _, err := p.Stat(); err != nil {
			t.Fatal(err)
		}
		if mypid == int(p.stat.Pid) {
			if p.stat.State != 'R' {
				t.Errorf("I am not running")
			}
			return
		}
	}
	t.Fatalf("Could not find myself")
}

func TestProcessByPid(t *testing.T) {
	p, err := ProcessByPid(mypid)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := p.Stat(); err != nil {
		t.Fatal(err)
	}
	if mypid != int(p.stat.Pid) {
		t.Fatalf("Got process %d, want %d", p.stat.Pid, mypid)
	}
}

func TestClean(t *testing.T) {
	p := &Process{
		ID:      1,
		dir:     "/foo",
		path:    "/bar",
		comm:    "foo",
		stat:    &Stat{},
		sysstat: &syscall.Stat_t{},
		groups:  []int{1},
		status:  map[string]StatusValue{},
	}
	p.Clean()
	if p.path != "" {
		t.Errorf("path not cleared")
	}
	if p.stat != nil {
		t.Errorf("stat not cleared")
	}
	if p.sysstat != nil {
		t.Errorf("sysstat not cleared")
	}
	if p.comm != "" {
		t.Errorf("comm not cleared")
	}
	if p.groups != nil {
		t.Errorf("groups not cleared")
	}
	if p.status != nil {
		t.Errorf("status not cleared")
	}
}

const dev003 = 0x1203

func initDev() {
        devMutex.Lock()
        devNames = map[DevT]string{
                dev003: "dev003",
                noDev:  "-",
        }
        devMutex.Unlock()
}

func TestTty(t *testing.T) {
        initDev()
        p := &Process{
                ID:    mypid,
                stat: &Stat{},
        }
        p.stat.TtyNr = 0x1203
        got, err := p.Tty()
        if err != nil {
                t.Error(err)
        }
        want := "dev003"
        if got != want {
                t.Errorf("Got tty %q, want %q", got, want)
        }
}

func TestEPerm(t *testing.T) {
        p := Process{
                ID: 1,
        }
        _, err := p.Path()
        if err != syscall.EPERM {
                t.Errorf("Got %v, want %v", err, syscall.EPERM)
        }
}
