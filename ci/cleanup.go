package ci

import "os"

func cleanupWorkspace(workdir string) error {
	return os.Remove(workdir)
}
