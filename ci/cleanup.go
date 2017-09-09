package ci

import "os"

func cleanupWorkspace(workdir string) error {
	return os.RemoveAll(workdir)
}
