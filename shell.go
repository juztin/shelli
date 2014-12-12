package shelli

import (
	"io"

	"code.google.com/p/go.crypto/ssh"
)

// Stream the io.Reader to a channel, passing any errors on the given error channel.
func streamToChan(stream io.Reader, ch chan<- []byte, errCh chan<- error) {
	//buf := make([]byte, 1<<10)
	for {
		buf := make([]byte, 1<<10)
		n, err := stream.Read(buf[0:])
		if n > 0 {
			/*c := make([]byte, n)
			copy(c, buf)
			ch <- c*/
			ch <- buf[0:n]
		}
		if err != nil {
			errCh <- err
		}
	}
}

// Get stdin, stdout, and stderr from the SSH session.
func sessionStreams(session *ssh.Session) (stdin io.WriteCloser, stdout io.Reader, stderr io.Reader, err error) {
	stdin, err = session.StdinPipe()
	if err != nil {
		return
	}
	stdout, err = session.StdoutPipe()
	if err != nil {
		return
	}
	stderr, err = session.StderrPipe()
	if err != nil {
		return
	}
	return
}

// NewShell creates a new SSH session to the client using inCh as stdin, and quiting on a receive to quit.
// Returns stdout, stderr channels and error channel
func NewPty(client *ssh.Client, inCh <-chan []byte, quit <-chan struct{}) (<-chan []byte, <-chan []byte, <-chan error) {
	errCh := make(chan error, 1)
	// Create a new session from the given configuration func.
	session, err := client.NewSession()
	if err != nil {
		errCh <- err
		return nil, nil, nil
	}
	// Get session std in/out/err
	stdin, stdout, stderr, err := sessionStreams(session)
	if err != nil {
		errCh <- err
		return nil, nil, nil
	}
	// Stream from readers to channels
	stdoutCh := make(chan []byte)
	stderrCh := make(chan []byte)
	go streamToChan(stdout, stdoutCh, errCh)
	go streamToChan(stderr, stderrCh, errCh)
	// Set terminal modes
	modes := ssh.TerminalModes{
		//ssh.ECHO:     0,
		ssh.VREPRINT: 0,
		//ssh.CS8:      1,
		//ssh.ECHOCTL: 0,
	}
	// Request pseudo terminal
	err = session.RequestPty("xterm", 40, 120, modes)
	//err = session.RequestPty("vt100", 24, 80, modes)
	if err != nil {
		errCh <- err
		return nil, nil, nil
	}
	// Invoke SSH shell
	err = session.Shell()
	if err != nil {
		errCh <- err
		return nil, nil, nil
	}
	// Listen for input/quit
	go func() {
		for {
			select {
			case b := <-inCh:
				stdin.Write(b)
			case <-quit:
				session.Close()
				return
			}
		}
	}()
	return stdoutCh, stderrCh, errCh
}
