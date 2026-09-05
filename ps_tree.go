package ps

// A ProcessMap is a map of processes and their children
type ProcessMap struct {
	Pids     map[int]*Process
	Children map[int][]int
}

// GetProcessMap returns a process map of all processes in the system.
// Additional information for each process is included including the
// Process.Children slice.
func GetProcessMap() (*ProcessMap, error) {
	procs, err := Processes(true)
	if err != nil {
		return nil, err
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
	return pm, nil
}

// GetChildren returns the list of PIDs of the direct children of the process
// specified by pid.  This is equivalent to GetProcessMap().GetChildren(pid).
func GetChildren(pid int) []int {
	pm, err := GetProcessMap()
	if err != nil {
		return nil
	}
	return pm.GetChildren(pid)
}

// GetDescendants returns the list of PIDs of all descendants of the process
// specified by pid.  This is equivalent to GetProcessMap().GetDescendants(pid).
func GetDescendants(pid int) []int {
	pm, err := GetProcessMap()
	if err != nil {
		return nil
	}
	return pm.GetDescendants(pid)
}

// GetDecendents is a synonym for GetDescendants.
func GetDecendents(pid int) []int {
	return GetDescendants(pid)
}

// GetChildren returns the list of PIDs of the direct children of the process
// specified by pid.
func (pm *ProcessMap) GetChildren(pid int) []int {
	if pm == nil {
		return nil
	}
	return pm.Children[pid]
}

// GetDescendants returns the list of PIDs of all descendants of the process
// specified by pid.
func (pm *ProcessMap) GetDescendants(pid int) []int {
	if pm == nil {
		return nil
	}
	return pm.appendDescendants(pid, nil, map[int]bool{pid: true})
}

// GetDecendents is a synonym for GetDescendants.
func (pm *ProcessMap) GetDecendents(pid int) []int {
	return pm.GetDescendants(pid)
}

func (pm *ProcessMap) appendDescendants(pid int, dst []int, seen map[int]bool) []int {
	for _, child := range pm.Children[pid] {
		if seen[child] {
			continue
		}
		seen[child] = true
		dst = append(dst, child)
		dst = pm.appendDescendants(child, dst, seen)
	}
	return dst
}
