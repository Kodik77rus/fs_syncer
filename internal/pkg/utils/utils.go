package utils

import (
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/Kodik77rus/fs_syncer/internal/pkg/constants"
)

func CheckOrCreateFile(filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			if err := CreateFile(filePath); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func CheckOrCreateFolders(foldersPaths ...string) error {
	for _, path := range foldersPaths {
		if err := CheckOrCreateFolder(path); err != nil {
			return err
		}
	}
	return nil
}

func CheckOrCreateFolder(folderPath string) error {
	_, err := GetOsStat(folderPath)
	if os.IsNotExist(err) {
		return CreateFolder(folderPath)
	}
	if err != nil {
		return err
	}
	return nil
}

func GetOsStat(path string) (fs.FileInfo, error) {
	return os.Stat(path)
}

func GetCurrentTime() time.Time {
	return time.Now()
}

func OsRemove(path string) error {
	return os.Remove(path)
}

func OsCreate(path string) (*os.File, error) {
	return os.Create(path)
}

func CreateFolder(folderPath string) error {
	return os.MkdirAll(folderPath, os.ModePerm)
}

func ChangeOsPerm(path string, mode fs.FileMode) error {
	return os.Chmod(path, mode)
}

func CreateFile(path string) error {
	f, err := OsCreate(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error(
				"failed close file",
				slog.String(
					"error",
					err.Error(),
				),
			)
		}
	}()
	return nil
}

func JoinOsPath(paths ...string) string {
	return filepath.Join(paths...)
}

func FormatTimeISO8601(time time.Time) string {
	return time.Format(constants.ISO8601TimeFormat)
}

func CopyFileContent(srcFilePath, destFilePath string) error {
	dstFile, err := OsCreate(destFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := dstFile.Close(); err != nil {
			slog.Error(
				"failed close file",
				slog.String(
					"error",
					err.Error(),
				),
			)
		}
	}()
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := srcFile.Close(); err != nil {
			slog.Error(
				"failed close file",
				slog.String(
					"error",
					err.Error(),
				),
			)
		}
	}()
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	return dstFile.Sync()
}

func WriteStringIntoFile(path string, data *string) error {
	file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(*data); err != nil {
		return err
	}
	return nil
}
