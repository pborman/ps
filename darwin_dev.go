//go:build darwin

package ps

import (
	"fmt"
	"os"
	"sync"
	"syscall"
)

// A DevT is a unix device number
type DevT uint32
const noDev = 0xffffffff

func (d DevT) Major() int {
	if d == noDev {
		return -1
	}
	return int((d >> 24) & 0xff)
}

func (d DevT) Minor() int {
	if d == noDev {
		return -1
	}
	return int(d & 0xffffff)
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
		if (stat.Mode & (syscall.S_IFCHR|syscall.S_IFBLK)) != 0 {
			devNames[DevT(stat.Rdev)] = de.Name()
		}
	}

	return devNames
}
