//go:build linux

package ps

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// A Process represents a process.  A Process will cache any information
// retrieved.  Use Clean to clear the cache.  Use Processes to fetch information
// about all processes at the time of its call.
//
// Note: Changing the value of ID will not automatically clear cached
// information.
type Process struct {
	ID      int // The process ID
	dir     string
	stat    *Stat
	sysstat *syscall.Stat_t
	comm    string
}

type Stat struct {
	Pid                 int
	Comm                string
	State               byte
	Ppid                int
	Pgrp                int
	Session             int
	TtyNr               int
	Tpgid               int
	Flags               uint
	Minflt              uint64
	Cminflt             uint64
	Majflt              uint64
	Cmajflt             uint64
	Utime               uint64
	Stime               uint64
	Cutime              int64
	Cstime              int64
	Priority            int64
	Nice                int64
	NumThreads          int64
	Itrealvalue         int64
	Starttime           uint64
	Vsize               uint64
	Rss                 int64
	Rsslim              uint64
	Startcode           uint64
	Endcode             uint64
	Startstack          uint64
	Kstkesp             uint64
	Kstkeip             uint64
	Signal              uint64
	Blocked             uint64
	Sigignore           uint64
	Sigcatch            uint64
	Wchan               uint64
	Nswap               uint64
	Cnswap              uint64
	ExitSignal          int
	Processor           int
	RtPriority          uint
	Policy              uint
	DelayacctBlkioTicks uint64
	GuestTime           uint64
	CguestTime          int64
	StartData           uint64
	EndData             uint64
	StartBrk            uint64
	ArgStart            uint64
	ArgEnd              uint64
	EnvStart            uint64
	EnvEnd              uint64
	ExitCode            int
}

func (p *Process) readStat() error {
	if p.stat != nil {
		return nil
	}
	data, err := ioutil.ReadFile(p.dirname() + "/stat")
	if err != nil {
		return err
	}
	var rerr error
	getstring := func() string {
		if rerr != nil || len(data) == 0 {
			return ""
		}
		var i int
		var s string
		if data[0] == '(' {
			i = bytes.Index(data, []byte(") "))
			if i < 0 {
				data = nil
				return ""
			}
			s = string(data[1:i])
			data = data[i+2:]
			return s
		}
		i = bytes.IndexByte(data, ' ')
		if i < 0 {
			s = string(data)
			data = nil
		} else {
			s = string(data[:i])
			data = data[i+1:]
		}
		return s
	}
	getint := func() int {
		s := getstring()
		if s == "" {
			return 0
		}
		var i int
		i, rerr = strconv.Atoi(s)
		return i
	}
	getbyte := func() byte {
		s := getstring()
		if s == "" {
			return 0
		}
		return s[0]
	}
	getuint64 := func() uint64 {
		s := getstring()
		if s == "" {
			return 0
		}
		var u uint64
		u, rerr = strconv.ParseUint(s, 10, 64)
		return u
	}
	getuint := func() uint {
		return uint(getuint64())
	}
	getint64 := func() int64 {
		s := getstring()
		if s == "" {
			return 0
		}
		var u int64
		u, rerr = strconv.ParseInt(s, 10, 64)
		return u
	}
	var s Stat
	s.Pid = getint()
	s.Comm = getstring()
	s.State = getbyte()
	s.Ppid = getint()
	s.Pgrp = getint()
	s.Session = getint()
	s.TtyNr = getint()
	s.Tpgid = getint()
	s.Flags = getuint()
	s.Minflt = getuint64()
	s.Cminflt = getuint64()
	s.Majflt = getuint64()
	s.Cmajflt = getuint64()
	s.Utime = getuint64()
	s.Stime = getuint64()
	s.Cutime = getint64()
	s.Cstime = getint64()
	s.Priority = getint64()
	s.Nice = getint64()
	s.NumThreads = getint64()
	s.Itrealvalue = getint64()
	s.Starttime = getuint64()
	s.Vsize = getuint64()
	s.Rss = getint64()
	s.Rsslim = getuint64()
	s.Startcode = getuint64()
	s.Endcode = getuint64()
	s.Startstack = getuint64()
	s.Kstkesp = getuint64()
	s.Kstkeip = getuint64()
	s.Signal = getuint64()
	s.Blocked = getuint64()
	s.Sigignore = getuint64()
	s.Sigcatch = getuint64()
	s.Wchan = getuint64()
	s.Nswap = getuint64()
	s.Cnswap = getuint64()
	s.ExitSignal = getint()
	s.Processor = getint()
	s.RtPriority = getuint()
	s.Policy = getuint()
	s.DelayacctBlkioTicks = getuint64()
	s.GuestTime = getuint64()
	s.CguestTime = getint64()
	s.StartData = getuint64()
	s.EndData = getuint64()
	s.StartBrk = getuint64()
	s.ArgStart = getuint64()
	s.ArgEnd = getuint64()
	s.EnvStart = getuint64()
	s.EnvEnd = getuint64()
	s.ExitCode = getint()
	if rerr == nil {
		p.stat = &s
	}
	return rerr
}

