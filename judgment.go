package main

import (
	"log"
	"os"
	"runtime"
	"strings"
)

type JudgeCvss struct {
	Value  string `json:"value"`
	Rating string `json:"rating"`
}

type Judgment struct {
	Score       int    `json:"score"`
	Reward      int    `json:"reward"`
	Description string `json:"description"`
	Cvss        JudgeCvss
}

func DcrptJudgment() (Judgment, error) {
	var judgment Judgment
	path, err := projectpath()
	if err != nil {
		return judgment, err
	}
	lstJudge, _ := WalkMatch(path, "*.ravro")
	for _, name := range lstJudge {
		if runtime.GOOS == "windows" {
			if !strings.Contains(name, "\\judgment\\") {
				continue
			}
		} else {
			if !strings.Contains(name, "/judgment/") {
				continue
			}
		}
		Process, err := fileProccessing(name)
		if err != nil {
			return judgment, err
		}
		_, err = SslDecrypt(name, Process.filename)
		if err != nil {
			return judgment, err
		}
		process := CheckPlatform(Process)
		err = os.Rename(process.newNamePath, process.oldNamePath)
		if err != nil {
			return judgment, err
		}
		_, err = JsonParser(process, &judgment)
		if err = os.Remove(process.oldNamePath); err != nil {
			log.Fatal(err)
		}
		return judgment, nil
	}
	return judgment, nil
}
