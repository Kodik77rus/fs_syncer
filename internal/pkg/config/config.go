package config

import "flag"

type Config struct {
	StdoutLogEnable  bool
	SrcFolderPath    string
	OutputFolderPath string
	LogFilePath      string
	LogLvl           string
}

func InitConfig() *Config {
	var (
		stdoutLogEnable  bool
		srcFolderPath    string
		outputFolderPath string
		logFilePath      string
		logLvl           string
	)
	flag.BoolVar(&stdoutLogEnable, "stdout-log-enable", false, "log to stdout")
	flag.StringVar(&srcFolderPath, "src-folder-path", "./src", "path to source dir")
	flag.StringVar(&outputFolderPath, "output-folder-path", "./out", "path to output folder")
	flag.StringVar(&logFilePath, "log-file-path", "./fs-watcher.log", "path to log file")
	flag.StringVar(&logLvl, "log-lvl", "debug", "app log lvl")
	flag.Parse()
	return &Config{
		StdoutLogEnable:  stdoutLogEnable,
		SrcFolderPath:    srcFolderPath,
		OutputFolderPath: outputFolderPath,
		LogFilePath:      logFilePath,
		LogLvl:           logLvl,
	}
}
