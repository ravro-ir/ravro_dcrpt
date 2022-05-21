package main

import (
	"log"
	"os"
	"runtime"
	"strings"
)

type RavroVulnerability struct {
	Name    string `json:"name"`
	Define  string `json:"define"`
	Fix     string `json:"fix"`
	Writeup string `json:"writeup"`
}

type JudgeCvss struct {
	Value  string `json:"value"`
	Rating string `json:"rating"`
}

type Judgment struct {
	Score         int    `json:"score"`
	Reward        int    `json:"reward"`
	Description   string `json:"description"`
	Cvss          JudgeCvss
	Vulnerability RavroVulnerability
}

func DcrptJudgment(currentPath, keyFixPath, outFixpath string, checkStatus bool) (Judgment, error) {
	var judgment Judgment
	var (
		path     string
		err      error
		lstJudge []string
	)
	if currentPath == "" {
		path, err = projectpath()
		if err != nil {
			return judgment, err
		}
		lstJudge, err = WalkMatch(path, "*.ravro")
		if err != nil {
			return judgment, err
		}
		if len(lstJudge) == 0 {
			return judgment, err
		}
	} else {
		lstJudge, err = WalkMatch(currentPath, "*.ravro")
		if err != nil {
			return judgment, err
		}
		if len(lstJudge) == 0 {
			return judgment, err
		}
	}
	for _, name := range lstJudge {
		if runtime.GOOS == "windows" {
			if !checkStatus {
				if !strings.Contains(name, "\\encrypt\\") {
					continue
				}
			}
			if !strings.Contains(name, "\\judgment\\") {
				continue
			}
		} else {
			if !strings.Contains(name, "/encrypt/") {
				continue
			}
			if !strings.Contains(name, "/judgment/") {
				continue
			}
		}
		Process, err := fileProccessing(name)
		if err != nil {
			return judgment, err
		}
		if runtime.GOOS == "windows" {
			_, err = SslDecrypt(Process.name, outFixpath+"\\"+Process.filename, keyFixPath)
		} else {
			_, err = SslDecrypt(Process.name, outFixpath+"/"+Process.filename, keyFixPath)
		}
		if err != nil {
			return judgment, err
		}
		process := CheckPlatform(outFixpath, Process)
		err = os.Rename(process.newNamePath, process.oldNamePath)
		if err != nil {
			return judgment, err
		}
		_, err = JsonParser(process, &judgment)
		if err = os.Remove(process.oldNamePath); err != nil {
			log.Fatal(err)
		}
	}
	return judgment, nil
}
