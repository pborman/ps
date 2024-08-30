//go:build darwin

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
		if p.kinfo != nil {
			t.Fatalf("Process[%d] has kinfo filled", p.ID)
		}
		if p.rusage != nil {
			t.Fatalf("Process[%d] has rusage filled", p.ID)
		}
		if p.cpath != "" {
			t.Fatalf("Process[%d] has path filled", p.ID)
		}
		if p.argenv != nil {
			t.Fatalf("Process[%d] has argenv filled", p.ID)
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
		if p.kinfo == nil {
			t.Fatalf("Process[%d] does not have kinfo filled", p.ID)
		}
		if p.ID != int(p.kinfo.Pid) {
			t.Fatalf("Process[%d] shoudl be %d", p.ID, p.kinfo.Pid)
		}
		if p.rusage != nil {
			t.Fatalf("Process[%d] has rusage filled", p.ID)
		}
		if p.cpath != "" {
			t.Fatalf("Process[%d] has path filled", p.ID)
		}
		if p.argenv != nil {
			t.Fatalf("Process[%d] has argenv filled", p.ID)
		}
	}
	if !pids[mypid] {
		t.Fatalf("My PID was not found")
	}

	// This has to be the last test.
	for _, p := range procs {
		if mypid == int(p.kinfo.Pid) {
			if p.kinfo.Stat != SRUN {
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
	if mypid != int(p.kinfo.Pid) {
		t.Fatalf("Got process %d, want %d", p.kinfo.Pid, mypid)
	}
}

func TestClean(t *testing.T) {
	p := &Process{
		ID:     1,
		kinfo:  &KInfoProc{},
		rusage: &RUsage{},
		argenv: &argenv{},
		cpath:  ".",
	}
	p.Clean()
	if p.kinfo != nil {
		t.Errorf("kinfo not cleared")
	}
	if p.rusage != nil {
		t.Errorf("rusage not cleared")
	}
	if p.argenv != nil {
		t.Errorf("argenv not cleared")
	}
	if p.cpath != "" {
		t.Errorf("path not cleared")
	}
}

func TestKInfo(t *testing.T) {
	p := &Process{ID: mypid}
	ki1, err := p.KInfo()
	if err != nil {
		t.Fatal(err)
	}
	ki2, err := p.KInfo(false)
	if err != nil {
		t.Fatal(err)
	}
	ki3, err := p.KInfo(true)
	if err != nil {
		t.Fatal(err)
	}
	if ki1 != ki2 {
		t.Errorf("Did not return cached structure")
	}
	if ki1 == ki3 {
		t.Errorf("Returned cached structure")
	}
}

func TestRUsage(t *testing.T) {
	p := &Process{ID: mypid}
	ru1, err := p.RUsage()
	if err != nil {
		t.Fatal(err)
	}
	ru2, err := p.RUsage(false)
	if err != nil {
		t.Fatal(err)
	}
	ru3, err := p.RUsage(true)
	if err != nil {
		t.Fatal(err)
	}
	if ru1 != ru2 {
		t.Errorf("Did not return cached structure")
	}
	if ru1 == ru3 {
		t.Errorf("Returned cached structure")
	}
}

const dev003 = 0x12000003

func initDev() {
	devMutex.Lock()
	devNames = map[DevT]string{
		dev003: "dev003",
		noDev:  "-",
	}
	devMutex.Unlock()
}

func TestDevT(t *testing.T) {
	initDev()
	dev := DevT(0x12345678)
	if got, want := dev.Major(), 0x12; got != want {
		t.Errorf("Major() got %02x, want %02x", got, want)
	}
	if got, want := dev.Minor(), 0x345678; got != want {
		t.Errorf("Minor() got %06x, want %06x", got, want)
	}

	dev = DevT(dev003)

	if got, want := dev.String(), "dev003"; got != want {
		t.Errorf("device 0x%08x got name %q, want %q", uint32(dev), got, want)
	}
	dev++
	if got, want := dev.String(), "18/4"; got != want {
		t.Errorf("device 0x%08x got name %q, want %q", uint32(dev), got, want)
	}
	dev = noDev
	if got, want := dev.String(), "-"; got != want {
		t.Errorf("device 0x%08x got name %q, want %q", uint32(dev), got, want)
	}
}

func TestTty(t *testing.T) {
	initDev()
	p := &Process{
		ID:    mypid,
		kinfo: &KInfoProc{},
	}
	ki := p.kinfo
	ki.Tdev = 0x12000003
	got, err := p.Tty()
	if err != nil {
		t.Error(err)
	}
	want := "dev003"
	if got != want {
		t.Errorf("Got tty %q, want %q", got, want)
	}
}

func TestStat(t *testing.T) {
	for _, tt := range []struct {
		in   Stat
		want string
	}{
		{0, "-"},
		{SIDL, "I"},
		{SRUN, "R"},
		{SSLEEP, "S"},
		{SSTOP, "T"},
		{SZOMB, "Z"},
		{42, "42"},
	} {
		got := tt.in.String()
		if got != tt.want {
			t.Errorf("%d: got %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestEPerm(t *testing.T) {
	p := Process{
		ID: 1,
	}
	_, err := p.Argv()
	if err != syscall.EPERM {
		t.Errorf("Got %v, want %v", err, syscall.EPERM)
	}
}
