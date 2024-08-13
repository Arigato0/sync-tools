package clipdb

import "time"

const (
	TYPE_TEXT = iota
	TYPE_IMG_FILE
	TYPE_BIN_FILE
	TYPE_FILE
	TYPE_DIR
)

type Entry struct {
	Type int
	// the filename of where the actual data is stored
	Filename string
	// if the type is of a file or dir it will be stored here otherwise empty
	Name string
	// date of when the entry was added mostly used for sorting
	Date time.Time
	// metadata specific to the type
	MetaData map[string]any
}
