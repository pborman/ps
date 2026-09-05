package ps

import (
	"os"
	"sort"
	"strings"
	"syscall"
	"testing"
)

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
	extra := make([]int, 16*1024*1024)
	for i := range extra {
		extra[i] = 1
	}
	second, err := p.Footprint(true)
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Errorf("Footprint did not change from %d", first)
	}
	if extra[0] != 1 {
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
	if len(argv) != len(os.Args) {
		t.Fatalf("got %d parameters, want %d", len(argv), len(os.Args))
	}
	for i, arg := range os.Args {
		if argv[i] != arg {
			t.Errorf("argv[%d] is %q, want %q", i, argv[i], arg)
		}
	}
}

func deadPid(t *testing.T) int {
	sp, err := os.StartProcess("/bin/sleep", []string{"sleep", "60"}, &os.ProcAttr{})
	if err != nil {
		t.Fatal(err)
	}
	sp.Kill()
	sp.Wait()
	return sp.Pid
}

func TestMissingPid(t *testing.T) {
	pid := deadPid(t)
	p := &Process{ID: pid}
	if _, err := ProcessByPid(pid); err != syscall.ESRCH {
		t.Errorf("ProcessByPid: got %T %v, want ESRCH", err, err)
	}
	if _, err := p.Argv(); err != syscall.ESRCH {
		t.Errorf("Argv: got %T %v, want ESRCH", err, err)
	}
	if _, err := p.Footprint(); err != syscall.ESRCH {
		t.Errorf("Footprint: got %T %v, want ESRCH", err, err)
	}
	if _, err := p.Footprint(); err != syscall.ESRCH {
		t.Errorf("Footprint again: got %T %v, want ESRCH", err, err)
	}
	if _, err := p.Fds(); err != syscall.ESRCH {
		t.Errorf("Fds: got %T %v, want ESRCH", err, err)
	}
}

func TestFds(t *testing.T) {
	f, err := os.CreateTemp("", "ps-fd-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	p := &Process{ID: mypid}
	fds, err := p.Fds()
	if err != nil {
		t.Fatal(err)
	}
	if len(fds) == 0 {
		t.Fatal("no file descriptors")
	}
	want := int(f.Fd())
	found := false
	for _, fd := range fds {
		if fd.Fd != want {
			continue
		}
		found = true
		if fd.Path == "" {
			t.Errorf("fd %d has empty path", fd.Fd)
		} else if !strings.HasSuffix(fd.Path, f.Name()) && !strings.Contains(fd.Path, "ps-fd-") {
			t.Errorf("fd %d path %q does not match %q", fd.Fd, fd.Path, f.Name())
		}
	}
	if !found {
		t.Errorf("did not find fd %d (%s) in %d descriptors", want, f.Name(), len(fds))
	}
}

func TestProcessByName(t *testing.T) {
	p := &Process{ID: mypid}
	cmd, err := p.Command()
	if err != nil {
		t.Fatal(err)
	}
	procs, err := ProcessByName(cmd)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, proc := range procs {
		if proc.ID == mypid {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("ProcessByName(%q) did not find pid %d", cmd, mypid)
	}

	path, err := p.Path()
	if err != nil {
		t.Fatal(err)
	}
	procs, err = ProcessByName(path)
	if err != nil {
		t.Fatal(err)
	}
	found = false
	for _, proc := range procs {
		if proc.ID == mypid {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("ProcessByName(%q) did not find pid %d", path, mypid)
	}

	procs, err = ProcessByName("")
	if err != nil || procs != nil {
		t.Errorf("ProcessByName(\"\") got %v %v, want nil nil", procs, err)
	}
}

func TestMatchTrailingPath(t *testing.T) {
	for _, tt := range []struct {
		path, name string
		want       bool
	}{
		{"/bin/foo", "bin/foo", true},
		{"/usr/bin/foo", "bin/foo", true},
		{"/sbin/foo", "bin/foo", false},
		{"/bin/foo", "n/foo", false},
		{"/n/foo", "n/foo", true},
	} {
		if got := matchTrailingPath(tt.path, tt.name); got != tt.want {
			t.Errorf("matchTrailingPath(%q, %q) = %v, want %v", tt.path, tt.name, got, tt.want)
		}
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
