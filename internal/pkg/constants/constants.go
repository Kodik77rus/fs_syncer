package constants

type FsEventName string

const (
	Create       FsEventName = "CREATE"
	Delete       FsEventName = "DELETE"
	Write        FsEventName = "WRITE"
	Unknown      FsEventName = "UNKNOWN"
	ChangeAttrib FsEventName = "CHANGE_ATTRIB"
)

const (
	ISO8601TimeFormat = "2006-01-02 15:04:05"
)
