package storage

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const (
	DataFileName = "lumora.data"
	headerSize   = 12
	MaxValueSize = 40 * 1024 * 1024
)

type DataRecord struct {
	Timestamp uint32
	KeySize   uint32
	ValueSize uint32
	Key       []byte
	Value     []byte
}

type DataManager struct {
	File     *os.File
	filepath string
}

func NewDataManager(path string) (*DataManager, error) {
	filepath := fmt.Sprintf("%s/%s", path, DataFileName)

	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to load data file: %w", err)
	}

	return &DataManager{
		File:     file,
		filepath: filepath,
	}, nil
}

func (dm *DataManager) WriteRecord(record *DataRecord) (int64, error) {

	if record == nil {
		return -1, ErrInvalidArgument
	}

	offset, err := dm.File.Seek(0, io.SeekEnd)
	if err != nil {
		return -1, fmt.Errorf("seek failed: %w", err)
	}

	header := make([]byte, 12)

	binary.BigEndian.PutUint32(header[0:4], record.Timestamp)
	binary.BigEndian.PutUint32(header[4:8], record.KeySize)
	binary.BigEndian.PutUint32(header[8:12], record.ValueSize)

	if _, err := dm.File.Write(header); err != nil {
		return -1, fmt.Errorf("header write failed: %w", err)
	}
	if _, err := dm.File.Write(record.Key); err != nil {
		return -1, fmt.Errorf("key write failed: %w", err)
	}
	if _, err := dm.File.Write(record.Value); err != nil {
		return -1, fmt.Errorf("value write failed: %w", err)
	}

	if err := dm.File.Sync(); err != nil {
		return -1, fmt.Errorf("sync failed: %w", err)
	}

	return offset, nil
}

func (dm *DataManager) ReadRecord(offset int64) (*DataRecord, error) {

	if _, err := dm.File.Seek(offset, io.SeekStart); err != nil {
		return nil, fmt.Errorf("seek to offset %d failed: %w", offset, err)
	}

	header := make([]byte, 12)
	if _, err := io.ReadFull(dm.File, header); err != nil {
		return nil, fmt.Errorf("header read failed: %w", err)
	}

	record := &DataRecord{
		Timestamp: binary.BigEndian.Uint32(header[0:4]),
		KeySize:   binary.BigEndian.Uint32(header[4:8]),
		ValueSize: binary.BigEndian.Uint32(header[8:12]),
	}

	// 16mb max key size
	if record.KeySize > 1<<24 {
		return nil, fmt.Errorf("%w: invalid key size %d", ErrDataCorruption, record.KeySize)
	}

	record.Key = make([]byte, record.KeySize)
	if _, err := io.ReadFull(dm.File, record.Key); err != nil {
		return nil, fmt.Errorf("key read failed: %w", err)
	}

	if record.ValueSize > MaxValueSize {
		return nil, fmt.Errorf("%w: value size exceeds maximum allowed size of %d bytes", ErrDataCorruption, MaxValueSize)
	}

	record.Value = make([]byte, record.ValueSize)
	if _, err := io.ReadFull(dm.File, record.Value); err != nil {
		return nil, fmt.Errorf("value read failed: %w", err)
	}

	return record, nil
}

func (dm *DataManager) Close() error {

	if dm.File != nil {
		err := dm.File.Close()
		dm.File = nil
		return err
	}
	return nil
}
