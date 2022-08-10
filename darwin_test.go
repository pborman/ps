//go:build darwin

package ps

import (
	"os"
	"sort"
	"strings"
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
		if p.path != "" {
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
		if p.path != "" {
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
		path:   ".",
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
	if p.path != "" {
		t.Errorf("path not cleared")
	}
}

func TestPid(t *testing.T) {
	p := &Process{ID: 42}
	if p.Pid() != 42 {
		t.Errorf("Got id %d, want 42", p.Pid())
	}
}

func TestPpid(t *testing.T) {
	p := &Process{ID: mypid}
	ppid, err := p.Ppid()
	if err != nil {
		t.Fatal(err)
	}
	if ppid != os.Getppid() {
		t.Errorf("Got ppid %d, want %d", ppid, os.Getppid())
	}
}

func TestUid(t *testing.T) {
	p := &Process{ID: mypid}
	uid, err := p.Uid()
	if err != nil {
		t.Fatal(err)
	}
	if uid != os.Getuid() {
		t.Errorf("Got uid %d, want %d", uid, os.Getuid())
	}
}

func TestGid(t *testing.T) {
	p := &Process{ID: mypid}
	gid, err := p.Gid()
	if err != nil {
		t.Fatal(err)
	}
	if gid != os.Getgid() {
		t.Errorf("Got uid %d, want %d", gid, os.Getgid())
	}
}

func TestGroups(t *testing.T) {
	p := &Process{ID: mypid}
	groups, err := p.Groups()
	if err != nil {
		t.Fatal(err)
	}
	ogroups, err := os.Getgroups()
	if err != nil {
		t.Fatal(err)
	}
	sort.Ints(groups)
	sort.Ints(ogroups)
	if len(groups) != len(ogroups) {
		t.Fatalf("Got %d groups, want %d", len(groups), len(ogroups))
	}
	for i, gid := range groups {
		if ogroups[i] != gid {
			t.Fatalf("Got group %d, want %d", gid, ogroups[i])
		}
	}
}

func TestFootprint(t *testing.T) {
	p := &Process{ID: mypid}
	first, err := p.Footprint()
	if err != nil {
		t.Fatal(err)
	}
	// Force our footprint to increase
	extra := make([]int, 1024*1024)
	second, err := p.Footprint(true)
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Errorf("Footprint did not change from %d", first)
	}
	if extra[0] != 0 {
		t.Fatalf("this cannot happen")
	}
}

func TestPath(t *testing.T) {
	p := &Process{ID: mypid}
	path, err := p.Path()
	if err != nil {
		t.Fatal(err)
	}
	if path == "" {
		t.Errorf("returned empty path")
	}
	command, err := p.Command()
	if err != nil {
		t.Fatal(err)
	}
	if command == "" {
		t.Errorf("returned empty command")
	}
	if !strings.HasSuffix(path, "/"+command) {
		t.Errorf("Command %q not the suffix of %q", "/"+command, path)
	}
}

func TestArgv(t *testing.T) {
	p := &Process{ID: mypid}
	argv, err := p.Argv()
	if err != nil {
		t.Fatal(err)
	}
	if len(argv) == 0 {
		t.Errorf("failed to get argv")
	}
	command, err := p.Command()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(argv[0], "/"+command) {
		t.Errorf("Arg[0] %q does not have the suffix of %q", argv[0], "/"+command)
	}
}

func TestEnviron(t *testing.T) {
	p := &Process{ID: mypid}
	env, err := p.Environ()
	if err != nil {
		t.Fatal(err)
	}
	if env["HOME"] == "" {
		t.Errorf("Could not find HOME")
	}
}

func TestValue(t *testing.T) {
	p := &Process{ID: mypid}
	value, err := p.Value("HOME")
	if err != nil {
		t.Fatal(err)
	}
	if value == "" {
		t.Errorf("Could not find HOME")
	}
	const badName = "==="
	_, err = p.Value(badName)
	if err == nil {
		t.Fatalf("returned value for impossible name")
	}
	if !IsUnset(err) {
		t.Fatalf("returned error type %T, want %T", err, ErrUnset("x"))
	}
	if string(err.(ErrUnset)) != badName {
		t.Fatalf("got error for %q, want it for %q", string(err.(ErrUnset)), badName)
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

func TestGetDevNames(t *testing.T) {
	devMutex.Lock()
	devNames = nil
	devMutex.Unlock()
	devMap := getDevNames()
	if len(devMap) == 0 {
		t.Errorf("Did not get any devices")
	}
	var v DevT
	for v, _ = range devMap {
		break
	}
	devMutex.Lock()
	devNames[v] = "changed"
	devMutex.Unlock()
	if v.String() != "changed" {
		t.Errorf("Got dev name %q, want %q", v.String(), "changed")
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
