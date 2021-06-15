/**
  Create by guoxin 2021.06.07
*/
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// ConflictsString 异常 cherry-pick
const ConflictsString = "error: could not apply"

type Commit struct {
	Id          string
	Time        time.Time
	Name        string
	Description string
}

type Status struct {
	Clean  bool
	Branch string
}

//CMDWrapper 执行命令并输出日志 See CMD
func CMDWrapper(command string, stdin, stdout, stderr io.Writer) (result string, err error) {
	fmt.Println()
	Console("command: " + command)
	result, err = CMD(command, stdin, stdout, stderr)
	Console("result: \n" + result)
	if err != nil {
		Console("err: " + err.Error())
	}
	Console("command: " + command + " done")
	return
}

// CMD 执行命令
func CMD(command string, stdin, stdout, stderr io.Writer) (result string, err error) {
	comm := exec.Command(`bash`, `-c`, command)
	out := bytes.Buffer{}
	err0 := bytes.Buffer{}
	comm.Stdout = &out
	comm.Stderr = &err0
	isError := false
	var message = ""
	err = comm.Start()
	if err != nil {
		isError = true
		result = err.Error()
		return
	}
	if err = comm.Wait(); err != nil {
		isError = true
		result = err0.String()
		return
	}

	header := fmt.Sprintf(`pid: %v, exit code: %v`,
		comm.ProcessState.Pid(),
		comm.ProcessState.Sys().(syscall.WaitStatus).ExitStatus())
	Console(header)
	message = out.String()
	if isError {
		message = err0.String()
	} else {
		message = out.String()
	}
	result = message
	return
}

func Checkout(branch string) bool {

	_, err := CMDWrapper(`git checkout `+branch, nil, &bytes.Buffer{}, &bytes.Buffer{})
	return err != nil
}

func DeleteBranch(branch string) bool {
	_, err := CMDWrapper(`git branch -d `+branch, nil, &bytes.Buffer{}, &bytes.Buffer{})
	return err != nil

}

// CreateBranch git branch {branch name}
func CreateBranch(branch string) bool {
	_, err := CMDWrapper(`git branch `+branch, nil, &bytes.Buffer{}, &bytes.Buffer{})
	return err != nil
}

func GitStatus() Status {
	const nothingToCommit = "nothing to commit, working tree clean"
	result, err := CMDWrapper(`git status`, nil, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		os.Exit(0)
	}
	lines := lines(result)
	split := strings.Split(lines[0], " ")
	clean := strings.TrimSpace(lines[1]) == nothingToCommit
	status := Status{Clean: clean, Branch: split[2]}
	return status
}

func getCommits(sourceBranch, keyword string) []Commit {
	result, err := CMDWrapper(`git log `+sourceBranch+` --oneline --reverse --pretty=format:"%h|%ad|%an|%s" --date=format:"%Y-%m-%d %H:%M:%S" --grep="`+keyword+`"`, nil, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		ConsoleError(err.Error(), 1)
	}
	lines := lines(result)
	if len(lines) == 0 {
		ConsoleError("未检索到符合 keyword: "+keyword+"的 commit", 1)
	}
	var commits = make([]Commit, len(lines)-1)
	lines = lines[:len(lines)-1]
	for i, line := range lines {
		commitLine := strings.Split(line, "|")
		parse, _ := time.Parse("2006-01-02 15:04:05", commitLine[1])
		commit := &Commit{}
		commit.Id = commitLine[0]
		commit.Time = parse
		commit.Name = commitLine[2]
		commit.Description = commitLine[3]
		commits[i] = *commit
	}
	return commits
}

var xxx syscall.WaitStatus = 0x100

func CherryPick(commitId string) bool {
	result, err := CMDWrapper(`git cherry-pick `+commitId, nil, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok {
			if exitError.ProcessState.Sys() == xxx {
				Console("合并%v过程中发生冲突", 1, commitId)
			} else {
				ConsoleError(err.Error(), 1)
			}
		} else {
			ConsoleError(err.Error(), 1)
		}
	}

	lines := lines(result)
	lines = lines[:len(lines)-1]
	return !strings.HasPrefix(lines[0], ConflictsString)
}

func CherryPickContinue() bool {
	_, err := CMDWrapper(`git cherry-pick --continue`, nil, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		ConsoleError(err.Error(), 1)
	}
	return false
}

func CherryPickAbort() bool {
	_, err := CMDWrapper(`git cherry-pick --abort`, nil, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		ConsoleError(err.Error(), 1)
	}
	return true
}

func AddAll() {
	_, err := CMDWrapper(`git add .`, nil, &bytes.Buffer{}, &bytes.Buffer{})
	if err != nil {
		ConsoleError(err.Error(), 1)
	}
}

func GitCommit() {
	_, err := CMDWrapper(`git commit`, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		ConsoleError(err.Error(), 1)
	}
}
