package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/Kodik77rus/fs_syncer/internal/pkg/config"
	"github.com/Kodik77rus/fs_syncer/internal/pkg/utils"
)

var logLvlMp map[string]slog.Level = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"error": slog.LevelError,
	"info":  slog.LevelInfo,
}

func InitLogger(conf *config.Config) error {
	if err := utils.CheckOrCreateFile(conf.LogFilePath); err != nil {
		return err
	}
	f, err := os.OpenFile(
		conf.LogFilePath,
		os.O_APPEND|os.O_RDWR,
		os.ModePerm,
	)
	if err != nil {
		return err
	}
	var h slog.Handler
	logLvl := parseLogLvl(conf.LogLvl)
	if conf.StdoutLogEnable {
		wrt := io.MultiWriter(os.Stdout, f)
		h = slog.NewJSONHandler(wrt, &slog.HandlerOptions{
			Level: logLvl,
		})
	} else {
		h = slog.NewJSONHandler(f, &slog.HandlerOptions{
			Level: logLvl,
		})
	}
	logger := slog.New(h)
	slog.SetDefault(logger)
	return nil
}

func parseLogLvl(logLvl string) slog.Level {
	lvl, ok := logLvlMp[logLvl]
	if !ok {
		return slog.LevelDebug
	}
	return lvl
}
