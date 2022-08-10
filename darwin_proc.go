//go:build darwin

package ps

// #include <sys/sysctl.h>
import "C"

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"syscall"
	"unsafe"
)

const kinfoProcSize = C.sizeof_struct_kinfo_proc

func init() {
	// Panic if our kinfo_proc structure is the wrong size
	if unsafe.Sizeof(KInfoProc{}) != kinfoProcSize {
		panic(fmt.Sprintf("KInfoProc has size %d, want %d", unsafe.Sizeof(KInfoProc{}), kinfoProcSize))
	}
}

// Stat is the current run state
type Stat uint8

// A pointer is a pointer in the C structure, it is unused
type pointer int64

// We should provide a function to convert this value to a float32
type Fixpt uint32

// KInfoProc is the kinfo_proc from sysctl.h
// Some public fields might not be filled in.
// The fields p_ruid, p_rgid, and cr_uid are renamed to Uid, Gid, and Euid
// for convience.
type KInfoProc struct {
	ExternProc
	EProc
}

// ITimerval is an interval timer
type ITimerval struct {
	Interval syscall.Timeval
	Value    syscall.Timeval
}

const (
	/* Status values. */
	SIDL   = Stat(1) // Process being created by fork.
	SRUN   = Stat(2) // Currently runnable.
	SSLEEP = Stat(3) // Sleeping on an address.
	SSTOP  = Stat(4) // Process debugging or suspension.
	SZOMB  = Stat(5) // Awaiting collection by parent.
)

var states = []string{
	0: "-",
	SIDL: "I", /* Process being created by fork. */
	SRUN: "R", /* Currently runnable. */
	SSLEEP: "S", /* Sleeping on an address. */
	SSTOP: "T", /* Process debugging or suspension. */
	SZOMB: "Z", /* Awaiting collection by parent. */
}

func (s Stat) String() string {
	if uint(s) < uint(len(states)) {
		return states[s]
	}
	return fmt.Sprintf("%d", s)
}

// ExternProc is extern_proc from proc.h
type ExternProc struct {
	Starttime syscall.Timeval // p_starttime - process start time
	_         pointer         // p_vmspace
	_         pointer         // p_sigacts
	Flag      uint32          // p_flag - P_* flags
	Stat      Stat            // p_stat - S* process status
	_         [3]byte
	Pid       uint32 // p_pid - Process identifier
	_         uint32 // p_oppid - Save parent pid during ptrace. XXX
	_         int32  // p_odupfc - Sideways return value from fdopen. XXX
	_         [4]byte
	_         pointer // user_stack
	_         pointer // exit_thread
	Debugger  uint32  // p_debugger - allow to debug
	Sigwait   int32   // sigwait - indication to suspend
	Estcpu    uint32  // p_estcpu - Time averaged value of p_cpticks
	Cpticks   int32   // p_cpticks - Ticks of cpu time
	Pctcpu    Fixpt   // p_pctcpu - %cpu for this process during p_swtime
	_         [4]byte
	_         pointer         // p_wchan
	_         pointer         // p_wmesg
	Swtime    uint32          // p_swtime - Time swapped in or out
	Slptime   uint32          // p_slpime - Time since last blocked
	Realtimer ITimerval       // p_realtimer - Alarm Timer
	Rtime     syscall.Timeval // p_rtime - Real time
	Uticks    uint64          // p_uticks - Statclock hits in user mode
	Sticks    uint64          // p_sticks - Statclock hits in system mode
	Iticks    uint64          // p_iticks - Statclock hits processing intr
	Traceflag int32           // p_traceflag - Kernel trace points
	_         [4]byte
	_         pointer // p_tracep
	_         int32   // p_siglist - DEPRECATED
	_         [4]byte
	_         pointer  // p_textvp
	Holdcnt   int32    // p_holdcnt - If non-zero, don't swap
	_         uint32   // p_sigmask - DEPRECATED
	Sigignore uint32   // p_sigignore - Signals being ignored
	Sigcatch  uint32   // p_sigcatch - Signals being caught by user
	Priority  uint8    // p_priority - Process priority
	Usrpri    uint8    // p_usrpri - User-priority based on p_cpu and p_nice
	Nice      int8     // p_nice - Process "nice" value
	Comm      [17]byte // p_comm
	_         [4]byte
	Pgrp      pointer // p_pgrp
	_         pointer // p_addr
	Xstat     uint16  // p_xstat - Exit status for wait; also stop signal
	Acflag    uint16  // p_acflag - Accounting flags
	_         [4]byte
	_         pointer // p_ru
}

