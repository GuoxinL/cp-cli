/**
  Create by guoxin 2021.06.07
*/
package main

import (
	"bufio"
	"fmt"
	"os"
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
func Console(s string, a ...interface{}) {
	str := ""
	if len(a) != 0 {
		str = fmt.Sprintf(s, a)
	} else {
		str = s
	}
	fmt.Println(CMDPrefix + " " + str)
}

func ConsoleError(s string, code int, a ...interface{}) {
	str := ""
	if len(a) != 0 {
		str = fmt.Sprintf(s, a)
	} else {
		str = s
	}
	commandError := CommandError{
		err:  str,
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

type StdFuncs map[string]func() bool

func stdinReader(stdFuncs StdFuncs) bool {
	var keywords []string
	for keyword, _ := range stdFuncs {
		keywords = append(keywords, keyword)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		in, _ := reader.ReadString('\n')
		// convert CRLF to LF
		in = strings.TrimSpace(strings.Replace(in, "\n", "", -1))

		for keyword, f := range stdFuncs {
			if in == keyword {
				return f()
			}
		}
	}
}
