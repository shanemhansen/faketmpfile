// Package faketmp is a library that turns an io.Reader into an *os.File which
// has a name in the linux proc filesystem. The primary use case is to pass in memory buffers or readers to legacy apps which require filename arguments.
package faketmpfile

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"strconv"
)

type FakeTempFile struct {
	*os.File
	errorChan chan error
}

// ExtraFiles ensures this FakeTempFile is passed to the cmd
func (this *FakeTempFile) ExtraFiles(cmd *exec.Cmd) {
	cmd.ExtraFiles = append(cmd.ExtraFiles, this.File)
}

// Name returns the name of the file in the linux proc filesystem.
// This probably does not work on mac/windows.
func (this *FakeTempFile) Name() string {
	return fmt.Sprintf("/proc/self/fd/" + strconv.Itoa(int(this.Fd())))
}

// Close closes the underlying io.Pipe and then waits for the status of the
// Copy
func (this *FakeTempFile) Close() error {
	this.File.Close()
	return <-this.errorChan
}

//FakeTemp returns a path which can be opened and read the contents of rdr.
//It also returns a filehandle, suitable for appending to the ExtraFiles attribute
func FakeTemp(rdr io.Reader) (*FakeTempFile, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create pipe")
	}
	errChan := make(chan error, 1)
	go func() {
		_, err := io.Copy(w, rdr)
		w.Close()
		errChan <- err
	}()
	return &FakeTempFile{File: r, errorChan: errChan}, nil
}
