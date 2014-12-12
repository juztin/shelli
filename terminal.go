package shelli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

// TerminalChan pumps from a terminal reader to the returned channel.
//func TerminalChan() <-chan byte {
func TerminalChan() <-chan []byte {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	//ch := make(chan byte)
	ch := make(chan []byte)
	go func() {
		defer exec.Command("stty", "-F", "/dev/tty", "echo")
		//var b []byte = make([]byte, 1)
		reader := bufio.NewReader(os.Stdin)
		for {
			//b, err := reader.ReadBytes('\n')
			/*for i := range b {
				in <- b[i]
			}*/
			/*i, err := reader.Read(b)
			if err != nil {
				fmt.Println(err)
				return
			}
			if i == 1 {
				ch <- []byte{b[0]}
			}*/
			b, err := reader.ReadByte()
			if err != nil {
				fmt.Println(err)
				return
			}
			ch <- []byte{b}
		}
	}()
	return ch
}
