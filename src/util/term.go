package term

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type TermColor string

const (
	Red      TermColor = "31"
	Green    TermColor = "32"
	Yellow   TermColor = "33"
	White    TermColor = "37"
	Cyan     TermColor = "36"
	Blue     TermColor = "34"
	Orange   TermColor = "35"
	DarkGray TermColor = "90"
)

func AskForInt(prompt string) int {
	var i int
	fmt.Println(prompt)
	fmt.Scanf("%d", &i)

	return i
}

func Colorize(text string, color TermColor) string {
	return fmt.Sprintf("\033[%sm%s\033[0m", color, text)
}

func ReadChar() byte {
	var c byte
	fmt.Scanf("%c", &c)

	return c
}

// setTerminalToRaw sets the terminal to raw mode.
func SetTerminalToRaw() error {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, os.Stdin.Fd(), syscall.TIOCGETA, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	if err != 0 {
		return err
	}

	termios.Lflag &^= syscall.ICANON | syscall.ECHO

	_, _, err = syscall.Syscall6(syscall.SYS_IOCTL, os.Stdin.Fd(), syscall.TIOCSETA, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	if err != 0 {
		return err
	}

	return nil
}

// restoreTerminal restores the terminal to its original mode.
func RestoreTerminal() {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, os.Stdin.Fd(), syscall.TIOCGETA, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	if err != 0 {
		fmt.Println("Error getting terminal settings:", err)
		return
	}

	termios.Lflag |= syscall.ICANON | syscall.ECHO

	_, _, err = syscall.Syscall6(syscall.SYS_IOCTL, os.Stdin.Fd(), syscall.TIOCSETA, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	if err != 0 {
		fmt.Println("Error restoring terminal settings:", err)
		return
	}
}

func GetCharStr() (string, int, error) {
	var char [1]byte
	n, err := os.Stdin.Read(char[:])
	if err != nil {
		return "0", n, err
	}
	return string(rune(char[0])), n, nil
}

func Clear(rows, cols int) {
	fmt.Printf("\033[%dA", rows)
	fmt.Printf("\033[%dD", cols)
}

func ClearLastTerminalLine() {

}
