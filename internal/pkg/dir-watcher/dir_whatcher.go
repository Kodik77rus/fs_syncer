package dir_watcher

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"time"
	"unsafe"

	"github.com/Kodik77rus/fs_syncer/internal/pkg/config"
	"github.com/Kodik77rus/fs_syncer/internal/pkg/constants"
	"github.com/Kodik77rus/fs_syncer/internal/pkg/utils"

	"golang.org/x/sys/unix"
)

var ErrReadTimeout = errors.New("read timeout")

type DirWatcher struct {
	cancelFunc    func()
	srcFolderPath string
}

type FSNotifyEvent struct {
	EventTime time.Time
	OsPath    string
	EventName constants.FsEventName
	Error     error
}

func InitDirWatcher(config *config.Config) *DirWatcher {
	return &DirWatcher{
		srcFolderPath: config.SrcFolderPath,
	}
}

func (d *DirWatcher) Start() <-chan FSNotifyEvent {
	fsNotifyChan := make(chan FSNotifyEvent, 1000)
	go func() {
		defer close(fsNotifyChan)
		fd, err := unix.InotifyInit1(0)
		if err != nil {
			fsNotifyChan <- FSNotifyEvent{
				Error: err,
			}
			return
		}
		defer unix.Close(fd)
		ctx, cancel := context.WithCancel(context.Background())
		d.cancelFunc = cancelFunc(cancel)
		defer d.cancelFunc()
		if _, err = unix.InotifyAddWatch(
			fd,
			d.srcFolderPath,
			unix.IN_CREATE|
				unix.IN_DELETE|
				unix.IN_DELETE_SELF|
				unix.IN_IGNORED|
				unix.IN_EXCL_UNLINK|
				unix.IN_CLOSE_WRITE|
				unix.IN_ATTRIB,
		); err != nil {
			fsNotifyChan <- FSNotifyEvent{
				Error: err,
			}
			return
		}
		errorChan := initFdWatcher(ctx, fd, fsNotifyChan)
	Loop:
		for {
			select {
			case <-ctx.Done():
				break Loop
			case err := <-errorChan:
				fsNotifyChan <- FSNotifyEvent{
					Error: err,
				}
				return
			}
		}
	}()
	return fsNotifyChan
}

func (d *DirWatcher) Stop() {
	d.cancelFunc()
}

func cancelFunc(ctxCancelFunc context.CancelFunc) func() {
	var once sync.Once
	return func() {
		once.Do(func() { ctxCancelFunc() })
	}
}

func initFdWatcher(ctx context.Context, fd int, fsEventChan chan<- FSNotifyEvent) <-chan error {
	errorChan := make(chan error)
	go func() {
		defer close(errorChan)
		var buff [(unix.SizeofInotifyEvent + unix.NAME_MAX + 1) * 20]byte
		for {
			offset := 0
			n, err := unix.Read(fd, buff[:])
			if err != nil {
				errorChan <- err
				return
			}
			if n == 0 {
				continue
			}
			select {
			case <-ctx.Done():
				return
			default:
				for offset < n {
					var fSNotifyEvent FSNotifyEvent
					fSNotifyEvent.EventTime = utils.GetCurrentTime()
					e := (*unix.InotifyEvent)(unsafe.Pointer(&buff[offset]))
					nameBs := buff[offset+unix.SizeofInotifyEvent : offset+unix.SizeofInotifyEvent+int(e.Len)]
					name := string(bytes.TrimRight(nameBs, "\x00"))
					fSNotifyEvent.OsPath = name
					if len(name) > 0 && e.Mask&unix.IN_ISDIR == unix.IN_ISDIR {
						// skip sub dir events
						continue
					}
					switch {
					case e.Mask&unix.IN_CREATE == unix.IN_CREATE:
						fSNotifyEvent.EventName = constants.Create
					case e.Mask&unix.IN_DELETE == unix.IN_DELETE ||
						e.Mask&unix.IN_DELETE_SELF == unix.IN_DELETE_SELF ||
						e.Mask&unix.IN_IGNORED == unix.IN_IGNORED:
						fSNotifyEvent.EventName = constants.Delete
					case e.Mask&unix.IN_CLOSE_WRITE == unix.IN_CLOSE_WRITE:
						fSNotifyEvent.EventName = constants.Write
					case e.Mask&unix.IN_ATTRIB == unix.IN_ATTRIB:
						fSNotifyEvent.EventName = constants.ChangeAttrib
					default:
						fSNotifyEvent.EventName = constants.Unknown
					}
					offset += int(unix.SizeofInotifyEvent + e.Len)
					fsEventChan <- fSNotifyEvent
				}
			}
		}
	}()
	return errorChan
}