func ProcessByPid(pid int) (*Process, error) {
	p := &Process{
		ID: pid,
	}
	if err := p.getStat(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Process) getStat() error {
	if p.sysstat != nil {
		return nil
	}
	var stat syscall.Stat_t
	if err := syscall.Stat(p.dirname(), &stat); err != nil {
		return err
	}
	p.sysstat = &stat
	return nil
}

func (p *Process) dirname() string {
	if p.dir == "" {
		p.dir = "/proc/" + strconv.Itoa(p.ID)
	}
	return p.dir
}

// Processes returns a list of all processes on the system.  If filled is true
// then additional information will be included.
// parameter is ignored on Linux.
func Processes(filled bool) ([]*Process, error) {
	pids, err := listallpids()
	if err != nil {
		return nil, err
	}
	p := make([]*Process, 0, len(pids))
	for _, pid := range pids {
		pr := &Process{
			ID:  pid,
			dir: "/proc/" + strconv.Itoa(pid),
		}
		if filled {
			if err := pr.getStat(); err != nil {
				continue
			}
		}
		p = append(p, pr)
	}
	return p, nil
}

func (p *Process) Pid() int {
	return p.ID
}

// Ppid returns the process's parent process id.
func (p *Process) Ppid() (int, error) {
	if err := p.readStat(); err != nil {
		return 0, err
	}
	return p.stat.Ppid, nil
}

func (p *Process) Uid() (int, error) {
	if err := p.getStat(); err != nil {
		return 0, err
	}
	return int(p.sysstat.Uid), nil
}

func (p *Process) Gid() (int, error) {
	if err := p.getStat(); err != nil {
		return 0, err
	}
	return int(p.sysstat.Gid), nil
}

func (p *Process) Path() (string, error) {
	return os.Readlink(p.dirname() + "/exe")
}

func (p *Process) Command() (string, error) {
	var err error
	if p.comm == "" {
		p.comm, err = p.stringFile("/comm")
	}
	return p.comm, err
}

func (p *Process) getStrings(name string) ([]string, error) {
	s, err := p.stringFile(name)
	if err != nil {
		return nil, err
	}
	if s[len(s)-1] == 0 {
		s = s[:len(s)-1]
	}
	return strings.Split(s, "\000"), nil
}

func (p *Process) Argv() ([]string, error) {
	// cache this
	return p.getStrings("/cmdline")
}

func (p *Process) Environ() (map[string]string, error) {
	// cache this
	fields, err := p.getStrings("/environ")
	if err != nil {
		return nil, err
	}
	env := map[string]string{}
	for _, s := range fields {
		i := strings.Index(s, "=")
		if i < 0 {
			env[s] = ""
		} else {
			env[s[:i]] = s[i+1:]
		}
	}
	return env, nil
}

// An ErrUnset is returned when requesting the value of a variable that is not
// set.
type ErrUnset string

func (e ErrUnset) Error() string {
	return fmt.Sprintf("variable not set: %s", string(e))
}

// IsUnset returns true if err is of type ErrUnset.
func IsUnset(err error) bool {
	_, ok := err.(ErrUnset)
	return ok
}

// Value returns the value p's environment variable name.
func (p *Process) Value(name string) (string, error) {
	env, err := p.Environ()
	if err != nil {
		return "", err
	}
	if v, ok := env[name]; ok {
		return v, nil
	}
	return "", ErrUnset(name)
}

func (p *Process) stringFile(name string) (string, error) {
	data, err := ioutil.ReadFile(p.dirname() + name)
	return string(data), err
}

func listallpids() ([]int, error) {
	f, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	var pids []int
	for _, name := range names {
		if pid, err := strconv.Atoi(name); err == nil {
			pids = append(pids, pid)
		}
	}
	return pids, nil
}
