package public

import (
	"os"
	"path/filepath"
	"time"
)

func CreateDir() (logDir string,err error) {
	dir,err := os.Getwd()
	if err!=nil{
		return
	}
	logDir = filepath.Join(dir,"logs")
	_,err =os.Lstat(logDir)
	if err != nil{
		err = os.MkdirAll(logDir,0755)
		if err!=nil{
			return
		}
	}
	return
}

//转换过期时间
func  MakeTime(t int) (dt time.Duration) {
	return time.Duration(t)*time.Second
}