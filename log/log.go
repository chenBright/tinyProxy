package log

import (
	"log"
	"path/filepath"
	"os"
)

const LogDir string = "logs"

// 初始化
func Init(name string) {
	filename := filepath.Join(LogDir, name)
	_ = os.MkdirAll(LogDir, os.ModePerm) // 新建文件夹和文件
	logFile, err := os.OpenFile(filename, os.O_CREATE | os.O_RDWR | os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Println("cannot create log file: ", err)
		return
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("success create log file: ", logFile.Name())
}
