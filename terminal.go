package shelli

/*import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

// TerminalChan pumps from a terminal reader to the returned channel.
func TerminalChan() <-chan []byte {
	// Disable input buffering.
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// Do not display entered characters on the screen.
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	ch := make(chan []byte)
	go func() {
		defer exec.Command("stty", "-F", "/dev/tty", "echo")
		reader := bufio.NewReader(os.Stdin)
		for {
			b, err := reader.ReadByte()
			if err != nil {
				fmt.Println(err)
				return
			}
			ch <- []byte{b}
		}
	}()
	return ch
}*/
