// package data is only used for testing.
package data

//#include <sys/sysctl.h>
//#include <sys/resource.h>
import "C"
import "unsafe"

type KInfo struct {
	Data  []byte
	kinfo *C.struct_kinfo_proc
}

func KInfoData() KInfo {
	data := make([]byte, C.sizeof_struct_kinfo_proc)
	for i := range data {
		data[i] = byte(i)
	}
	return KInfo{
		Data:  data,
		kinfo: (*C.struct_kinfo_proc)(unsafe.Pointer(&data[0])),
	}
}

func (k KInfo) Flag() uint32      { return uint32(k.kinfo.kp_proc.p_flag) }
func (k KInfo) Stat() uint32      { return uint32(k.kinfo.kp_proc.p_stat) }
func (k KInfo) Pid() uint32       { return uint32(k.kinfo.kp_proc.p_pid) }
func (k KInfo) Debugger() uint32  { return uint32(k.kinfo.kp_proc.p_debugger) }
func (k KInfo) Sigwait() int32    { return int32(k.kinfo.kp_proc.sigwait) }
func (k KInfo) Estcpu() uint32    { return uint32(k.kinfo.kp_proc.p_estcpu) }
func (k KInfo) Cpticks() int32    { return int32(k.kinfo.kp_proc.p_cpticks) }
func (k KInfo) Pctcpu() uint32    { return uint32(k.kinfo.kp_proc.p_pctcpu) }
func (k KInfo) Swtime() uint32    { return uint32(k.kinfo.kp_proc.p_swtime) }
func (k KInfo) Slptime() uint32   { return uint32(k.kinfo.kp_proc.p_slptime) }
func (k KInfo) Uticks() uint64    { return uint64(k.kinfo.kp_proc.p_uticks) }
func (k KInfo) Sticks() uint64    { return uint64(k.kinfo.kp_proc.p_sticks) }
func (k KInfo) Iticks() uint64    { return uint64(k.kinfo.kp_proc.p_iticks) }
func (k KInfo) Traceflag() int32  { return int32(k.kinfo.kp_proc.p_traceflag) }
func (k KInfo) Holdcnt() int32    { return int32(k.kinfo.kp_proc.p_holdcnt) }
func (k KInfo) Sigignore() uint32 { return uint32(k.kinfo.kp_proc.p_sigignore) }
func (k KInfo) Sigcatch() uint32  { return uint32(k.kinfo.kp_proc.p_sigcatch) }
func (k KInfo) Priority() uint8   { return uint8(k.kinfo.kp_proc.p_priority) }
func (k KInfo) Usrpri() uint8     { return uint8(k.kinfo.kp_proc.p_usrpri) }
func (k KInfo) Nice() int8        { return int8(k.kinfo.kp_proc.p_nice) }
func (k KInfo) Xstat() uint16     { return uint16(k.kinfo.kp_proc.p_xstat) }
func (k KInfo) Acflag() uint16    { return uint16(k.kinfo.kp_proc.p_acflag) }
func (k KInfo) Uid() int32        { return int32(k.kinfo.kp_eproc.e_pcred.p_ruid) }
func (k KInfo) Svuid() int32      { return int32(k.kinfo.kp_eproc.e_pcred.p_svuid) }
func (k KInfo) Gid() int32        { return int32(k.kinfo.kp_eproc.e_pcred.p_rgid) }
func (k KInfo) Svgid() int32      { return int32(k.kinfo.kp_eproc.e_pcred.p_svgid) }
func (k KInfo) Refcnt() int32     { return int32(k.kinfo.kp_eproc.e_pcred.p_refcnt) }
func (k KInfo) Ref() int32        { return int32(k.kinfo.kp_eproc.e_ucred.cr_ref) }
func (k KInfo) Euid() int32       { return int32(k.kinfo.kp_eproc.e_ucred.cr_uid) }
func (k KInfo) Ppid() int32       { return int32(k.kinfo.kp_eproc.e_ppid) }
func (k KInfo) Pgid() int32       { return int32(k.kinfo.kp_eproc.e_pgid) }
func (k KInfo) Tpgid() int32      { return int32(k.kinfo.kp_eproc.e_tpgid) }
func (k KInfo) Tdev() uint32      { return uint32(k.kinfo.kp_eproc.e_tdev) }
func (k KInfo) Xsize() uint32     { return uint32(k.kinfo.kp_eproc.e_xsize) }
func (k KInfo) Eflag() uint32     { return uint32(k.kinfo.kp_eproc.e_flag) }
func (k KInfo) Xrssize() int16    { return int16(k.kinfo.kp_eproc.e_xrssize) }
func (k KInfo) Xccount() int16    { return int16(k.kinfo.kp_eproc.e_xccount) }
func (k KInfo) Xswrss() int16     { return int16(k.kinfo.kp_eproc.e_xswrss) }
func (k KInfo) Ngroups() int16    { return int16(k.kinfo.kp_eproc.e_ucred.cr_ngroups) }
func (k KInfo) Jobc() int16       { return int16(k.kinfo.kp_eproc.e_jobc) }
func (k KInfo) Comm() [17]byte    {
	var comm [17]byte
	for i := range comm {
		comm[i] = byte(k.kinfo.kp_proc.p_comm[i])
	}
	return comm
}
func (k KInfo) Groups() [16]int32 {
	var groups [16]int32
	for i := range groups {
		groups[i] = int32(k.kinfo.kp_eproc.e_ucred.cr_groups[i])
	}
	return groups
}

