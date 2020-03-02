package standard

import (
	"context"
	"io"
	"os"
	"os/exec"
	"sync"
)

// Exec starts the named process with the given args, and streams STDOUT/STDERR
// back to the caller.
func Exec(ctx context.Context, name string, arg ...string) error {
	cmd := exec.CommandContext(ctx, name, arg...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	// Launch the process.
	if err := cmd.Start(); err != nil {
		return err
	}

	var errStdout, errStderr error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		// Copy from process STDOUT to our STDOUT.
		defer wg.Done()
		_, errStdout = io.Copy(os.Stdout, stdout)
	}()

	wg.Add(1)
	go func() {
		// Copy from process STDERR to our STDERR.
		defer wg.Done()
		_, errStderr = io.Copy(os.Stderr, stderr)
	}()

	// Wait for the process and output copying to finish.
	err := cmd.Wait()
	wg.Wait()

	// Return any encountered errors.
	switch {
	case err != nil:
		return err

	case errStdout != nil:
		return errStdout

	case errStderr != nil:
		return errStderr

	default:
		return nil
	}
}
