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
		if p.cpath != "" {
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
		if p.cgroups != nil {
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
		if p.cpath != "" {
			t.Fatalf("Process[%d] has path filled", p.ID)
		}
		if p.comm != "" {
			t.Fatalf("Process[%d] has comm filled", p.ID)
		}
		if p.stat == nil {
			t.Fatalf("Process[%d] does not have stat filled", p.ID)
		}
		if p.sysstat == nil {
			t.Fatalf("Process[%d] does not have sysstat filled", p.ID)
		}
		if p.cgroups != nil {
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
		cpath:   "/bar",
		comm:    "foo",
		stat:    &Stat{},
		sysstat: &syscall.Stat_t{},
		cgroups: []int{1},
		status:  map[string]StatusValue{},
		cargv:   []string{"x"},
		cenv:    map[string]string{"A": "B"},
	}
	p.Clean()
	if p.cpath != "" {
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
	if p.cgroups != nil {
		t.Errorf("groups not cleared")
	}
	if p.status != nil {
		t.Errorf("status not cleared")
	}
	if p.cargv != nil {
		t.Errorf("argv not cleared")
	}
	if p.cenv != nil {
		t.Errorf("env not cleared")
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
		ID:   mypid,
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
	p.stat.TtyNr = 0
	got, err = p.Tty()
	if err != nil {
		t.Error(err)
	}
	if got != "-" {
		t.Errorf("Got tty %q, want %q", got, "-")
	}
}

func TestDevT(t *testing.T) {
	initDev()
	dev := DevT((0x34567 & 0xff) | (0x12 << 8) | ((0x34567 &^ 0xff) << 12))
	if got, want := dev.Major(), 0x12; got != want {
		t.Errorf("Major() got %02x, want %02x", got, want)
	}
	if got, want := dev.Minor(), 0x34567; got != want {
		t.Errorf("Minor() got %05x, want %05x", got, want)
	}

	dev = DevT(dev003)
	if got, want := dev.String(), "dev003"; got != want {
		t.Errorf("device 0x%08x got name %q, want %q", uint32(dev), got, want)
	}
	dev = noDev
	if got, want := dev.String(), "-"; got != want {
		t.Errorf("device 0x%08x got name %q, want %q", uint32(dev), got, want)
	}
}

func TestParseStatComm(t *testing.T) {
	data := []byte("1234 (foo) bar) S 10 20 30 0 40 41 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0\n")
	s, err := parseStat(data)
	if err != nil {
		t.Fatal(err)
	}
	if s.Pid != 1234 {
		t.Errorf("Pid = %d, want 1234", s.Pid)
	}
	if s.Comm != "foo) bar" {
		t.Errorf("Comm = %q, want %q", s.Comm, "foo) bar")
	}
	if s.State != 'S' {
		t.Errorf("State = %q, want S", s.State)
	}
	if s.Ppid != 10 {
		t.Errorf("Ppid = %d, want 10", s.Ppid)
	}
}

func TestStatusValue(t *testing.T) {
	got, err := StatusValue("12 kB").AsSize()
	if err != nil || got != 12*1024 {
		t.Errorf("AsSize: got %d %v, want %d", got, err, 12*1024)
	}
	u, err := StatusValue("ffffffffffffffff").AsHex()
	if err != nil || u != ^uint64(0) {
		t.Errorf("AsHex: got %x %v", u, err)
	}
}

func TestEPerm(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("root can read pid 1")
	}
	if os.Getpid() == 1 {
		t.Skip("we are pid 1")
	}
	p := Process{
		ID: 1,
	}
	_, err := p.Path()
	if err != syscall.EPERM {
		t.Errorf("Got %v, want %v", err, syscall.EPERM)
	}
}