/*
	The following are not being tested:

	Starttime syscall.Timeval // p_starttime - process start time
	Realtimer ITimerval       // p_realtimer - Alarm Timer
	Rtime     syscall.Timeval // p_rtime - Real time
	Usrpri    uint8    // p_usrpri - User-priority based on p_cpu and p_nice
	Nice      int8     // p_nice - Process "nice" value
	Pgrp      pointer // p_pgrp
	Xstat     uint16  // p_xstat - Exit status for wait; also stop signal
	Acflag    uint16  // p_acflag - Accounting flags
*/

type RInfo struct {
	Data  []byte
	ru *C.struct_rusage_info_v5
}

func RInfoData() RInfo {
	data := make([]byte, C.sizeof_struct_rusage_info_v5)
	for i := range data {
		data[i] = byte(i)
	}
	return RInfo{
		Data:  data,
		ru: (*C.struct_rusage_info_v5)(unsafe.Pointer(&data[0])),
	}
}

func (r RInfo) Uuid() [16]byte                  {
	var uuid [16]byte
	for i := range uuid {
		uuid[i] = byte(r.ru.ri_uuid[i])
	}
	return uuid
}
func (r RInfo) UserTime() uint64                { return uint64(r.ru.ri_user_time) }
func (r RInfo) SystemTime() uint64              { return uint64(r.ru.ri_system_time) }
func (r RInfo) PkgIdleWkups() uint64            { return uint64(r.ru.ri_pkg_idle_wkups) }
func (r RInfo) InterruptWkups() uint64          { return uint64(r.ru.ri_interrupt_wkups) }
func (r RInfo) Pageins() uint64                 { return uint64(r.ru.ri_pageins) }
func (r RInfo) WiredSize() uint64               { return uint64(r.ru.ri_wired_size) }
func (r RInfo) ResidentSize() uint64            { return uint64(r.ru.ri_resident_size) }
func (r RInfo) PhysFootprint() uint64           { return uint64(r.ru.ri_phys_footprint) }
func (r RInfo) ProcStartAbstime() uint64        { return uint64(r.ru.ri_proc_start_abstime) }
func (r RInfo) ProcExitAbstime() uint64         { return uint64(r.ru.ri_proc_exit_abstime) }
func (r RInfo) ChildUserTime() uint64           { return uint64(r.ru.ri_child_user_time) }
func (r RInfo) ChildSystemTime() uint64         { return uint64(r.ru.ri_child_system_time) }
func (r RInfo) ChildPkgIdleWkups() uint64       { return uint64(r.ru.ri_child_pkg_idle_wkups) }
func (r RInfo) ChildInterruptWkups() uint64     { return uint64(r.ru.ri_child_interrupt_wkups) }
func (r RInfo) ChildPageins() uint64            { return uint64(r.ru.ri_child_pageins) }
func (r RInfo) ChildElapsedAbstime() uint64     { return uint64(r.ru.ri_child_elapsed_abstime) }
func (r RInfo) DiskioBytesread() uint64         { return uint64(r.ru.ri_diskio_bytesread) }
func (r RInfo) DiskioByteswritten() uint64      { return uint64(r.ru.ri_diskio_byteswritten) }
func (r RInfo) CpuTimeQosDefault() uint64       { return uint64(r.ru.ri_cpu_time_qos_default) }
func (r RInfo) CpuTimeQosMaintenance() uint64   { return uint64(r.ru.ri_cpu_time_qos_maintenance) }
func (r RInfo) CpuTimeQosBackground() uint64    { return uint64(r.ru.ri_cpu_time_qos_background) }
func (r RInfo) CpuTimeQosUtility() uint64       { return uint64(r.ru.ri_cpu_time_qos_utility) }
func (r RInfo) CpuTimeQosLegacy() uint64        { return uint64(r.ru.ri_cpu_time_qos_legacy) }
func (r RInfo) CpuTimeQosUserInitiated() uint64 { return uint64(r.ru.ri_cpu_time_qos_user_initiated) }
func (r RInfo) CpuTimeQosUserInteractive() uint64 {
	return uint64(r.ru.ri_cpu_time_qos_user_interactive)
}
func (r RInfo) BilledSystemTime() uint64         { return uint64(r.ru.ri_billed_system_time) }
func (r RInfo) ServicedSystemTime() uint64       { return uint64(r.ru.ri_serviced_system_time) }
func (r RInfo) LogicalWrites() uint64            { return uint64(r.ru.ri_logical_writes) }
func (r RInfo) LifetimeMaxPhysFootprint() uint64 { return uint64(r.ru.ri_lifetime_max_phys_footprint) }
func (r RInfo) Instructions() uint64             { return uint64(r.ru.ri_instructions) }
func (r RInfo) Cycles() uint64                   { return uint64(r.ru.ri_cycles) }
func (r RInfo) BilledEnergy() uint64             { return uint64(r.ru.ri_billed_energy) }
func (r RInfo) ServicedEnergy() uint64           { return uint64(r.ru.ri_serviced_energy) }
func (r RInfo) IntervalMaxPhysFootprint() uint64 { return uint64(r.ru.ri_interval_max_phys_footprint) }
func (r RInfo) RunnableTime() uint64             { return uint64(r.ru.ri_runnable_time) }
func (r RInfo) Flags() uint64                    { return uint64(r.ru.ri_flags) }
