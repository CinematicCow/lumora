package storage

import (
	"encoding/binary"
	"io"
	"os"
)

const (
	DataFileName = "lumora.data"
)

type DataRecord struct {
	Timestamp uint32
	KeySize   uint32
	ValueSize uint32
	Key       []byte
	Value     []byte
}

type DataManager struct {
	file     *os.File
	filepath string
}

func NewDataManager(path string) (*DataManager, error) {
	filepath := path + "/" + DataFileName

	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return &DataManager{
		file:     file,
		filepath: filepath,
	}, nil
}

func (dm *DataManager) WriteRecord(record *DataRecord) (int64, error) {
	offset, _ := dm.file.Seek(0, io.SeekCurrent)

	header := make([]byte, 12)

	binary.BigEndian.PutUint32(header[0:4], record.Timestamp)
	binary.BigEndian.PutUint32(header[4:8], record.KeySize)
	binary.BigEndian.PutUint32(header[8:12], record.ValueSize)

	if _, err := dm.file.Write(header); err != nil {
		return -1, err
	}
	if _, err := dm.file.Write(record.Key); err != nil {
		return -1, err
	}
	if _, err := dm.file.Write(record.Value); err != nil {
		return -1, err
	}

	return offset, nil
}

func (dm *DataManager) ReadRecord(offset int64) (*DataRecord, error) {

	_, err := dm.file.Seek(offset, io.SeekStart)

	if err != nil {
		return nil, err
	}

	header := make([]byte, 12)

	if _, err := dm.file.Read(header); err != nil {
		return nil, err
	}

	record := &DataRecord{
		Timestamp: binary.BigEndian.Uint32(header[0:4]),
		KeySize:   binary.BigEndian.Uint32(header[4:8]),
		ValueSize: binary.BigEndian.Uint32(header[8:12]),
	}

	record.Key = make([]byte, record.KeySize)
	if _, err := dm.file.Read(record.Key); err != nil {
		return nil, err
	}

	record.Value = make([]byte, record.ValueSize)
	if _, err := dm.file.Read(record.Value); err != nil {
		return nil, err
	}

	return record, nil
}
