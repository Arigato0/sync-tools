package clipdb

import (
	"encoding/json"
	"errors"
	"io/fs"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	TYPE_TEXT = iota
	TYPE_FILE
	TYPE_DIR
)

func TypeString(t int) string {
	switch t {
	case TYPE_TEXT:
		return "text"
	case TYPE_FILE:
		return "file"
	case TYPE_DIR:
		return "directory"
	default:
		return "unknown"
	}
}

type Entry struct {
	Type int
	// if the type is of a file or dir it will be stored here otherwise empty
	Filename string
	// date of when the entry was added mostly used for sorting
	Date time.Time
	// metadata specific to the type
	MetaData map[string]any
}

func NewTextEntry() Entry {
	return Entry{
		Type: TYPE_TEXT,
		Date: time.Now(),
	}
}

func NewFsEntry(path string, entryType int) Entry {
	return Entry{
		Filename: filepath.Base(path),
		Date:     time.Now(),
		Type:     entryType,
	}
}

func NewDirEntry(path string) Entry {
	return NewFsEntry(path, TYPE_DIR)
}

func NewFileEntry(path string) Entry {
	return NewFsEntry(path, TYPE_FILE)
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func makeId(fmt string) string {

	builder := strings.Builder{}

	for _, c := range fmt {
		if c == 'x' {
			value := byte(randRange(48, 57))
			builder.WriteByte(value)
		} else {
			builder.WriteRune(c)
		}
	}

	return builder.String()
}

func validateCacheDir(filename string) string {
	cacheDir, err := os.UserCacheDir()

	if err != nil {
		return ""
	}

	cacheDir = filepath.Join(cacheDir, "nclipdb")

	os.Mkdir(cacheDir, os.ModePerm)

	return filepath.Join(cacheDir, filename)
}

func copyDir(src, dst string) error {

	srcBase := filepath.Base(src)
	dstDir := filepath.Dir(dst)

	idx := strings.LastIndex(src, srcBase)

	err := filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {

		newPath := filepath.Join(dstDir, path[idx:])

		if info.IsDir() {
			os.Mkdir(newPath, os.ModePerm)
			return nil
		}

		data, err := os.ReadFile(path)

		if err != nil {
			return err
		}

		return os.WriteFile(newPath, data, os.ModePerm)
	})

	if err != nil {
		return err
	}

	return os.Rename(filepath.Join(dstDir, srcBase), dst)
}

func (entry *Entry) Save(data []byte) error {

	filename := TypeString(entry.Type) + makeId("-xxxx-xxxx")

	if entry.Type != TYPE_TEXT && entry.Filename == "" {
		return errors.New("for none text based entries a name is required to restore the original file")
	}

	fullPath := validateCacheDir(filename)

	if entry.Type == TYPE_DIR {
		path := string(data)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			return errors.New("the given directory path does not exist")
		}

		err := copyDir(path, fullPath)

		if err != nil {
			return err
		}

	} else if entry.Type == TYPE_FILE {
		contents, err := os.ReadFile(string(data))

		if err != nil {
			return err
		}

		os.WriteFile(fullPath, contents, os.ModePerm)

	} else {
		os.WriteFile(fullPath, data, os.ModePerm)
	}

	jsonData, err := json.Marshal(entry)

	if err != nil {
		return err
	}

	return os.WriteFile(fullPath+".clip_entry", jsonData, os.ModePerm)
}
