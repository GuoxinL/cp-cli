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
		PanicError(`请检查当前仓库状态，保证当前工作区没有修改`, 1)
	}
	currentBranch = status.Branch
	// 检出源分支
	if Checkout(targetBranch) {
		PanicError(`检出源分支异常，目标分支 `+targetBranch+` 不存在`, 1)
	}
	// 创建并且检出分支 {source branch}_{keyword}_{target-branch}_{yyyyMMdd}
	newBranch := fmt.Sprintf("%v_to_%v_%v_%v", sourceBranch, targetBranch, keyword, getTime())
	if CreateBranch(newBranch) {
		PanicError(`创建分支异常，新分支 `+newBranch, 1)
	}
	// 检出新目标分支
	if Checkout(newBranch) {
		PanicError(fmt.Sprintf("检出新分支 %v 异常", newBranch), 1)
	}
	commits := getCommits(sourceBranch, keyword)
	if len(commits) == 0 {
		PanicError(fmt.Sprintf("检出新分支 %v 异常", newBranch), 1)
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

}