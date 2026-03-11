package chroot

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
)

// chrootContext holds the file descriptors needed to escape back to the real root
type chrootContext struct {
	realRootFd int // FD of the real root directory (escape hatch)
	realCwdFd  int // FD of the real working directory
}

// enterChroot pivots into the given chroot directory.
// Returns a ChrootContext that can be used to return to the original root.
func enterChroot(newRoot string) (*chrootContext, error) {
	// Lock this goroutine to the current OS thread — chroot is per-thread on Linux
	runtime.LockOSThread()

	// Open a FD to the real root BEFORE chrooting (escape hatch)
	realRootFd, err := syscall.Open("/", syscall.O_RDONLY|syscall.O_DIRECTORY, 0)
	if err != nil {
		runtime.UnlockOSThread()
		return nil, fmt.Errorf("open real root: %w", err)
	}

	// Open a FD to the real CWD
	realCwdFd, err := syscall.Open(".", syscall.O_RDONLY|syscall.O_DIRECTORY, 0)
	if err != nil {
		syscall.Close(realRootFd)
		runtime.UnlockOSThread()
		return nil, fmt.Errorf("open real cwd: %w", err)
	}

	// Chroot into the new root
	if err := syscall.Chroot(newRoot); err != nil {
		syscall.Close(realRootFd)
		syscall.Close(realCwdFd)
		runtime.UnlockOSThread()
		return nil, fmt.Errorf("chroot to %q: %w", newRoot, err)
	}

	// Change directory to "/" inside the chroot so CWD is consistent
	if err := syscall.Chdir("/"); err != nil {
		// Best-effort escape back before returning the error
		_ = syscall.Fchdir(realRootFd)
		_ = syscall.Chroot(".")
		syscall.Close(realRootFd)
		syscall.Close(realCwdFd)
		runtime.UnlockOSThread()
		return nil, fmt.Errorf("chdir inside chroot: %w", err)
	}

	return &chrootContext{
		realRootFd: realRootFd,
		realCwdFd:  realCwdFd,
	}, nil
}

// exit restores the original root and working directory.
func (c *chrootContext) exit() error {
	defer runtime.UnlockOSThread()
	defer syscall.Close(c.realRootFd)
	defer syscall.Close(c.realCwdFd)

	// fchdir into the real root FD, then chroot(".") to break out
	if err := syscall.Fchdir(c.realRootFd); err != nil {
		return fmt.Errorf("fchdir to real root fd: %w", err)
	}
	if err := syscall.Chroot("."); err != nil {
		return fmt.Errorf("chroot('.') to escape: %w", err)
	}
	// Restore original working directory
	if err := syscall.Fchdir(c.realCwdFd); err != nil {
		return fmt.Errorf("fchdir to real cwd: %w", err)
	}
	return nil
}

// RunInChroot is a convenience wrapper: enters chroot, calls fn, then exits.
func RunInChroot(newRoot string, fn func() error) error {
	ctx, err := enterChroot(newRoot)
	if err != nil {
		return fmt.Errorf("enter chroot: %w", err)
	}
	// Always exit the chroot, even if fn panics
	defer func() {
		if exitErr := ctx.exit(); exitErr != nil {
			fmt.Fprintf(os.Stderr, "WARNING: failed to exit chroot: %v\n", exitErr)
		}
	}()
	return fn()
}
