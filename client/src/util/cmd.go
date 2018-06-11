package util

import (
	"errors"
	"fmt"
	"os/exec"
	"time"
)

// Cmd wraps the stdlib function exec.Command to log the command and command
// result
func Cmd(name string, args ...string) (string, error) {
	printCmd(name, args...)

	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()

	outString := string(out)

	fmt.Println(outString)

	return outString, err
}

// CmdWithTimeoutOptions is a struct to be passwd to CmdWithTimeout
type CmdWithTimeoutOptions struct {
	Timeout   time.Duration
	LogOutput bool
}

// CmdWithTimeout is similar to Cmd, but allows to specify a timeout, and will
// kill the command if it doesn't exit before the timeout is reached
func CmdWithTimeout(
	options CmdWithTimeoutOptions,
	name string,
	args ...string) (string, error) {
	type result struct {
		out string
		err error
	}

	c := make(chan result, 1)
	cmd := exec.Command(name, args...)

	go func() {
		if options.LogOutput {
			printCmd(name, args...)
		}

		out, err := cmd.CombinedOutput()
		outString := string(out)

		if options.LogOutput {
			fmt.Println(outString)
		}

		c <- result{outString, err}
	}()

	select {
	case res := <-c:
		return res.out, res.err
	case <-time.After(options.Timeout):
		cmd.Process.Kill()
		return "", errors.New("Command timed out")
	}
}

func printCmd(name string, args ...string) {
	fmt.Printf("> " + name)

	for _, arg := range args {
		fmt.Printf(" " + arg)
	}

	fmt.Printf("\n")
}
