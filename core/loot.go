package core

import (
	"fmt"
	"os"
	"syscall"

	"github.com/muraenateam/necrobrowser/log"
)

// CheckLoot checks the path is writable
func CheckLoot(path string) (isWritable bool, err error) {

	isWritable = false
	info, err := os.Stat(path)
	if err != nil {
		log.Debug("Loot path doesn't exist, creating")
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return
		}
	}

	info, err = os.Stat(path)
	if !info.IsDir() {
		err = fmt.Errorf("Loot path isn't a directory")
		return
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		err = fmt.Errorf("Write permission bit is not set on this file for user")
		return
	}

	var stat syscall.Stat_t
	if err = syscall.Stat(path, &stat); err != nil {
		err = fmt.Errorf("Unable to get stat")
		return
	}

	err = nil
	if uint32(os.Geteuid()) != stat.Uid {
		err = fmt.Errorf("User doesn't have permission to write to loot directory (%s)", path)
		return
	}

	isWritable = true
	return
}
