package ci

import "os/exec"

func runTests(workdir string) error {
	args := []string{"test"}
	cmd := exec.Command("make", args...)
	cmd.Dir = workdir

	out, err := cmd.Output()
	if err != nil {
		log().Errorf("Unable run tests: %v", err)
		return err
	}

	log().Debugf("%s tests are passed. %s", workdir, out)

	return nil
}
