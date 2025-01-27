package storage

import (
	"encoding/gob"
	"io"
	"os"
)

const (
	IndexFileName = "lumora.index"
)

type IndexEntry struct {
	Offset    int64
	KeySize   uint32
	ValueSize uint32
	Timestamp uint32
}

type IndexManager struct {
	Index    map[string]IndexEntry
	file     *os.File
	filepath string
}

func NewIndexManager(path string) (*IndexManager, error) {

	filepath := path + "/" + IndexFileName

	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	index := make(map[string]IndexEntry)

	if info, _ := file.Stat(); info.Size() > 0 {
		decoder := gob.NewDecoder(file)
		if err := decoder.Decode(&index); err != nil {
			return nil, err
		}
	}

	return &IndexManager{
		Index:    index,
		file:     file,
		filepath: filepath,
	}, nil
}

func (im *IndexManager) Save() error {
	im.file.Truncate(0)
	im.file.Seek(0, io.SeekCurrent)
	encoder := gob.NewEncoder(im.file)
	return encoder.Encode(im.Index)
}
