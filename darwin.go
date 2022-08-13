//go:build darwin

package ps

import (
	"strings"
)

// commLen is the maximum length of name that Command() will return.
// zero means no limit.
const commLen = 0

// A Process represents a process.  A Process will cache any information
// retrieved.  Use Clean to clear the cache.  Use Processes to fetch information
// about all processes at the time of its call.
//
// Note: Changing the value of ID will not automatically clear cached
// information.
type Process struct {
	ID     int // The process ID
	kinfo  *KInfoProc
	rusage *RUsage
	path   string
	argenv *argenv
}

// Clean discards all information known about p other than the process id.
func (p *Process) Clean() {
	p.kinfo = nil
	p.rusage = nil
	p.argenv = nil
	p.path = ""
}

func (p *Process) Pid() int {
	return p.ID
}

// Ppid returns the process's parent process id.
func (p *Process) Ppid() (int, error) {
	if err := p.fillKinfo(); err != nil {
		return 0, err
	}
	return int(p.kinfo.Ppid), nil
}

// Uid returns the user id of the process.
func (p *Process) Uid() (int, error) {
	if err := p.fillKinfo(); err != nil {
		return 0, err
	}
	return int(p.kinfo.Uid), nil
}

// Gid returns the user id of the process.
func (p *Process) Gid() (int, error) {
	if err := p.fillKinfo(); err != nil {
		return 0, err
	}
	return int(p.kinfo.Gid), nil
}

// Groups returns the list of groups the process is in
func (p *Process) Groups() ([]int, error) {
	if err := p.fillKinfo(); err != nil {
		return nil, err
	}
	groups := make([]int, p.kinfo.Ngroups)
	for i := range groups {
		groups[i] = int(p.kinfo.Groups[i])
	}
	return groups, nil
}

// Tty returns the controlling tty associated with p.  "-" is returned if there
// is no associated tty.
func (p *Process) Tty() (string, error) {
	if err := p.fillKinfo(); err != nil {
		return "", err
	}
	return p.kinfo.Tdev.String(), nil
}

// Footprint returns the phsycial memory footprint of p in bytes.
// Pass in the value "true" to refresh the information.
func (p *Process) Footprint(refresh ...bool) (int, error) {
	if isTrue(refresh) {
		p.rusage = nil
	}
	if err := p.fillRUsage(); err != nil {
		return 0, err
	}
	return int(p.rusage.PhysFootprint), nil

}

// Path returns the full pathname of the binary associated with p.
func (p *Process) Path() (string, error) {
	if p.path != "" {
		return p.path, nil
	}
	var err error
	p.path, err = pidpath(p.ID)
	return p.path, err
}

// Command returns the command name of the binary associated with p.
func (p *Process) Command() (string, error) {
	if p.path == "" {
		var err error
		p.path, err = pidpath(p.ID)
		if err != nil {
			return "", err
		}
	}
	return p.path[strings.LastIndex(p.path, "/")+1:], nil
}

// Argv returns p's arguments.  Non-root users will receive an error when
// requesting information about a process with a different UID.
func (p *Process) Argv() ([]string, error) {
	if err := p.fillArgenv(); err != nil {
		return nil, err
	}
	return p.argenv.argv, nil
}

// Environ returns a map of p's environment variables at time of launch.
// Non-root users will receive an error when requesting information about a
// process with a different UID.
func (p *Process) Environ() (map[string]string, error) {
	if err := p.fillArgenv(); err != nil {
		return nil, err
	}
	return p.argenv.env, nil
}

// Value returns the value p's environment variable name.
func (p *Process) Value(name string) (string, error) {
	if err := p.fillArgenv(); err != nil {
		return "", err
	}
	if v, ok := p.argenv.env[name]; ok {
		return v, nil
	}
	return "", ErrUnset(name)
}

// KInfo returns the KInfoProc structure associated with p.
// Pass in the value "true" to refresh the information.
func (p *Process) KInfo(refresh ...bool) (*KInfoProc, error) {
	if isTrue(refresh) {
		p.kinfo = nil
	}
	err := p.fillKinfo()
	return p.kinfo, err
}

// RUsage returns the RUsage structure associated with p.
// Pass in the value "true" to refresh the information.
func (p *Process) RUsage(refresh ...bool) (*RUsage, error) {
	if isTrue(refresh) {
		p.rusage = nil
	}
	err := p.fillRUsage()
	return p.rusage, err
}

func (p *Process) fillKinfo() error {
	if p.kinfo != nil {
		return nil
	}
	var err error
	p.kinfo, err = getKInfoPid(p.ID)
	return err
}

func (p *Process) fillRUsage() error {
	if p.rusage != nil {
		return nil
	}
	var err error
	p.rusage, err = pidrusage(p.ID)
	return err
}

func (p *Process) fillArgenv() error {
	if p.argenv != nil {
		return nil
	}
	var err error
	p.argenv, err = getProcArgs(p.ID)
	return err
}

// ProcessByPid returns the Process associated with pid.  It is shorthand for:
//
//	p := Process{ID:pid}
//	_, err := p.KInfo(true)
func ProcessByPid(pid int) (*Process, error) {
	ki, err := getKInfoPid(pid)
	if err != nil {
		return nil, err
	}
	return &Process{
		ID:    pid,
		kinfo: ki,
	}, nil
}

// Processes returns a list of all processes on the system.  Setting filled to
// true will also gather the kproc_info structures for each process.  This is
// much more efficient than requesting the kproc_info structure for each
// process.
func Processes(filled bool) ([]*Process, error) {
	if filled {
		return fullProcesses()
	}
	pids, err := listallpids()
	if err != nil {
		return nil, err
	}
	p := make([]*Process, len(pids))
	for i, pid := range pids {
		p[i] = &Process{
			ID: int(pid),
		}
	}
	return p, nil
}

func fullProcesses() ([]*Process, error) {
	procs, err := getKInfoAll()
	if err != nil {
		return nil, err
	}
	p := make([]*Process, len(procs))
	for i, ki := range procs {
		p[i] = &Process{
			ID:    int(ki.Pid),
			kinfo: ki,
		}
	}
	return p, nil
}

func isTrue(b []bool) bool {
	if len(b) == 0 {
		return false
	}
	return b[0]
}
