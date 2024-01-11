package dir_watcher

import (
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Kodik77rus/fs_syncer/internal/pkg/config"
	"github.com/Kodik77rus/fs_syncer/internal/pkg/constants"
	"github.com/Kodik77rus/fs_syncer/internal/pkg/utils"
)

func TestDirWatcher(t *testing.T) {
	t.Parallel()
	req := require.New(t)
	c := &config.Config{}
	tmpDirPath := t.TempDir()
	defer t.Cleanup(func() {
		os.Remove(tmpDirPath)
	})
	c.SrcFolderPath = tmpDirPath
	testFileName := "test1.txt"
	testFilePath := utils.JoinOsPath(tmpDirPath, testFileName)
	dirWatcher := InitDirWatcher(c)
	fSNotifyEventChan := dirWatcher.Start()
	defer dirWatcher.Stop()
	dataToWrite := "test"
	cases := []struct {
		testName string
		helpFunc func() error
		want     FSNotifyEvent
	}{
		{
			testName: "expect create event",
			helpFunc: createFile(testFilePath),
			want: FSNotifyEvent{
				OsPath:    testFileName,
				EventName: constants.Create,
			},
		},
		{
			testName: "expect Write event",
			helpFunc: writeStringIntoFile(testFilePath, &dataToWrite),
			want: FSNotifyEvent{
				OsPath:    testFileName,
				EventName: constants.Write,
			},
		},
		{
			testName: "expect ChangeAttrib event",
			helpFunc: changeFileAttr(testFilePath, 0777),
			want: FSNotifyEvent{
				OsPath:    testFileName,
				EventName: constants.ChangeAttrib,
			},
		},
		{
			testName: "expect Delete event",
			helpFunc: deleteFile(testFilePath),
			want: FSNotifyEvent{
				OsPath:    testFileName,
				EventName: constants.Delete,
			},
		},
	}
	time.Sleep(100 * time.Millisecond)
	for _, testCase := range cases {
		t.Run(testCase.testName, func(t *testing.T) {
			err := testCase.helpFunc()
			req.NoError(err)
			fsEvent := <-fSNotifyEventChan
			req.NoError(fsEvent.Error)
			req.Equal(testCase.want.EventName, fsEvent.EventName)
			req.Equal(testCase.want.OsPath, fsEvent.OsPath)
		})
	}
}

func createFile(path string) func() error {
	return func() error {
		_, err := utils.OsCreate(path)
		return err
	}
}

func writeStringIntoFile(path string, data *string) func() error {
	return func() error {
		return utils.WriteStringIntoFile(path, data)
	}
}

func changeFileAttr(path string, perm fs.FileMode) func() error {
	return func() error {
		return utils.ChangeOsPerm(path, perm)
	}
}

func deleteFile(path string) func() error {
	return func() error {
		return utils.OsRemove(path)
	}
}
