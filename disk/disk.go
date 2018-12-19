package disk

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	redislight "github.com/ITandElectronics/GoHomework"
)

// New construct new storage
func New(path string) (*Disk, error) {
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.FileMode(0640))
	if err != nil {
		return nil, err
	}
	data := make(map[string]string)
	if err := json.NewDecoder(fd).Decode(&data); err != nil && err != io.EOF {
		return nil, err
	}
	return &Disk{
		m:    &sync.RWMutex{},
		fd:   fd,
		data: data,
	}, nil
}

// Disk represent disk storage
type Disk struct {
	m    *sync.RWMutex
	fd   *os.File
	data map[string]string
}

// Get checks whether key exists and either return appropriate value or KeyIsNotExists error
func (d *Disk) Get(key string) (string, error) {
	d.m.RLock()
	defer d.m.RUnlock()
	val, ok := d.data[key]
	if !ok {
		return "", redislight.ErrKeyIsNotExists
	}
	return val, nil
}

// Set create new or update existing key value
func (d *Disk) Set(key, value string) error {
	d.m.Lock()
	defer d.m.Unlock()
	d.data[key] = value
	return d.sync()
}

// Del remove entry related to provided key from the storage
func (d *Disk) Del(key string) error {
	d.m.Lock()
	defer d.m.Unlock()
	if _, ok := d.data[key]; !ok {
		return redislight.ErrKeyIsNotExists
	}
	delete(d.data, key)
	return d.sync()
}

// all sync calls should be protected by mutex
func (d *Disk) sync() error {
	_, err := d.fd.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("coudn't seek to the start of the file: %v", err)
	}
	return json.NewEncoder(d.fd).Encode(d.data)
}
