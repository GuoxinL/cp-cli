/**
  Create by guoxin 2021.06.07
*/
package main

import (
	"fmt"
	"strings"
	"time"
)

type CommandError struct {
	error
	err  string
	code int
}

func (c CommandError) Error() string {
	return fmt.Sprintf("code: %v, error: %v", c.code, c.err)
}

func (c CommandError) Console() {
	Console(c.err)
}

// Console 输出
func Console(s string) {
	fmt.Println(CMDPrefix + " " + s)
}

func ConsoleError(s string, code int) {
	commandError := CommandError{
		err:  s,
		code: code,
	}
	panic(commandError)
}

func empty(s string) bool {
	return len(s) == 0
}

func getTime() string {
	t := time.Now()
	return t.Format("20060102150405")
}

func lines(str string) []string {
	if empty(str) {
		return []string{}
	}
	return strings.Split(str, "\n")
}
