// The ps package provides programatic access to values normall associated with
// the ps command.  By its nature, different architctures may support different
// functionality.  The Process structure is available for all supported
// architectures.
//
// COMMON FUNCTIONS
//
// Processes - returns a slice of all processes on the system
// Process.Path - returns the pathname of the process
// Process.Command - returns the name of the process
// Process.Pid - returns the process id
// Process.PPid - returns the parent process id
//
// DARWIN (macOS) SPECIFIC FUNCTIONS
package ps
