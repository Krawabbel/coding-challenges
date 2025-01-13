package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"golang.org/x/sys/unix"
)

type shell struct {
	history []string
	histPtr int
}

func newShell() *shell {
	return &shell{
		history: make([]string, 0),
	}
}

func (s *shell) deleteLastNChars(n int) {
	if n > 0 {
		s.printf("\033[%dD\033[K", n)
	}
}

func (s *shell) print(args ...any) {
	str := fmt.Sprint(args...)
	_, err := io.WriteString(os.Stdout, str)
	must(err)
}

func (s *shell) println(args ...any) {
	s.print(args...)
	s.print("\n")
}

func (s *shell) printf(format string, args ...any) {
	str := fmt.Sprintf(format, args...)
	s.print(str)
}

func (s *shell) printErr(err error) {
	s.println(err.Error())
}

var builtins = map[string]func(*shell, []string) error{
	"cd": cd,
}

type tooManyArgumentsError struct {
	name       string
	want, have int
}

func (e tooManyArgumentsError) Error() string {
	return fmt.Sprintf("'%s' error: too many arguments (want %d, have %d)", e.name, e.want, e.have)
}

func cd(s *shell, args []string) error {
	if len(args) > 1 {
		return tooManyArgumentsError{"cd", 1, len(args)}
	}
	dir := args[0]
	return os.Chdir(dir)
}

func (s *shell) run(line string) error {
	s.histAdd(line)
	return s.runLine(line)
}

func (s *shell) runLine(line string) error {
	cmds := strings.Split(line, "|")

	var stdin io.Reader = os.Stdin
	var wg sync.WaitGroup
	for i, c := range cmds {
		fmt.Fprint(os.Stderr, i, c, "\n")

		wg.Add(1)
		r, w := io.Pipe()

		var stdout io.Writer
		if i == len(cmds)-1 {
			stdout = os.Stdout
		} else {
			stdout = bufio.NewWriter(w)
		}

		go func() {
			defer wg.Done()
			defer w.Close()
			if err := s.runCmd(c, stdin, stdout); err != nil {
				s.printErr(err)
			}
		}()

		stdin = bufio.NewReader(r)
	}

	wg.Wait()
	return nil
}

func (s *shell) runCmd(command string, stdin io.Reader, stdout io.Writer) error {
	parts := strings.Split(strings.TrimSpace(command), " ")

	name := parts[0]
	args := parts[1:]

	builtin, found := builtins[name]
	if found {
		return builtin(s, args)
	}

	cmd := exec.Command(name, args...)
	if cmd.Err != nil {
		return cmd.Err
	}

	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func match(have []byte, want ...byte) bool {
	if len(have) != len(want) {
		return false
	}
	for i := range have {
		if have[i] != want[i] {
			return false
		}
	}
	return true
}

func (s *shell) escape(line string, esc []byte) (string, error) {
	switch {
	case match(esc, 27, 91, 65):
		return s.histGet(-1), nil
	case match(esc, 27, 91, 66):
		return s.histGet(+1), nil
	}
	return "", fmt.Errorf("unexpected escape sequence %v", esc)
}

func (s *shell) histAdd(line string) {
	s.history = append(s.history, line)
	s.histPtr = len(s.history)
}

func (s *shell) histGet(offset int) string {
	s.histPtr = max(0, min(len(s.history), s.histPtr+offset))
	if s.histPtr == len(s.history) {
		return ""
	}
	return s.history[s.histPtr]
}

func (s *shell) prompt() (string, error) {
	s.print("ccsh> ")
	line := ""
	buf := make([]byte, 1024)
	for {
		s.print(line)
		l := len(line)

		n, err := os.Stdin.Read(buf)
		if err != nil {
			return "", err
		}

		if n == 0 {
			return "", fmt.Errorf("err: no input")
		}

		c := buf[0]
		switch c {
		case 27: // ESCAPE
			next, err := s.escape(line, buf[:n])
			if err != nil {
				return "", err
			}
			line = next
		case '\n':
			s.println()
			return line, nil
		case '\t':
			return "", fmt.Errorf("tab")
		case '\x7f':
			if l > 0 {
				line = line[0 : l-1]
			} else {
				s.print("\a")
			}
		default:
			line += string(buf[:n])
		}

		s.deleteLastNChars(l)
	}
}

func (s *shell) repl() error {

	for {
		line, err := s.prompt()
		if err != nil {
			return err
		}

		if line == "exit" {
			return nil
		}

		if err := s.run(line); err != nil {
			s.printErr(err)
		}
	}
}

func runShell() error {
	s := newShell()
	return s.repl()
}

func main() {
	fd := int(os.Stdin.Fd())
	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	must(err)

	lflag := termios.Lflag
	cc := termios.Cc

	termios.Lflag &^= unix.ICANON | unix.ECHO
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0
	must(unix.IoctlSetTermios(fd, unix.TCSETS, termios))

	defer func() {
		termios.Lflag = lflag
		termios.Cc = cc
		must(unix.IoctlSetTermios(fd, unix.TCSETS, termios))
	}()

	must(runShell())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
