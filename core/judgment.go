package core

import (
	"errors"
	"log"
	"os"
	"ravro_dcrpt/entity"
	"ravro_dcrpt/utils"
	"runtime"
	"strings"
)

func DcrptJudgment(currentPath, keyFixPath, outFixpath string, checkStatus bool) (entity.Judgment, error) {
	var judgment entity.Judgment
	var (
		err      error
		lstJudge []string
	)
	lstJudge, err = utils.WalkMatch(currentPath, "*.ravro")
	if err != nil {
		return judgment, err
	}
	if len(lstJudge) == 0 {
		return judgment, errors.New("judge is empty")
	}
	//}
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
		Process, err := utils.FileProccessing(name)
		if err != nil {
			return judgment, err
		}
		if runtime.GOOS == "windows" {
			_, err = utils.SslDecrypt(Process.Name, outFixpath+"\\"+Process.Filename, keyFixPath)
		} else {
			_, err = utils.SslDecrypt(Process.Name, outFixpath+"/"+Process.Filename, keyFixPath)
		}
		if err != nil {
			return judgment, err
		}
		process := utils.CheckPlatform(outFixpath, Process)
		err = os.Rename(process.NewNamePath, process.OldNamePath)
		if err != nil {
			return judgment, err
		}
		_, err = utils.JsonParser(process, &judgment)
		if err = os.Remove(process.OldNamePath); err != nil {
			log.Fatal(err)
		}
	}
	return judgment, nil
}
