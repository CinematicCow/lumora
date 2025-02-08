package storage

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"
)

const (
	DataFileName = "lumora.data"
	headerSize   = 12
	MaxValueSize = 40 * 1024 * 1024
	MaxKeySize   = 40 * 1024 * 1024
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
	mu       sync.Mutex
	closed   bool
}

func NewDataManager(path string) (*DataManager, error) {
	filepath := fmt.Sprintf("%s/%s", path, DataFileName)

	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to load data file: %w", err)
	}

	return &DataManager{
		file:     file,
		filepath: filepath,
	}, nil
}

func (dm *DataManager) WriteRecord(record *DataRecord) (int64, error) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.closed {
		return -1, ErrDataClosed
	}

	if record == nil {
		return -1, ErrInvalidArgument
	}

	offset, err := dm.file.Seek(0, io.SeekEnd)
	if err != nil {
		return -1, fmt.Errorf("seek failed: %w", err)
	}

	header := make([]byte, 12)

	binary.BigEndian.PutUint32(header[0:4], record.Timestamp)
	binary.BigEndian.PutUint32(header[4:8], record.KeySize)
	binary.BigEndian.PutUint32(header[8:12], record.ValueSize)

	if _, err := dm.file.Write(header); err != nil {
		return -1, fmt.Errorf("header write failed: %w", err)
	}
	if _, err := dm.file.Write(record.Key); err != nil {
		return -1, fmt.Errorf("key write failed: %w", err)
	}
	if _, err := dm.file.Write(record.Value); err != nil {
		return -1, fmt.Errorf("value write failed: %w", err)
	}

	if err := dm.file.Sync(); err != nil {
		return -1, fmt.Errorf("sync failed: %w", err)
	}

	return offset, nil
}

func (dm *DataManager) ReadRecord(offset int64) (*DataRecord, error) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.closed {
		return nil, ErrDataClosed
	}

	if _, err := dm.file.Seek(offset, io.SeekStart); err != nil {
		return nil, fmt.Errorf("%w: seek failed: %v", ErrDataCorruption, err)
	}

	header := make([]byte, headerSize)
	n, err := io.ReadFull(dm.file, header)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil, fmt.Errorf("%w: incomplete header (%d bytes)",
				ErrDataCorruption, n)
		}
		return nil, fmt.Errorf("header read failed: %w", err)
	}

	record := &DataRecord{
		Timestamp: binary.BigEndian.Uint32(header[0:4]),
		KeySize:   binary.BigEndian.Uint32(header[4:8]),
		ValueSize: binary.BigEndian.Uint32(header[8:12]),
	}

	if record.KeySize > MaxKeySize {
		return nil, fmt.Errorf("%w: key size %d exceeds maximum",
			ErrDataCorruption, record.KeySize)
	}

	record.Key = make([]byte, record.KeySize)
	if _, err := io.ReadFull(dm.file, record.Key); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil, fmt.Errorf("%w: incomplete key data", ErrDataCorruption)
		}
		return nil, fmt.Errorf("key read failed: %w", err)
	}

	if record.ValueSize > MaxValueSize {
		return nil, fmt.Errorf("%w: value size %d exceeds maximum",
			ErrDataCorruption, record.ValueSize)
	}

	record.Value = make([]byte, record.ValueSize)
	if _, err := io.ReadFull(dm.file, record.Value); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil, fmt.Errorf("%w: incomplete value data", ErrDataCorruption)
		}
		return nil, fmt.Errorf("value read failed: %w", err)
	}

	return record, nil
}

func (dm *DataManager) Close() error {

	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.closed {
		return nil
	}

	if err := dm.file.Close(); err != nil {
		return fmt.Errorf("data file close failed: %w", err)
	}

	dm.closed = true
	return nil
}
