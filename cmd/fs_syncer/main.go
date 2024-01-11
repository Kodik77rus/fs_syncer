package main

import (
	"fmt"
	"log/slog"
	"os"

	fs_syncer "gitlab.rebrainme.com/golang_users_repos/5266/fs-watcher/internal/app/fs-syncer"
	"gitlab.rebrainme.com/golang_users_repos/5266/fs-watcher/internal/pkg/config"
	dir_watcher "gitlab.rebrainme.com/golang_users_repos/5266/fs-watcher/internal/pkg/dir-watcher"
	"gitlab.rebrainme.com/golang_users_repos/5266/fs-watcher/internal/pkg/logger"
	"gitlab.rebrainme.com/golang_users_repos/5266/fs-watcher/internal/pkg/utils"
)

func main() {
	if err := run(); err != nil {
		slog.Error(
			"main shutdown",
			slog.String(
				"error",
				err.Error(),
			),
		)
		os.Exit(1)
	}
}

func run() error {
	config := config.InitConfig()
	if err := utils.CheckOrCreateFile(config.LogFilePath); err != nil {
		return fmt.Errorf("failed check or create log file: %v", err)
	}
	if err := utils.CheckOrCreateFolders(
		config.SrcFolderPath,
		config.OutputFolderPath,
	); err != nil {
		return err
	}
	if err := logger.InitLogger(config); err != nil {
		return err
	}
	dirWatcher := dir_watcher.InitDirWatcher(config)
	fsSyncer := fs_syncer.InitFsWatcher(config, dirWatcher)
	return fsSyncer.Start()
}
