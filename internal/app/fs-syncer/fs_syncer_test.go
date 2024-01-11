package fs_syncer

import (
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Kodik77rus/fs_syncer/internal/pkg/config"
	dir_watcher "github.com/Kodik77rus/fs_syncer/internal/pkg/dir-watcher"
	"github.com/Kodik77rus/fs_syncer/internal/pkg/utils"
)

func TestFsSyncer(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	c := &config.Config{}
	srcTmpDirPath := t.TempDir()
	outputTmpDirPath := t.TempDir()
	defer t.Cleanup(func() {
		os.RemoveAll(srcTmpDirPath)
		os.RemoveAll(outputTmpDirPath)
	})
	testFileName := "test1.txt"
	sourceFilePath := utils.JoinOsPath(srcTmpDirPath, testFileName)
	outPutFilePath := utils.JoinOsPath(outputTmpDirPath, testFileName)
	c.SrcFolderPath = srcTmpDirPath
	c.OutputFolderPath = outputTmpDirPath
	fsSyncer := InitFsWatcher(c, dir_watcher.InitDirWatcher(c))
	go func() {
		req.NoError(fsSyncer.Start())
	}()
	defer fsSyncer.Stop()
	time.Sleep(1 * time.Second)
	t.Run("crete file successful", func(t *testing.T) {
		file, err := utils.OsCreate(sourceFilePath)
		req.NoError(err)
		defer file.Close()
		originalFileInfo, err := file.Stat()
		req.NoError(err)
		time.Sleep(100 * time.Millisecond)
		outputFileInfo, err := os.Stat(outPutFilePath)
		req.NoError(err)
		req.Equal(originalFileInfo.Name(), outputFileInfo.Name())
		req.Equal(originalFileInfo.Size(), outputFileInfo.Size())
		req.Equal(originalFileInfo.Mode(), outputFileInfo.Mode())
	})
	t.Run("write file successful", func(t *testing.T) {
		data := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed nisi lorem, elementum at odio et, aliquet consectetur enim. Ut laoreet mattis nunc id convallis. Proin id erat tempor, faucibus enim sed, vehicula nulla. Mauris scelerisque dui in turpis ultricies, ut maximus ligula pellentesque. Sed turpis nisl, finibus sit amet scelerisque sit amet, suscipit eu ante. Nunc varius malesuada risus consequat fringilla. Etiam neque mi, congue ut mauris ac, gravida varius nulla. Etiam bibendum massa est, ac facilisis tortor dictum nec. Sed ut rhoncus enim. Sed ac malesuada augue. Quisque nec tempor turpis. In hac habitasse platea dictumst. Aliquam elementum lacus tincidunt vehicula pretium."
		err := utils.WriteStringIntoFile(sourceFilePath, &data)
		req.NoError(err)
		originalFileData, err := os.ReadFile(sourceFilePath)
		req.NoError(err)
		time.Sleep(100 * time.Millisecond)
		outputFileData, err := os.ReadFile(outPutFilePath)
		req.NoError(err)
		req.Equal(originalFileData, outputFileData)
	})
	t.Run("change file attribute successful", func(t *testing.T) {
		err := utils.ChangeOsPerm(sourceFilePath, os.ModePerm)
		req.NoError(err)
		originalFileStat, err := utils.GetOsStat(sourceFilePath)
		req.NoError(err)
		time.Sleep(100 * time.Millisecond)
		outPutFileStat, err := utils.GetOsStat(outPutFilePath)
		req.NoError(err)
		req.Equal(originalFileStat.Mode().Perm(), outPutFileStat.Mode().Perm())
	})
	t.Run("remove file successful", func(t *testing.T) {
		req.NoError(utils.OsRemove(sourceFilePath))
		time.Sleep(100 * time.Millisecond)
		_, err := os.Stat(outPutFilePath)
		req.ErrorIs(err, fs.ErrNotExist)
	})

}
