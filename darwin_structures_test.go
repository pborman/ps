//go:build darwin

package ps

import (
	"testing"
	"unsafe"

	data "github.com/pborman/ps/darwin_test_data"
)

func TestKInfoStructure(t *testing.T) {
	raw := data.KInfoData()
	ki, err := mkKInfoProc(raw.Data)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the Go and C structures, pointing to the same
	// data, have the same values.
	if got, want := ki.Flag, raw.Flag(); got != want {
		t.Errorf("Flag() got %v want %v", got, want)
	}
	if got, want := ki.Stat, Stat(raw.Stat()); got != want {
		t.Errorf("Stat() got %v want %v", got, want)
	}
	if got, want := ki.Pid, raw.Pid(); got != want {
		t.Errorf("Pid() got %v want %v", got, want)
	}
	if got, want := ki.Debugger, raw.Debugger(); got != want {
		t.Errorf("Debugger() got %v want %v", got, want)
	}
	if got, want := ki.Sigwait, raw.Sigwait(); got != want {
		t.Errorf("Sigwait() got %v want %v", got, want)
	}
	if got, want := ki.Estcpu, raw.Estcpu(); got != want {
		t.Errorf("Estcpu() got %v want %v", got, want)
	}
	if got, want := ki.Cpticks, raw.Cpticks(); got != want {
		t.Errorf("Cpticks() got %v want %v", got, want)
	}
	if got, want := ki.Pctcpu, Fixpt(raw.Pctcpu()); got != want {
		t.Errorf("Pctcpu() got %v want %v", got, want)
	}
	if got, want := ki.Swtime, raw.Swtime(); got != want {
		t.Errorf("Swtime() got %v want %v", got, want)
	}
	if got, want := ki.Slptime, raw.Slptime(); got != want {
		t.Errorf("Slptime() got %v want %v", got, want)
	}
	if got, want := ki.Uticks, raw.Uticks(); got != want {
		t.Errorf("Uticks() got %v want %v", got, want)
	}
	if got, want := ki.Sticks, raw.Sticks(); got != want {
		t.Errorf("Sticks() got %v want %v", got, want)
	}
	if got, want := ki.Iticks, raw.Iticks(); got != want {
		t.Errorf("Iticks() got %v want %v", got, want)
	}
	if got, want := ki.Traceflag, raw.Traceflag(); got != want {
		t.Errorf("Traceflag() got %v want %v", got, want)
	}
	if got, want := ki.Holdcnt, raw.Holdcnt(); got != want {
		t.Errorf("Holdcnt() got %v want %v", got, want)
	}
	if got, want := ki.Sigignore, raw.Sigignore(); got != want {
		t.Errorf("Sigignore() got %v want %v", got, want)
	}
	if got, want := ki.Sigcatch, raw.Sigcatch(); got != want {
		t.Errorf("Sigcatch() got %v want %v", got, want)
	}
	if got, want := ki.Priority, raw.Priority(); got != want {
		t.Errorf("Priority() got %v want %v", got, want)
	}
	if got, want := ki.Usrpri, raw.Usrpri(); got != want {
		t.Errorf("Usrpri() got %v want %v", got, want)
	}
	if got, want := ki.Nice, raw.Nice(); got != want {
		t.Errorf("Nice() got %v want %v", got, want)
	}
	if got, want := ki.Xstat, raw.Xstat(); got != want {
		t.Errorf("Xstat() got %v want %v", got, want)
	}
	if got, want := ki.Acflag, raw.Acflag(); got != want {
		t.Errorf("Acflag() got %v want %v", got, want)
	}

	if got, want := ki.Uid, raw.Uid(); got != want {
		t.Errorf("Uid got %v, want %v", got, want)
	}
	if got, want := ki.Svuid, raw.Svuid(); got != want {
		t.Errorf("Svuid got %v, want %v", got, want)
	}
	if got, want := ki.Gid, raw.Gid(); got != want {
		t.Errorf("Gid got %v, want %v", got, want)
	}
	if got, want := ki.Svgid, raw.Svgid(); got != want {
		t.Errorf("Svgid got %v, want %v", got, want)
	}
	if got, want := ki.Refcnt, raw.Refcnt(); got != want {
		t.Errorf("Refcnt got %v, want %v", got, want)
	}
	if got, want := ki.Ref, raw.Ref(); got != want {
		t.Errorf("Ref got %v, want %v", got, want)
	}
	if got, want := ki.Euid, raw.Euid(); got != want {
		t.Errorf("Euid got %v, want %v", got, want)
	}
	if got, want := ki.Ppid, raw.Ppid(); got != want {
		t.Errorf("Ppid got %v, want %v", got, want)
	}
	if got, want := ki.Pgid, raw.Pgid(); got != want {
		t.Errorf("Pgid got %v, want %v", got, want)
	}
	if got, want := ki.Tpgid, raw.Tpgid(); got != want {
		t.Errorf("Tpgid got %v, want %v", got, want)
	}
	if got, want := ki.Tdev, DevT(raw.Tdev()); got != want {
		t.Errorf("Tdev got %v, want %v", got, want)
	}
	if got, want := ki.Xsize, raw.Xsize(); got != want {
		t.Errorf("Xsize got %v, want %v", got, want)
	}
	if got, want := ki.Eflag, raw.Eflag(); got != want {
		t.Errorf("Eflag got %v, want %v", got, want)
	}
	if got, want := ki.Xrssize, raw.Xrssize(); got != want {
		t.Errorf("Xrssize got %v, want %v", got, want)
	}
	if got, want := ki.Xccount, raw.Xccount(); got != want {
		t.Errorf("Xccount got %v, want %v", got, want)
	}
	if got, want := ki.Xswrss, raw.Xswrss(); got != want {
		t.Errorf("Xswrss got %v, want %v", got, want)
	}
	if got, want := ki.Ngroups, raw.Ngroups(); got != want {
		t.Errorf("Ngroups got %v, want %v", got, want)
	}
	if got, want := ki.Jobc, raw.Jobc(); got != want {
		t.Errorf("Jobc got %v, want %v", got, want)
	}
	if got, want := ki.Groups, raw.Groups(); got != want {
		t.Errorf("Groups got %v, want %v", got, want)
	}
	if got, want := ki.Comm, raw.Comm(); got != want {
		t.Errorf("Comm got %v, want %v", got, want)
	}
}

