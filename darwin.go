//go:build darwin

package ps

import "strings"

// commLen is the maximum length of name that Command() will return.
// zero means no limit.
const commLen = 0

// A Process represents a process.  A Process caches any information retrieved
// for a process.  Use Clean to clear the cache.  Use Processes to fetch
// information about all processes at the time of its call.
//
// Note: Changing the value of ID will not automatically clear cached
// information.
type Process struct {
	ID       int        // The process ID
	Children []*Process // Only filled in by GetProcessMap
	kinfo    *KInfoProc
	rusage   *RUsage
	cpath    string
	argenv   *argenv
}

func (p *Process) clean() {
	p.kinfo = nil
	p.rusage = nil
	p.argenv = nil
	p.cpath = ""
}

func (p *Process) pid() int {
	return p.ID
}

func (p *Process) ppid() (int, error) {
	if err := p.fillKinfo(); err != nil {
		return 0, err
	}
	return int(p.kinfo.Ppid), nil
}

func (p *Process) uid() (int, error) {
	if err := p.fillKinfo(); err != nil {
		return 0, err
	}
	return int(p.kinfo.Uid), nil
}

func (p *Process) gid() (int, error) {
	if err := p.fillKinfo(); err != nil {
		return 0, err
	}
	return int(p.kinfo.Gid), nil
}

func (p *Process) groups() ([]int, error) {
	if err := p.fillKinfo(); err != nil {
		return nil, err
	}
	groups := make([]int, p.kinfo.Ngroups)
	for i := range groups {
		groups[i] = int(p.kinfo.Groups[i])
	}
	return groups, nil
}

func (p *Process) tty() (string, error) {
	if err := p.fillKinfo(); err != nil {
		return "", err
	}
	return p.kinfo.Tdev.String(), nil
}

func (p *Process) footprint(refresh ...bool) (int, error) {
	if isTrue(refresh) {
		p.rusage = nil
	}
	if err := p.fillRUsage(); err != nil {
		return 0, err
	}
	return int(p.rusage.PhysFootprint), nil

}

func (p *Process) path() (string, error) {
	if p.cpath != "" {
		return p.cpath, nil
	}
	var err error
	p.cpath, err = pidpath(p.ID)
	return p.cpath, err
}

func (p *Process) command() (string, error) {
	if p.cpath == "" {
		var err error
		p.cpath, err = pidpath(p.ID)
		if err != nil {
			return "", err
		}
	}
	return p.cpath[strings.LastIndex(p.cpath, "/")+1:], nil
}

func (p *Process) argv() ([]string, error) {
	if err := p.fillArgenv(); err != nil {
		return nil, err
	}
	return p.argenv.argv, nil
}

func (p *Process) environ() (map[string]string, error) {
	if err := p.fillArgenv(); err != nil {
		return nil, err
	}
	return p.argenv.env, nil
}

func (p *Process) value(name string) (string, error) {
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
// KInfo is only available on darwin.
func (p *Process) KInfo(refresh ...bool) (*KInfoProc, error) {
	if isTrue(refresh) {
		p.kinfo = nil
	}
	err := p.fillKinfo()
	return p.kinfo, err
}

// RUsage returns the RUsage structure associated with p.
// Pass in the value "true" to refresh the information.
// RUsage is only available ond darwin.
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
func processByPid(pid int) (*Process, error) {
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
func processes(filled bool) ([]*Process, error) {
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
