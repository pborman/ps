package ps

// A ProcessMap is a map of processes and their children
type ProcessMap struct {
	Pids     map[int]*Process
	Children map[int][]int
}

// GetProcessMap returns a process map of all processes in the system.
// Additional information for each process is included including the
// Process.Children slice.
func GetProcessMap() *ProcessMap {
	procs, err := Processes(true)
	if err != nil {
		return nil
	}
	pm := &ProcessMap{
		Pids:     map[int]*Process{},
		Children: map[int][]int{},
	}
	for _, p := range procs {
		pm.Pids[p.ID] = p
		ppid, err := p.Ppid()
		if err == nil && ppid != 0 {
			pm.Children[ppid] = append(pm.Children[ppid], p.ID)
		}
	}
	for _, p := range procs {
		ppid, err := p.Ppid()
		if err != nil || ppid == 0 {
			continue
		}
		pp := pm.Pids[ppid]
		if pp == nil {
			continue
		}
		pp.Children = append(pp.Children, p)
	}
	return pm
}

// GetChildren returns the list of PIDs of the direct children of the process
// specified by pid.  This is equvialent  to GetProcessMap().GetChildren(pid).
func GetChildren(pid int) []int {
	return GetProcessMap().GetChildren(pid)
}

// GetChildren returns the list of PIDs of all decendents of the process
// specified by pid.  This is equvialent to GetProcessMap().GetDecendents(pid).
func GetDecendents(pid int) []int {
	return GetProcessMap().GetDecendents(pid)
}

// GetChildren returns the list of PIDs of the direct children of the process
// specified by pid.
func (pm *ProcessMap) GetChildren(pid int) []int {
	if pm == nil {
		return nil
	}
	return pm.Children[pid]
}

// GetChildren returns the list of PIDs of all decendents of the process
// specified by pid.
func (pm *ProcessMap) GetDecendents(pid int) []int {
	if pm == nil {
		return nil
	}
	p := pm.Pids[pid]
	if p == nil {
		return nil
	}
	var children []int
	children = append(children, pm.Children[pid]...)
	for _, child := range pm.Children[pid] {
		children = append(children, pm.GetDecendents(child)...)
	}
	return children
}
