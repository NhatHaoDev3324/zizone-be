package factory

import (
	"log"

	"github.com/NhatHaoDev3324/goAuth/constant"
)

func LogSuccess(msg string) {
	log.Println(constant.Green+"[OK]"+constant.ResetAll, msg)
}

func LogInfo(msg string) {
	log.Println(constant.Blue+"[INFO]"+constant.ResetAll, msg)
}

func LogError(msg string) {
	log.Println(constant.Red+"[ERROR]"+constant.ResetAll, msg)
}

func LogWarn(msg string) {
	log.Println(constant.Yellow+"[WARN]"+constant.ResetAll, msg)
}