func TestRUsageStructure(t *testing.T) {
	raw := data.RInfoData()
	var ru = (*RUsage)(unsafe.Pointer(&raw.Data[0]))
	if got, want := ru.Uuid, raw.Uuid(); got != want {
		t.Errorf("Uuid got %v, want %v", got, want)
	}
	if got, want := ru.UserTime, raw.UserTime(); got != want {
		t.Errorf("UserTime got %v, want %v", got, want)
	}
	if got, want := ru.SystemTime, raw.SystemTime(); got != want {
		t.Errorf("SystemTime got %v, want %v", got, want)
	}
	if got, want := ru.PkgIdleWkups, raw.PkgIdleWkups(); got != want {
		t.Errorf("PkgIdleWkups got %v, want %v", got, want)
	}
	if got, want := ru.InterruptWkups, raw.InterruptWkups(); got != want {
		t.Errorf("InterruptWkups got %v, want %v", got, want)
	}
	if got, want := ru.Pageins, raw.Pageins(); got != want {
		t.Errorf("Pageins got %v, want %v", got, want)
	}
	if got, want := ru.WiredSize, raw.WiredSize(); got != want {
		t.Errorf("WiredSize got %v, want %v", got, want)
	}
	if got, want := ru.ResidentSize, raw.ResidentSize(); got != want {
		t.Errorf("ResidentSize got %v, want %v", got, want)
	}
	if got, want := ru.PhysFootprint, raw.PhysFootprint(); got != want {
		t.Errorf("PhysFootprint got %v, want %v", got, want)
	}
	if got, want := ru.ProcStartAbstime, raw.ProcStartAbstime(); got != want {
		t.Errorf("ProcStartAbstime got %v, want %v", got, want)
	}
	if got, want := ru.ProcExitAbstime, raw.ProcExitAbstime(); got != want {
		t.Errorf("ProcExitAbstime got %v, want %v", got, want)
	}
	if got, want := ru.ChildUserTime, raw.ChildUserTime(); got != want {
		t.Errorf("ChildUserTime got %v, want %v", got, want)
	}
	if got, want := ru.ChildSystemTime, raw.ChildSystemTime(); got != want {
		t.Errorf("ChildSystemTime got %v, want %v", got, want)
	}
	if got, want := ru.ChildPkgIdleWkups, raw.ChildPkgIdleWkups(); got != want {
		t.Errorf("ChildPkgIdleWkups got %v, want %v", got, want)
	}
	if got, want := ru.ChildInterruptWkups, raw.ChildInterruptWkups(); got != want {
		t.Errorf("ChildInterruptWkups got %v, want %v", got, want)
	}
	if got, want := ru.ChildPageins, raw.ChildPageins(); got != want {
		t.Errorf("ChildPageins got %v, want %v", got, want)
	}
	if got, want := ru.ChildElapsedAbstime, raw.ChildElapsedAbstime(); got != want {
		t.Errorf("ChildElapsedAbstime got %v, want %v", got, want)
	}
	if got, want := ru.DiskioBytesread, raw.DiskioBytesread(); got != want {
		t.Errorf("DiskioBytesread got %v, want %v", got, want)
	}
	if got, want := ru.DiskioByteswritten, raw.DiskioByteswritten(); got != want {
		t.Errorf("DiskioByteswritten got %v, want %v", got, want)
	}
	if got, want := ru.CpuTimeQosDefault, raw.CpuTimeQosDefault(); got != want {
		t.Errorf("CpuTimeQosDefault got %v, want %v", got, want)
	}
	if got, want := ru.CpuTimeQosMaintenance, raw.CpuTimeQosMaintenance(); got != want {
		t.Errorf("CpuTimeQosMaintenance got %v, want %v", got, want)
	}
	if got, want := ru.CpuTimeQosBackground, raw.CpuTimeQosBackground(); got != want {
		t.Errorf("CpuTimeQosBackground got %v, want %v", got, want)
	}
	if got, want := ru.CpuTimeQosUtility, raw.CpuTimeQosUtility(); got != want {
		t.Errorf("CpuTimeQosUtility got %v, want %v", got, want)
	}
	if got, want := ru.CpuTimeQosLegacy, raw.CpuTimeQosLegacy(); got != want {
		t.Errorf("CpuTimeQosLegacy got %v, want %v", got, want)
	}
	if got, want := ru.CpuTimeQosUserInitiated, raw.CpuTimeQosUserInitiated(); got != want {
		t.Errorf("CpuTimeQosUserInitiated got %v, want %v", got, want)
	}
	if got, want := ru.CpuTimeQosUserInteractive, raw.CpuTimeQosUserInteractive(); got != want {
		t.Errorf("CpuTimeQosUserInteractive got %v, want %v", got, want)
	}
	if got, want := ru.BilledSystemTime, raw.BilledSystemTime(); got != want {
		t.Errorf("BilledSystemTime got %v, want %v", got, want)
	}
	if got, want := ru.ServicedSystemTime, raw.ServicedSystemTime(); got != want {
		t.Errorf("ServicedSystemTime got %v, want %v", got, want)
	}
	if got, want := ru.LogicalWrites, raw.LogicalWrites(); got != want {
		t.Errorf("LogicalWrites got %v, want %v", got, want)
	}
	if got, want := ru.LifetimeMaxPhysFootprint, raw.LifetimeMaxPhysFootprint(); got != want {
		t.Errorf("LifetimeMaxPhysFootprint got %v, want %v", got, want)
	}
	if got, want := ru.Instructions, raw.Instructions(); got != want {
		t.Errorf("Instructions got %v, want %v", got, want)
	}
	if got, want := ru.Cycles, raw.Cycles(); got != want {
		t.Errorf("Cycles got %v, want %v", got, want)
	}
	if got, want := ru.BilledEnergy, raw.BilledEnergy(); got != want {
		t.Errorf("BilledEnergy got %v, want %v", got, want)
	}
	if got, want := ru.ServicedEnergy, raw.ServicedEnergy(); got != want {
		t.Errorf("ServicedEnergy got %v, want %v", got, want)
	}
	if got, want := ru.IntervalMaxPhysFootprint, raw.IntervalMaxPhysFootprint(); got != want {
		t.Errorf("IntervalMaxPhysFootprint got %v, want %v", got, want)
	}
	if got, want := ru.RunnableTime, raw.RunnableTime(); got != want {
		t.Errorf("RunnableTime got %v, want %v", got, want)
	}
	if got, want := ru.Flags, raw.Flags(); got != want {
		t.Errorf("Flags got %v, want %v", got, want)
	}

}