// EProc is eproc from sysctl.h
type EProc struct {
	_       pointer  // e_paddr
	_       pointer  // e_sess
	_       [72]byte // e_pcred.pc_lock
	_       pointer  // e_pcred.pc_ucred
	Uid     int32    // e_pcred.p_ruid - Real user id
	Svuid   int32    // e_pcred.p_svuid - Saved effective user id
	Gid     int32    // e_pcred.p_rgid - Real group id
	Svgid   int32    // e_pcred.p_svgid - Saved effective group id
	Refcnt  int32    // e_pcred.p_refcnt - Number of references
	_       [4]byte
	Ref     int32 // e_ucred.cr_ref - reference count
	Euid    int32 // e_ucred.cr_uid - effective user id
	Ngroups int16 // e_ucred.cr_ngroups - number of groups
	_       [2]byte
	Groups  [16]int32 // e_ucred.cr_groups - groups
	_       [4]byte

	_     [64]byte // e_vm - address space
	Ppid  int32    // e_ppid - parent process id
	Pgid  int32    // e_pgid - process group id
	Jobc  int16    // e_jobc - job control counter
	_     [2]byte
	Tdev  DevT  // e_tdev - controlling tty dev
	Tpgid int32 // e_tpgid - tty process group id

	// Everything from here on down is always zero
	_       [4]byte
	_       int16   // e_tsess
	_       int16   // e_tsess
	_       int32   // e_tsess
	_       [8]byte // e_wmesg - wchan message
	Xsize   uint32  // e_xsize - text size
	Xrssize int16   // e_xrssize - text rss
	Xccount int16   // e_xccount - text references
	Xswrss  int16   // e_xswrss
	_       [2]byte
	Eflag   uint32   // e_flag
	_       [12]byte // e_login - short setlogin() name
	_       [16]byte // e_spare
	_       [4]byte  // Padding
}

// EFlag constants
const (
	EPROC_CTTY    = 0x01 /* controlling tty vnode active */
	EPROC_SLEADER = 0x02 /* session leader */

)

const (
	/* These flags are kept in extern_proc.p_flag. */
	P_ADVLOCK   = 0x00000001 /* Process may hold POSIX adv. lock */
	P_CONTROLT  = 0x00000002 /* Has a controlling terminal */
	P_LP64      = 0x00000004 /* Process is LP64 */
	P_NOCLDSTOP = 0x00000008 /* No SIGCHLD when children stop */

	P_PPWAIT    = 0x00000010 /* Parent waiting for chld exec/exit */
	P_PROFIL    = 0x00000020 /* Has started profiling */
	P_SELECT    = 0x00000040 /* Selecting; wakeup/waiting danger */
	P_CONTINUED = 0x00000080 /* Process was stopped and continued */

	P_SUGID   = 0x00000100 /* Has set privileges since last exec */
	P_SYSTEM  = 0x00000200 /* Sys proc: no sigs, stats or swap */
	P_TIMEOUT = 0x00000400 /* Timing out during sleep */
	P_TRACED  = 0x00000800 /* Debugged process being traced */

	P_DISABLE_ASLR = 0x00001000 /* Disable address space layout randomization */
	P_WEXIT        = 0x00002000 /* Working on exiting */
	P_EXEC         = 0x00004000 /* Process called exec. */

	/* Should be moved to machine-dependent areas. */
	P_OWEUPC = 0x00008000 /* Owe process an addupc() call at next ast. */

	P_AFFINITY   = 0x00010000   /* xxx */
	P_TRANSLATED = 0x00020000   /* xxx */
	P_CLASSIC    = P_TRANSLATED /* xxx */

	P_DELAYIDLESLEEP = 0x00040000 /* Process is marked to delay idle sleep on disk IO */
	P_CHECKOPENEVT   = 0x00080000 /* check if a vnode has the OPENEVT flag set on open */

	P_DEPENDENCY_CAPABLE = 0x00100000 /* process is ok to call vfs_markdependency() */
	P_REBOOT             = 0x00200000 /* Process called reboot() */
	P_RESV6              = 0x00400000 /* used to be P_TBE */
	P_RESV7              = 0x00800000 /* (P_SIGEXC)signal exceptions */

	P_THCWD        = 0x01000000 /* process has thread cwd  */
	P_RESV9        = 0x02000000 /* (P_VFORK)process has vfork children */
	P_ADOPTPERSONA = 0x04000000 /* process adopted a persona (used to be P_NOATTACH) */
	P_RESV11       = 0x08000000 /* (P_INVFORK) proc in vfork */

	P_NOSHLIB = 0x10000000 /* no shared libs are in use for proc */
	/* flag set on exec */
	P_FORCEQUOTA   = 0x20000000 /* Force quota for root */
	P_NOCLDWAIT    = 0x40000000 /* No zombies when chil procs exit */
	P_NOREMOTEHANG = 0x80000000 /* Don't hang on remote FS ops */

	P_INMEM   = 0 /* Obsolete: retained for compilation */
	P_NOSWAP  = 0 /* Obsolete: retained for compilation */
	P_PHYSIO  = 0 /* Obsolete: retained for compilation */
	P_FSTRACE = 0 /* Obsolete: retained for compilation */
	P_SSTEP   = 0 /* Obsolete: retained for compilation */

	P_DIRTY_TRACK              = 0x00000001 /* track dirty state */
	P_DIRTY_ALLOW_IDLE_EXIT    = 0x00000002 /* process can be idle-exited when clean */
	P_DIRTY_DEFER              = 0x00000004 /* defer initial opt-in to idle-exit */
	P_DIRTY                    = 0x00000008 /* process is dirty */
	P_DIRTY_SHUTDOWN           = 0x00000010 /* process is dirty during shutdown */
	P_DIRTY_TERMINATED         = 0x00000020 /* process has been marked for termination */
	P_DIRTY_BUSY               = 0x00000040 /* serialization flag */
	P_DIRTY_MARKED             = 0x00000080 /* marked dirty previously */
	P_DIRTY_AGING_IN_PROGRESS  = 0x00000100 /* aging in one of the 'aging bands' */
	P_DIRTY_LAUNCH_IN_PROGRESS = 0x00000200 /* launch is in progress */
	P_DIRTY_DEFER_ALWAYS       = 0x00000400 /* defer going to idle-exit after every dirty->clean transition. For legacy jetsam policy only. This is the default with the other policies.*/
	P_DIRTY_IS_DIRTY           = (P_DIRTY | P_DIRTY_SHUTDOWN)
	P_DIRTY_IDLE_EXIT_ENABLED  = (P_DIRTY_TRACK | P_DIRTY_ALLOW_IDLE_EXIT)
)

