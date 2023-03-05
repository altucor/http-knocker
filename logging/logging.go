package logging

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// var (
// 	DebugLogger   *log.Entry
// 	InfoLogger    *log.Entry
// 	WarningLogger *log.Entry
// 	ErrorLogger   *log.Entry
// )

func init() {
	// TODO: Make possible to set logging to file or to stdout
	// TODO: Find more comfortable log formatter
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	// log.SetLevel(log.ErrorLevel)
	// file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.SetOutput(file)
	log.SetOutput(os.Stdout)

	// DebugLogger = log.WithFields(log.Fields{
	// 	"timestamp": "test",
	// 	"lol":       "aaa",
	// })
}

/*
Loggin needs:
timestamp level fileName functionName

*/

func CommonLog() *log.Entry {
	// pc, fullFileName, lineNum, _ := runtime.Caller(1)
	// funcNameE := runtime.FuncForPC(pc).Name()
	// _, fileName := filepath.Split(fullFileName)
	return log.WithFields(log.Fields{
		// "class": className,
		// "func": funcNameE,
		// "file": fileName + "@" + fmt.Sprint(lineNum),
	})
}
