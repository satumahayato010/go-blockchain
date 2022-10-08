package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

// ConfigList config.iniファイルの情報をgoで扱うために格納する構造体を定義。
type ConfigList struct {
	ApiKey    string
	ApiSecret string
	LogFile   string
}

// Config ConfigList構造体として定義しておく。
var Config ConfigList

// init関数は、main関数より先に呼び出される。
func init() {
	// ini.Loadメソッドは、引数に渡したファイルを読み込み解析をする。
	// 読み込んだ内容をFile構造体に格納し、返り値として返す。
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Panicf("Failed to read file %v", err)
		// プログラムを指定されたコードで終了させる。０が成功、それ以外はエラー。
		os.Exit(1)
	}
	// 先に指定した, ConfigList型のConfig変数に値を格納する。
	Config = ConfigList{
		// Sectionメソッドは、引数のセクション（部分）を解析して、中身をSection構造体として返す
		// keyで、引数に指定したセクションの値を解析してKey構造体に格納して返す。Stringで文字列として返している。
		ApiKey:    cfg.Section("bitflyer").Key("api_key").String(),
		ApiSecret: cfg.Section("bitflyer").Key("api_secret").String(),
		LogFile:   cfg.Section("gotrading").Key("log_file").String(),
	}
}
