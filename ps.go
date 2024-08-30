// The ps package provides programatic access to values normall associated with
// the ps command.  By its nature, different architctures may support different
// functionality.  The Process structure is available for all supported
// architectures.
package ps

import (
	"fmt"
	"strings"
)

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

// ProcessByName returns a list of processes with the provided name.  If name is
// an absolute pathname then the commands pathname must must match.  If name has
// no slashes then only basename of the comman must match.  If name is not
// absolute but has slashes then the trailing components of the pathname must
// mach name.
//
// Not all operating systems permit reading the full pathname of a process
// unless the caller is root or the owner of the process.  In these cases a
// different mechanism may be used to discover the command's name.  This
// mechanism may return a truncated result (e.g., linux only returns up to 15
// characters of the name).
//
// # EXAMPLES
//
// The name "/bin/ps" will only match commands whose pathname is "/bin/ps".  On
// systems such as linux only the processes owned by the caller are returned.
//
// The name "sh" will match any command named "sh" but will not match ksh.
//
// The name "bin/foo" will match "/bin/foo" but not "/sbin/foo".
//
// The name "systemd-timesynced", on linux, will match both "systemd-timesynced"
// and "systemd-timesyncd" for processes not owned by the caller (the command
// name will is truncated to "systemd-timesync")
func ProcessByName(name string) ([]*Process, error) {
	if name == "" {
		return nil, nil // Maybe EINVAL?
	}

	ps, err := Processes(false)
	if err != nil {
		return nil, err
	}
	var procs []*Process

	switch strings.Index(name, "/") {
	case 0:
		for _, p := range ps {
			path, err := p.Path()
			if err == nil && path == name {
				procs = append(procs, p)
			}
		}
	case -1:
		shortName := name
		if len(shortName) > commLen {
			shortName = name[:commLen]
		}
		for _, p := range ps {
			cmd, err := p.Command()
			if err != nil {
				continue
			}
			if name == cmd || shortName == cmd {
				procs = append(procs, p)
			}
		}
	default:
		// We have a slash that is not at the begining.
		name = "/" + name
		for _, p := range ps {
			path, err := p.Path()
			if err == nil && strings.HasSuffix(path, name) {
				procs = append(procs, p)
			}
		}
	}
	return procs, nil
}

// Argv returns p's arguments.  Non-root users will receive an error when
// requesting information about a process with a different UID.
func (p *Process) Argv() ([]string, error) {
	return p.argv()
}

// Clean discards all information known about p other than the process id.
func (p *Process) Clean() {
	p.clean()
}

// Command returns the command name of the binary associated with p.
func (p *Process) Command() (string, error) {
	return p.command()
}

// Environ returns a map of p's environment variables at time of launch.
// Non-root users will receive an error when requesting information about a
// process with a different UID.
func (p *Process) Environ() (map[string]string, error) {
	return p.environ()
}

// Footprint returns the phsycial memory footprint of p in bytes.
// Pass in the value "true" to refresh the information.
func (p *Process) Footprint(refresh ...bool) (int, error) {
	return p.footprint(refresh...)
}

// Gid returns the user id of the process.
func (p *Process) Gid() (int, error) {
	return p.gid()
}

// Groups returns the list of groups the process is in
func (p *Process) Groups() ([]int, error) {
	return p.groups()
}

// Path returns the full pathname of the binary associated with p.
func (p *Process) Path() (string, error) {
	return p.path()
}

// Pid returns the process's process ID.
func (p *Process) Pid() int {
	return p.pid()
}

// Ppid returns the process's parent process id.
func (p *Process) Ppid() (int, error) {
	return p.ppid()
}

// Tty returns the controlling tty associated with p.  "-" is returned if there
// is no associated tty.
func (p *Process) Tty() (string, error) {
	return p.tty()
}

// Uid returns the user id of the process.
func (p *Process) Uid() (int, error) {
	return p.uid()
}

// Value returns the value p's environment variable name.
func (p *Process) Value(name string) (string, error) {
	return p.value(name)
}

// ProcessByPid returns the Process associated with pid with additional
// information filled in.  The additional information is platform specific.
func ProcessByPid(pid int) (*Process, error) {
	return processByPid(pid)
}

// Processes returns a list of all processes on the system.  Setting filled to
// true will also gather the kproc_info structures for each process.  This is
// much more efficient than requesting the kproc_info structure for each
// process.
func Processes(filled bool) ([]*Process, error) {
	return processes(filled)
}
