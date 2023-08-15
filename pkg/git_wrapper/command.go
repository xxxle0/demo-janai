package git_wrapper

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"operarius/internal/utils"
	"os/exec"
	"strings"

	"github.com/guardrailsio/go-scan-helper/logger"
	helperlog "github.com/guardrailsio/go-scan-helper/logger"
)

//go:generate mockgen -source=./command.go -destination=../../mock/pkg/git_wrapper/command.go -package=mock

type CommandBuilder struct {
	baseCommand     string
	dir             string
	command         string
	args            []string
	baseCommandArgs []string
	logger          logger.ILogger
}

type ICommandBuilder interface {
	AddCommand(string)
	AddArg(string)
	AddBaseCommandArg(string)
	AddBaseCommandArgs([]string)
	AddArgs([]string)
	SetDir(dir string)
	SetLogger(logger.ILogger)
	Build() string
	Exec() (string, error)
	ExecStreamStdout() (string, error)
	ExecCommandPath(commandPath string, cb func(*exec.Cmd)) error
}

func NewCommandBuilder() ICommandBuilder {
	return &CommandBuilder{
		baseCommand: "git",
	}
}

func (c *CommandBuilder) SetLogger(logger helperlog.ILogger) {
	c.logger = logger
}

func (c *CommandBuilder) AddCommand(command string) {
	c.command = command
}

func (c *CommandBuilder) AddArgs(args []string) {
	c.args = append(c.args, args...)
}

func (c *CommandBuilder) AddArg(arg string) {
	c.args = append(c.args, arg)
}

func (c *CommandBuilder) AddBaseCommandArg(arg string) {
	c.baseCommandArgs = append(c.baseCommandArgs, arg)
}

func (c *CommandBuilder) AddBaseCommandArgs(args []string) {
	c.baseCommandArgs = append(c.baseCommandArgs, args...)
}

func (c *CommandBuilder) SetDir(dir string) {
	c.dir = dir
}

func (c *CommandBuilder) Build() string {
	args := strings.Trim(fmt.Sprintf("%s %s", c.command, strings.Join(c.args, " ")), " ")
	return fmt.Sprintf("%s %s", c.baseCommand, args)
}

func (c *CommandBuilder) Exec() (string, error) {
	var errb bytes.Buffer
	args := c.baseCommandArgs
	args = append(args, c.command)
	args = append(args, c.args...)
	if c.logger != nil {
		c.logger.Debugf("Exec at %s: Command = %s, Arguments = %v", c.dir, c.baseCommand, args)
	} else {
		log.Printf("Exec at %s: Command = %s, Arguments = %v", c.dir, c.baseCommand, args)
	}
	cmd := exec.Command(c.baseCommand, args...)
	cmd.Stderr = &errb
	if c.dir != "" {
		cmd.Dir = c.dir
	}
	stdout, err := cmd.Output()
	c.Reset()
	if err != nil {
		log.Println("err: ", strings.TrimSpace(errb.String()))
		return "", errors.New(strings.TrimSpace(errb.String()))
	}
	cmd.Wait()
	return string(stdout), nil
}

func (c *CommandBuilder) Reset() {
	c.baseCommand = ""
	c.dir = ""
	c.command = ""
	c.args = []string{}
}

func (c *CommandBuilder) ExecStreamStdout() (string, error) {
	args := []string{c.command}
	args = append(args, c.baseCommandArgs...)
	args = append(args, c.args...)
	cmd := exec.Command(c.baseCommand, args...)
	if c.logger != nil {
		c.logger.Debugf("Exec at %s: Command = %s, Arguments = %v", c.dir, c.baseCommand, args)
	} else {
		log.Printf("Exec at %s: Command = %s, Arguments = %v", c.dir, c.baseCommand, args)
	}
	if c.dir != "" {
		cmd.Dir = c.dir
	}
	stderr, _ := cmd.StderrPipe()
	cmd.Start()
	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Print(m)
	}
	cmd.Wait()
	output, err := cmd.Output()
	c.Reset()
	return string(output), err
}

func (c *CommandBuilder) ExecCommandPath(commandPath string, cb func(*exec.Cmd)) error {
	args := []string{strings.ToLower(c.baseCommand)}
	args = append(args, c.baseCommandArgs...)
	args = append(args, strings.ToLower(c.command))
	args = append(args, c.args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := &exec.Cmd{
		Path:   commandPath,
		Args:   args,
		Stdout: &stdout,
		Stderr: &stderr,
	}

	if cb != nil {
		cb(cmd)
	}
	if c.logger != nil {
		c.logger.Debugf("Exec at %s: Command = %s, Arguments = %v", c.dir, c.baseCommand, args)
	} else {
		log.Printf("Exec at %s: Command = %s, Arguments = %v", c.dir, c.baseCommand, args)
	}
	err := cmd.Run()
	c.Reset()
	if err != nil {
		if strings.Contains(stderr.String(), "fatal: fetch-pack") ||
			strings.Contains(stderr.String(), "fatal: early EOF") {
			return utils.ErrCloneFailedDueToLackofMemory
		} else {
			errMsg := stderr.String()
			p := strings.Split(errMsg, "\n")
			if len(p) >= 2 {
				errMsg = p[1]
				errMsg = strings.TrimSpace(errMsg)
				if errMsg == "" {
					errMsg = stderr.String()
				}
			}
			return fmt.Errorf("error: %s", errMsg)
		}
	}
	return nil
}
