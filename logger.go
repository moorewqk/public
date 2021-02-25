package public

import (
	"fmt"
	"go.uber.org/zap"
	_ "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path/filepath"
)


/*
Filename:   filePath,   // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
*/
func  InitLogger(appname string) (logger *zap.SugaredLogger) {
	logDir,err := CreateDir()
	if err!=nil{
		fmt.Printf("日志目录操作失败:%v",err)
		panic(err)
	}
	info_log := appname + ".log"
	info_log = filepath.Join(logDir,info_log)
	logger = NewLogger(info_log, zapcore.InfoLevel, MAX_SIZE, MAX_BACKUPS, MAX_AGE, true, "Main")
	return
}