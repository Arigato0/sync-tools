package clipdb

const (
	TYPE_TEXT = iota
	TYPE_IMG_FILE
	TYPE_BIN_FILE
	TYPE_FILE
	TYPE_DIR
)

type Entry struct {
	Type     int
	Filename string
}
