package loger

import (
	"log"
	"time"
)
import "os"

var (
	file, _ = os.OpenFile("/home/kryvonos/go/BunkerBot/BunkerBot/Bot/App/src/loger/log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	logFile = log.New(file, "", 0)
)

func LogErr(err error) {
	if err != nil {
		logFile.Print(err)
		logFile.Println(time.Now().String())
	}
}
