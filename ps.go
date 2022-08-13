// The ps package provides programatic access to values normall associated with
// the ps command.  By its nature, different architctures may support different
// functionality.  The Process structure is available for all supported
// architectures.
//
// # COMMON FUNCTIONS
//
// Processes - returns a slice of all processes on the system
// Process.Path - returns the pathname of the process
// Process.Command - returns the name of the process
// Process.Pid - returns the process id
// Process.PPid - returns the parent process id
//
// DARWIN (macOS) SPECIFIC FUNCTIONS
package ps

import (
	"strings"
)

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
// EXAMPLES
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
