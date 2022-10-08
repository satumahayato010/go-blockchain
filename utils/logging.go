package utils

import (
	"io"
	"log"
	"os"
)

func LoggingSettings(logFile string) {
	// os.OpenFile(ファイル名, フラグ, ファイルモード)返り値に、File構造体とエラー返る。
	logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("file=logFile err=%s", err.Error())
	}
	// io.MultiWriter() は引数で与えられたすべての io.Writer に対しての Write を行うような io.Writer を返します。
	//以下の例の場合は、 os.Stdout と fileの両方に書き込むようになります。
	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	// 出力する、フラグをしている。日付、時刻、ファイル名と行番号。
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// setOutputメソッドは、ログの出力先を設定する。引数に出力先を指定する。(io.writer型）
	log.SetOutput(multiLogFile)
}
