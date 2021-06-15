/**
  Create by guoxin 2021.06.07
*/
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var newBranch string
var currentBranch string

func Process(sourceBranch, targetBranch, keyword string) {
	defer func() {
		if err := recover(); err != nil {
			commandError, ok := err.(CommandError)
			if ok {
				commandError.Console()
			} else {
				err2, ok := err.(error)
				if ok {
					Console(err2.Error())
				} else {
					fmt.Println(err)
				}
			}
			Checkout(currentBranch)
			DeleteBranch(newBranch)
			os.Exit(commandError.code)
		}
	}()
	// 检查仓库状态
	status := GitStatus()
	if !status.Clean {
		ConsoleError(`请检查当前仓库状态，保证当前工作区没有修改`, 1)
	}
	currentBranch = status.Branch
	// 检出源分支
	if Checkout(targetBranch) {
		ConsoleError(`检出源分支异常，目标分支 %v 不存在`, 1, targetBranch)
	}
	// 创建并且检出分支 {source branch}_{keyword}_{target-branch}_{yyyyMMdd}
	newBranch = fmt.Sprintf("%v_to_%v_%v_%v", sourceBranch, targetBranch, keyword, getTime())
	if CreateBranch(newBranch) {
		ConsoleError(`创建分支异常，新分支 %v`, 1, newBranch)
	}
	// 检出新目标分支
	if Checkout(newBranch) {
		ConsoleError("检出新分支 %v 异常", 1, newBranch)
	}
	commits := getCommits(sourceBranch, keyword)
	if len(commits) == 0 {
		ConsoleError("检出新分支 %v 异常", 1, newBranch)
	}
	for _, commit := range commits {
		if CherryPick(commit.Id) {
			Console("Cherry-Pick 遇到冲突请处理...")
			Console("处理完成请输入：continue")
			Console("退出请输入：abort")
			stdinReader(StdFuncs{
				"continue": CherryPickContinue,
				"abort":    CherryPickAbort,
			})
			// git add .
			AddAll()
			// git commit
			GitCommit()
		}
	}
	marshal, _ := json.Marshal(commits)
	Console(string(marshal))
	// 阶段
	// 1. 获得 commit 信息
	// 2. 编排 commit-id
	// 3. 执行 cherry-pick 操作
	// git log origin/dev --oneline --reverse | grep "TOERP-4180"
	// git log master --oneline --reverse | grep "TOERP-0"
	//`git log master --oneline --reverse --grep="TOERP-0"`
	//
	ConsoleError("结束", 1)

}
