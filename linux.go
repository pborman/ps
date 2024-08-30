//go:build linux

package ps

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

// A Process represents a process.  A Process will cache any information
// retrieved.  Use Clean to clear the cache.  Use Processes to fetch information
// about all processes at the time of its call.
//
// Note: Changing the value of ID will not automatically clear cached
// information.
type Process struct {
	ID       int        // The process ID
	Children []*Process // Only filled in by GetProcessMap
	dir      string
	cpath    string
	stat     *Stat
	sysstat  *syscall.Stat_t
	comm     string
	cgroups  []int
	status   map[string]StatusValue
}

type StatusValue string

// commLen is the maximum length of name that Command() will return.
const commLen = 15

func (p *Process) clean() {
	p.cpath = ""
	p.stat = nil
	p.sysstat = nil
	p.comm = ""
	p.cgroups = nil
	p.status = nil
}

// A Stat contains the information from /proc/PID/stat.
// Stat is only defined for linux.
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

// Stat returns the information from /proc/PID/stat.
// Stat is only available on linux.
func (p *Process) Stat(refresh ...bool) (*Stat, error) {
	if len(refresh) > 0 && refresh[0] {
		p.stat = nil
	}
	if p.stat != nil {
		return p.stat, nil
	}
	data, err := ioutil.ReadFile(p.dirname() + "/stat")
	if err != nil {
		err = fixError(err)
		return nil, err
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
		if s != "" && s[len(s)-1] == '\n' {
			s = s[:len(s)-1]
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
	return p.stat, rerr
}

func (p *Process) footprint(refresh ...bool) (int, error) {
	if _, err := p.Stat(refresh...); err != nil {
		return 0, err
	}
	return int(p.stat.Vsize), nil
}

func (p *Process) groups() ([]int, error) {
	if p.cgroups != nil {
		return p.cgroups, nil
	}
	data, err := ioutil.ReadFile(p.dirname() + "/status")
	if err != nil {
		err = fixError(err)
		return nil, err
	}
	const Groups = "\nGroups:"
	x := bytes.Index(data, []byte(Groups))
	if x < 0 {
		return nil, errors.New("could not find Groups in status file")
	}
	data = data[x+len(Groups):]
	x = bytes.IndexByte(data, '\n')
	if x < 0 {
		return nil, errors.New("could not find Groups in status file")
	}
	groupList := bytes.Fields(data[:x])
	groups := make([]int, len(groupList))
	for i, f := range groupList {
		groups[i], err = strconv.Atoi(string(f))
		if err != nil {
			fmt.Printf("bad group: %q\n", f)
		}
	}
	p.cgroups = groups
	return p.cgroups, nil
}

const (
	stIgnore = iota
	stString
	stOctal
	stDecimal
	stHex
	stSize
	stByte
	stID // real, effective, saved, filesystem (unused)
	stArray
	stRange
	stHexList
	stSlash
)

var statusTypes = map[string]int{
	"Name":                       stString,
	"Umask":                      stOctal,
	"State":                      stByte,
	"Tgid":                       stDecimal,
	"Ngid":                       stDecimal,
	"Pid":                        stDecimal,
	"PPid":                       stDecimal,
	"TracerPid":                  stDecimal,
	"Uid":                        stArray,
	"Gid":                        stArray,
	"FDSize":                     stDecimal,
	"Groups":                     stArray,
	"NStgid":                     stDecimal,
	"NSpid":                      stDecimal,
	"NSpgid":                     stDecimal,
	"NSsid":                      stDecimal,
	"VmPeak":                     stSize,
	"VmSize":                     stSize,
	"VmLck":                      stSize,
	"VmPin":                      stSize,
	"VmHWM":                      stSize,
	"VmRSS":                      stSize,
	"RssAnon":                    stSize,
	"RssFile":                    stSize,
	"RssShmem":                   stSize,
	"VmData":                     stSize,
	"VmStk":                      stSize,
	"VmExe":                      stSize,
	"VmLib":                      stSize,
	"VmPTE":                      stSize,
	"VmSwap":                     stSize,
	"HugetlbPages":               stSize,
	"CoreDumping":                stDecimal,
	"Threads":                    stDecimal,
	"SigQ":                       stSlash,
	"SigPnd":                     stHex,
	"ShdPnd":                     stHex,
	"SigBlk":                     stHex,
	"SigIgn":                     stHex,
	"SigCgt":                     stHex,
	"CapInh":                     stHex,
	"CapPrm":                     stHex,
	"CapEff":                     stHex,
	"CapBnd":                     stHex,
	"CapAmb":                     stHex,
	"NoNewPrivs":                 stDecimal,
	"Seccomp":                    stDecimal,
	"Speculation_Store_Bypass":   stString,
	"Cpus_allowed":               stHex,
	"Cpus_allowed_list":          stRange,
	"Mems_allowed":               stHexList,
	"Mems_allowed_list":          stDecimal,
	"voluntary_ctxt_switches":    stDecimal,
	"nonvoluntary_ctxt_switches": stDecimal,
}

// StatusMap returns the values from /proc/PID/status.
// StatusMap is only availabe on linux.
func (p *Process) StatusMap(refresh ...bool) (map[string]StatusValue, error) {
	if p.status != nil && (len(refresh) == 0 || !refresh[0]) {
		return p.status, nil
	}
	data, err := ioutil.ReadFile(p.dirname() + "/status")
	if err != nil {
		err = fixError(err)
		p.status = nil
		return nil, err
	}
	p.status = map[string]StatusValue{}
	for _, line := range bytes.Split(data, []byte{'\n'}) {
		x := bytes.IndexByte(line, ':')
		if x <= 0 {
			continue // this should never happen.
		}
		p.status[string(line[:x])] = StatusValue(bytes.TrimSpace(line[x+1:]))
	}
	return p.status, nil
}

// StatusValue returns the StatusValue of the item name in /proc/PID/status.
// Each call to StatusValue reads the entire /proc/PID/status file.  Use
// StatusMap if multiple values are needed.
// StatusValue is only available on linux.
func (p *Process) StatusValue(name string) (StatusValue, error) {
	data, err := ioutil.ReadFile(p.dirname() + "/status")
	if err != nil {
		err = fixError(err)
		return "", err
	}
	bname := []byte("\n" + name + ":")
	if !bytes.HasPrefix(data, bname[1:]) {
		x := bytes.Index(data, bname)
		if x < 0 {
			return "", ErrUnset(name)
		}
		data = data[x+1:]
	}
	x := bytes.IndexByte(data, '\n')
	if x < 0 {
		x = len(data)
	}
	return StatusValue(bytes.TrimSpace(data[len(name)+1 : x])), nil
}

// AsOctal returns the numeric value of s assuming it is octal.
func (s StatusValue) AsOctal() (int64, error) {
	return strconv.ParseInt(string(s), 8, 64)
}

// AsDecimal returns the numeric value of s assuming it is decimal.
func (s StatusValue) AsDecimal() (int64, error) {
	return strconv.ParseInt(string(s), 10, 64)
}

// AsHex returns the numeric value of s assuming it is hexadecimal.
func (s StatusValue) AsHex() (int64, error) {
	return strconv.ParseInt(string(s), 16, 64)
}

// AsOctal returns the numeric value of s assuming it is an array of decimal.
func (s StatusValue) AsArray() ([]int64, error) {
	a := strings.Fields(string(s))
	vs := make([]int64, len(a))
	var err error
	for i, v := range a {
		vs[i], err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return vs, nil
}

// Creds contains the credentials for a process.
// Creds is only available on linux.
type Creds struct {
	Real      int
	Effective int
	Saved     int
	FS        int
}

// AsCreds returns s interpreted as a Creds structure.
func (s StatusValue) AsCreds() (Creds, error) {
	var creds Creds
	a := strings.Fields(string(s))
	if len(a) != 4 {
		return creds, errors.New("incorrect number of fields")
	}
	var err error
	creds.Real, err = strconv.Atoi(a[0])
	if err != nil {
		return creds, err
	}
	creds.Effective, err = strconv.Atoi(a[1])
	if err != nil {
		return creds, err
	}
	creds.Saved, err = strconv.Atoi(a[2])
	if err != nil {
		return creds, err
	}
	creds.FS, err = strconv.Atoi(a[3])
	if err != nil {
		return creds, err
	}
	return creds, nil
}

// AsSize returns the number of bytes represented by s.
func (s StatusValue) AsSize(value string) (int64, error) {
	a := strings.Fields(string(s))
	if len(a) != 2 {
		return 0, errors.New("incorrect number of fields")
	}
	v, err := strconv.ParseInt(string(a[0]), 10, 64)
	if err != nil {
		return 0, err
	}
	if string(a[1]) == "kB" {
		v *= 1024
	}
	return v, nil
}

func processByPid(pid int) (*Process, error) {
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

func processes(filled bool) ([]*Process, error) {
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

func (p *Process) pid() int {
	return p.ID
}

func (p *Process) ppid() (int, error) {
	if _, err := p.Stat(); err != nil {
		return 0, err
	}
	return p.stat.Ppid, nil
}

func (p *Process) uid() (int, error) {
	if err := p.getStat(); err != nil {
		return 0, err
	}
	return int(p.sysstat.Uid), nil
}

func (p *Process) gid() (int, error) {
	if err := p.getStat(); err != nil {
		return 0, err
	}
	return int(p.sysstat.Gid), nil
}

func (p *Process) path() (string, error) {
	var err error
	if p.cpath == "" {
		p.cpath, err = os.Readlink(p.dirname() + "/exe")
		err = fixError(err)
	}
	return p.cpath, err
}

func (p *Process) command() (string, error) {
	if _, err := p.Path(); err != nil {
		s, err := p.Stat()
		if err != nil {
			return "", err
		}
		return s.Comm, err
	}
	return p.cpath[strings.LastIndex(p.cpath, "/")+1:], nil
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

func (p *Process) argv() ([]string, error) {
	// cache this
	return p.getStrings("/cmdline")
}

func (p *Process) environ() (map[string]string, error) {
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

func (p *Process) value(name string) (string, error) {
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
	err = fixError(err)
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

func (p *Process) tty() (string, error) {
	if _, err := p.Stat(); err != nil {
		return "", err
	}
	return DevT(p.stat.TtyNr).String(), nil
}

// A DevT is a Linux device number
type DevT uint32

const noDev = 0xffffffff

func (d DevT) Major() int {
	if d == noDev {
		return -1
	}
	return int((d >> 8) & 0xff)
}

func (d DevT) Minor() int {
	if d == noDev {
		return -1
	}
	return int(d & 0xff)
}

// String returns the string form of d.  If d is -1 then "-" is returned.  The
// first call to String for any DevT caches all known device names from /dev.
func (d DevT) String() string {
	if name := getDevNames()[d]; name != "" {
		return name
	}
	return fmt.Sprintf("%d/%d", d.Major(), d.Minor())
}

var devMutex sync.RWMutex
var devNames map[DevT]string

// fillDevNames safely fills devNames if it is not already filled.
// One fillDevNames returns, devNames can be accessed without a lock.
func getDevNames() map[DevT]string {
	devMutex.RLock()
	d := devNames
	devMutex.RUnlock()
	if d != nil {
		return d
	}

	defer devMutex.Unlock()
	devMutex.Lock()
	if devNames != nil {
		return devNames
	}

	devNames = map[DevT]string{}
	devNames[noDev] = "-"
	des, err := os.ReadDir("/dev")
	if err != nil {
		return devNames
	}
	for _, de := range des {
		i, err := de.Info()
		if err != nil {
			continue
		}
		stat := i.Sys().(*syscall.Stat_t)
		if (stat.Mode & (syscall.S_IFCHR | syscall.S_IFBLK)) != 0 {
			devNames[DevT(stat.Rdev)] = de.Name()
		}
		if (stat.Mode & syscall.S_IFDIR) != 0 {
			name := de.Name()
			des, err := os.ReadDir("/dev/" + name)
			name += "/"
			if err != nil {
				continue
			}
			for _, de := range des {
				i, err := de.Info()
				if err != nil {
					continue
				}
				stat := i.Sys().(*syscall.Stat_t)
				if (stat.Mode & (syscall.S_IFCHR | syscall.S_IFBLK)) != 0 {
					devNames[DevT(stat.Rdev)] = name + de.Name()
				}
			}
		}
	}

	return devNames
}

func fixError(err error) error {
	if os.IsNotExist(err) {
		return syscall.ESRCH
	}
	if os.IsPermission(err) {
		return syscall.EPERM
	}
	return err
}
