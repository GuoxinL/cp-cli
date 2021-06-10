/**
  Create by guoxin 2021.06.02
*/
package main

import (
	"flag"
)

const CMDPrefix = "[CMD]"

// cp-cli origin/dev master TOERP-0
// 1. 检出目标分支
// 2. 将源分支的commit cherry-pick到{target-branch}_{keyword}_{yyyyMMdd}
// 3. 完成合并

func main() {
	// 参数
	// 参数
	// 1. source branch:
	// 代码来源从source branch 检索代码
	// 2. target branch:
	// 目标分支
	// 2. keyword:
	// 通过 keyword 从 source branch中检索出来 commit
	//
	Process(getParameters())
}

func getParameters() (string, string, string) {
	Console(`Parameters`)
	var sourceBranch = flag.String("s", "", "source branch:\t代码来源从 source branch 检索代码")
	var targetBranch = flag.String("t", "", "target branch:\t代码来源从 target branch 检索代码")
	var keyword = flag.String("k", "", "keyword:\t通过 keyword 从 source branch中检索出来 commit")

	flag.Parse()
	if empty(*sourceBranch) {
		PanicError(`source branch is required`, 1)
	} else {
		Console(`source branch: ` + *sourceBranch)
	}
	if empty(*targetBranch) {
		PanicError(`target branch is required`, 1)
	} else {
		Console(`target branch: ` + *targetBranch)
	}
	if empty(*keyword) {
		PanicError(`keyword is required`, 1)
	} else {
		Console(`keyword: ` + *keyword)
	}
	return *sourceBranch, *targetBranch, *keyword
}
