package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
	// never reached
	panic(true)
	return nil, nil
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: ding [command] [args]")
		return
	}
	args := make([]string, 0)
	for index, a := range os.Args {
		if index >= 2 {
			args = append(args, a)
		}
	}
	cmd := exec.Command(os.Args[1], args...)
	var stdout, stderr []byte
	var errStdout, errStderr error
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	cmd.Start()

	go func() {
		stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
	}()

	go func() {
		stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)
	}()

	cmd.Wait()
	ring()
}

func ring() {
	fmt.Print("\a")
}
