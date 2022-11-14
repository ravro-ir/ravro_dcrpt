package core

import (
	"os"
	"ravro_dcrpt/entity"
	"ravro_dcrpt/utils"
	"runtime"
	"strings"
)

func DcrptAmendment(currentPath, keyFixPath, outFixPath string) ([]string, error) {
	var amendment entity.Amendment
	var (
		path    string
		err     error
		lstMore []string
	)
	if currentPath == "" {
		path, err = utils.Projectpath()
		if err != nil {
			return lstMore, err
		}
	}

	lstAmendment, _ := utils.WalkMatch(path, "*.ravro")
	for _, name := range lstAmendment {
		if !strings.Contains(name, "amendment-") {
			continue
		}
		Process, err := utils.FileProccessing(name)
		if err != nil {
			return lstMore, err
		}
		if runtime.GOOS == "windows" {
			_, err = utils.SslDecrypt(Process.Name, outFixPath+"\\"+Process.Filename, keyFixPath)
		} else {
			_, err = utils.SslDecrypt(Process.Name, outFixPath+"/"+Process.Filename, keyFixPath)
		}
		if err != nil {
			return lstMore, err
		}
		process := utils.CheckPlatform(outFixPath, Process)
		err = os.Rename(process.NewNamePath, process.OldNamePath)
		if err != nil {
			return lstMore, err
		}
		if strings.Index(process.OldName, "data") != 0 {
			continue
		}
		_, err = utils.JsonParser(process, &amendment)
		if err = os.Remove(process.OldNamePath); err != nil {
			return lstMore, err
		}
		lstMore = append(lstMore, amendment.Description)
	}
	return lstMore, nil
}