func getKInfoAll() ([]*KInfoProc, error) {
	data, err := sysctl([]int32{_CTL_KERN, _KERN_PROC, _KERN_PROC_ALL})
	if err != nil {
		return nil, err
	}
	procs := make([]*KInfoProc, 0, len(data)/kinfoProcSize)
	k := 0
	for i := kinfoProcSize; i <= len(data); i += kinfoProcSize {
		proc := &KInfoProc{}
		proc, err = mkKInfoProc(data[k:i])
		k = i
		if err != nil {
			return nil, err
		}
		procs = append(procs, proc)
	}
	return procs, nil
}

func getKInfoPid(pid int) (*KInfoProc, error) {
	data, err := sysctl([]int32{_CTL_KERN, _KERN_PROC, _KERN_PROC_PID, int32(pid)})
	if err != nil {
		return nil, err
	}
	if len(data) != kinfoProcSize {
		// XXX
		return nil, errors.New("bad return from sysctl")
	}
	return mkKInfoProc(data)
}

// mkKInfoProc returns data as a pointer to a KInfoProc structure.
// An error is returned if length of data is not the same size as
// a KInfoProc.
func mkKInfoProc(data []byte) (*KInfoProc, error) {
	if uintptr(len(data)) != kinfoProcSize {
		return nil, errors.New(fmt.Sprintf("KInfoProc is %d bytes, got %d", kinfoProcSize, len(data)))
	}
	var ki *KInfoProc
	sh := *(*reflect.SliceHeader)(unsafe.Pointer(&data))
	*(*uintptr)(unsafe.Pointer(&ki)) = sh.Data
	return ki, nil
}

// maxProc returns the maximum number of procces on the system.
func maxProc() (int, error) {
	buf, err := sysctl([]int32{_CTL_KERN, _KERN_MAXPROC})
	if err != nil {
		return 0, err
	}
	return int(*(*uint32)(unsafe.Pointer(&buf[0]))), nil
}

type argenv struct {
	command string
	argv    []string
	env     map[string]string
}

func getProcArgs(pid int) (*argenv, error) {
	buf, err := sysctl([]int32{_CTL_KERN, _KERN_PROCARGS2, int32(pid)})
	if err != nil {
		if err != syscall.EINVAL {
			return nil, err
		}
		// EINVAL can mean the process does not exist
		// or it is owned by someone else.
		pids, _ := listallpids()
		for _, p := range pids {
			if int(p) == pid {
				return nil, syscall.EPERM
			}
		}
		return nil, syscall.ESRCH
	}
	getInt32 := func() int {
		if len(buf) < 4 {
			buf = nil
			return 0
		}
		i := binary.LittleEndian.Uint32(buf)
		buf = buf[4:]
		return int(i)
	}
	getString := func() string {
		if len(buf) == 0 {
			return ""
		}
		i := bytes.IndexByte(buf, 0)
		if i >= 0 {
			s := string(buf[:i])
			buf = buf[i+1:]
			return s
		}
		s := string(buf)
		buf = nil
		return s
	}
	var ae argenv
	ae.env = map[string]string{}
	argc := getInt32()
	ae.command = getString()
	// The first string is padded to a boundary of 8 bytes
	if i := (len(ae.command) + 1) & 7; i != 0 {
		buf = buf[8-i:]
	}
	if argc > 0 {
		ae.argv = make([]string, argc)
		for i := range ae.argv {
			ae.argv[i] = getString()
		}
	}
	for {
		s := getString()
		if s == "" {
			return &ae, nil
		}
		i := strings.Index(s, "=")
		if i < 0 {
			ae.env[s] = ""
		} else {
			ae.env[s[:i]] = s[i+1:]
		}
	}
}
