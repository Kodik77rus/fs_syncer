package fs_syncer

import (
	"io/fs"
	"log/slog"

	"github.com/Kodik77rus/fs_syncer/internal/pkg/config"
	"github.com/Kodik77rus/fs_syncer/internal/pkg/constants"
	dir_watcher "github.com/Kodik77rus/fs_syncer/internal/pkg/dir-watcher"
	"github.com/Kodik77rus/fs_syncer/internal/pkg/utils"
)

type FsSyncer struct {
	dirWatcher       *dir_watcher.DirWatcher
	outputFolderPath string
	srcFolderPath    string
}

func InitFsWatcher(
	config *config.Config,
	dirWatcher *dir_watcher.DirWatcher,
) *FsSyncer {
	return &FsSyncer{
		dirWatcher:       dirWatcher,
		outputFolderPath: config.OutputFolderPath,
		srcFolderPath:    config.SrcFolderPath,
	}
}

func (f *FsSyncer) Start() error {
	watcherChan := f.dirWatcher.Start()
	for fsEvent := range watcherChan {
		if err := fsEvent.Error; err != nil {
			return err
		}
		newEventLog(fsEvent)
		outPutFolderEntityPath := utils.JoinOsPath(f.outputFolderPath, fsEvent.OsPath)
		fsEventName := fsEvent.EventName
		switch fsEventName {
		case constants.Delete:
			if err := utils.OsRemove(outPutFolderEntityPath); err != nil {
				errorDeleteLog(outPutFolderEntityPath, err)
				continue
			}
			successDeleteLog(outPutFolderEntityPath)
		case constants.Create:
			if err := utils.CreateFile(outPutFolderEntityPath); err != nil {
				errorCreateLog(outPutFolderEntityPath, err)
				continue
			}
			successCreateLog(outPutFolderEntityPath)
		case constants.Write:
			srcFilePath := utils.JoinOsPath(f.srcFolderPath, fsEvent.OsPath)
			if err := utils.CopyFileContent(srcFilePath, outPutFolderEntityPath); err != nil {
				failedSyncFileContentLog(srcFilePath, outPutFolderEntityPath)
				continue
			}
			successfulSyncFileContentLog(srcFilePath, outPutFolderEntityPath)
		case constants.ChangeAttrib:
			srcFilePath := utils.JoinOsPath(f.srcFolderPath, fsEvent.OsPath)
			currentEntityInfo, err := utils.GetOsStat(srcFilePath)
			if err != nil {
				failedGetOsEntityStatLog(srcFilePath, err)
				continue
			}
			lastEntityInfo, err := utils.GetOsStat(outPutFolderEntityPath)
			if err != nil {
				failedGetOsEntityStatLog(outPutFolderEntityPath, err)
				continue
			}
			if err := utils.ChangeOsPerm(
				outPutFolderEntityPath,
				currentEntityInfo.Mode().Perm(),
			); err != nil {
				failedChangeEntityPermLog(outPutFolderEntityPath, err)
				continue
			}
			changeEntityPermLog(
				outPutFolderEntityPath,
				lastEntityInfo.Mode(),
				currentEntityInfo.Mode(),
			)
		default:
			unsupportedFsEventLog(fsEventName)
			continue
		}
		syncSuccessfulLog(outPutFolderEntityPath, fsEventName)
	}
	return nil
}

func (f *FsSyncer) Stop() {
	f.dirWatcher.Stop()
}

func newEventLog(fsEvent dir_watcher.FSNotifyEvent) {
	slog.Info(
		"detected new fs event",
		slog.String(
			"event name",
			string(fsEvent.EventName),
		),
		slog.String(
			"event time",
			utils.FormatTimeISO8601(fsEvent.EventTime),
		),
		slog.String(
			"os path",
			fsEvent.OsPath,
		),
	)
}

func successDeleteLog(path string) {
	slog.Info(
		"remove successful",
		slog.String(
			"path",
			path,
		),
	)
}

func errorDeleteLog(path string, err error) {
	slog.Error(
		"failed delete",
		slog.String(
			"path",
			path,
		),
		slog.String(
			"error",
			err.Error(),
		),
	)
}

func successCreateLog(path string) {
	slog.Info(
		"create successful",
		slog.String(
			"path",
			path,
		),
	)
}

func errorCreateLog(path string, err error) {
	slog.Error(
		"failed create",
		slog.String(
			"path",
			path,
		),
		slog.String(
			"error",
			err.Error(),
		),
	)
}

func changeEntityPermLog(path string, lastPerm, currentPerm fs.FileMode) {
	slog.Info(
		"change perm entity successful",
		slog.String(
			"path",
			path,
		),
		slog.String(
			"last entity perm",
			lastPerm.String(),
		),
		slog.String(
			"current entity perm",
			currentPerm.String(),
		),
	)
}

func failedChangeEntityPermLog(path string, err error) {
	slog.Info(
		"failed change perm entity",
		slog.String(
			"path",
			path,
		),
		slog.String(
			"error",
			err.Error(),
		),
	)
}

func failedGetOsEntityStatLog(path string, err error) {
	slog.Info(
		"failed get os entity info",
		slog.String(
			"path",
			path,
		),
		slog.String(
			"error",
			err.Error(),
		),
	)
}

func syncSuccessfulLog(path string, eventName constants.FsEventName) {
	slog.Info(
		"sync successful",
		slog.String(
			"path",
			path,
		),
		slog.String(
			"event name",
			string(eventName),
		),
	)
}

func unsupportedFsEventLog(eventName constants.FsEventName) {
	slog.Warn(
		"unsupported fs event name",
		slog.String(
			"event name",
			string(eventName),
		),
	)
}

func successfulSyncFileContentLog(srcFilePath, distFilePath string) {
	slog.Info(
		"successful sync files",
		slog.String(
			"source file path",
			srcFilePath,
		),
		slog.String(
			"destination file path",
			distFilePath,
		),
	)
}

func failedSyncFileContentLog(srcFilePath, distFilePath string) {
	slog.Info(
		"failed sync files",
		slog.String(
			"source file path",
			srcFilePath,
		),
		slog.String(
			"destination file path",
			distFilePath,
		),
	)
}
