package public

const (
	MAX_SIZE = 100    	//每个日志文件保存的最大尺寸 单位：M
	MAX_BACKUPS = 10  	//日志文件最多保存多少个备份
	MAX_AGE = 7			//文件最多保存多少天
	LETTER_BYTES = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	REDIS_EXPTIME_DEFAULT = 300
)