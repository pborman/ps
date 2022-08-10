//go:build darwin

package ps

//#include <sys/types.h>
//#include <sys/resource.h>
//#include <libproc.h>
import "C"

import (
	"unsafe"
)

func pidpath(pid int) (string, error) {
	const (
		MAXPATHLEN               = 1024
		PROC_PIDPATHINFO_MAXSIZE = 4 * MAXPATHLEN
	)
	var buf [PROC_PIDPATHINFO_MAXSIZE]byte
	i, err := C.proc_pidpath(C.int(pid), unsafe.Pointer(&buf), C.uint32_t(len(buf)))
	if err != nil {
		return "", err
	}
	return string(buf[:i]), nil
}

const (
	_RUSAGE_INFO_V0 = iota
	_RUSAGE_INFO_V1
	_RUSAGE_INFO_V2
	_RUSAGE_INFO_V3
	_RUSAGE_INFO_V4
	_RUSAGE_INFO_V5
	_RUSAGE_INFO_CURRENT = _RUSAGE_INFO_V5
)

// This is the Go representaion of struct rusage_info_v5 from sys/resource.h
type RUsage struct {
	Uuid                      [16]byte
	UserTime                  uint64
	SystemTime                uint64
	PkgIdleWkups              uint64
	InterruptWkups            uint64
	Pageins                   uint64
	WiredSize                 uint64
	ResidentSize              uint64
	PhysFootprint             uint64
	ProcStartAbstime          uint64
	ProcExitAbstime           uint64
	ChildUserTime             uint64
	ChildSystemTime           uint64
	ChildPkgIdleWkups         uint64
	ChildInterruptWkups       uint64
	ChildPageins              uint64
	ChildElapsedAbstime       uint64
	DiskioBytesread           uint64
	DiskioByteswritten        uint64
	CpuTimeQosDefault         uint64
	CpuTimeQosMaintenance     uint64
	CpuTimeQosBackground      uint64
	CpuTimeQosUtility         uint64
	CpuTimeQosLegacy          uint64
	CpuTimeQosUserInitiated   uint64
	CpuTimeQosUserInteractive uint64
	BilledSystemTime          uint64
	ServicedSystemTime        uint64
	LogicalWrites             uint64
	LifetimeMaxPhysFootprint  uint64
	Instructions              uint64
	Cycles                    uint64
	BilledEnergy              uint64
	ServicedEnergy            uint64
	IntervalMaxPhysFootprint  uint64
	RunnableTime              uint64
	Flags                     uint64
}

func pidrusage(pid int) (*RUsage, error) {
	var ri RUsage
	var vri = (*C.rusage_info_t)(unsafe.Pointer(&ri))
	_, err := C.proc_pid_rusage(C.int(pid), _RUSAGE_INFO_V5, vri)
	return &ri, err
}

func listallpids() ([]int32, error) {
	maxproc, err := maxProc()
	if err != nil {
		return nil, err
	}
	pids := make([]int32, maxproc)
	i, err := C.proc_listallpids(unsafe.Pointer(&pids[0]), C.int(len(pids)*4))
	if err != nil {
		return nil, err
	}
	return pids[:i], nil
}
