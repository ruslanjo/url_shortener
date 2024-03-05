package disk

import (
	"io"
	"os"
	"sync"
)

type URLModel struct {
	FullLink  string `json:"full_link"`
	ShortLink string `json:"short_link"`
}

type URLStorage interface {
	Persist(URLModel) error
	ReadAll() ([]URLModel, error)
}

type DiskStorage struct {
	Path string
}

func (d DiskStorage) readAll() ([]byte, error) {
	var data []byte

	m := sync.Mutex{}
	m.Lock()
	defer m.Unlock()

	f, err := os.OpenFile(d.Path, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err = io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (d DiskStorage) persist(data []byte) error {
	m := sync.Mutex{}
	m.Lock()
	defer m.Unlock()
	f, err := os.OpenFile(d.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	pData := make([]byte, len(data), len(data)+1)
	_ = copy(pData, data)
	pData = append(pData, '\n')

	if _, err = f.Write(pData); err != nil {
		return err
	}
	return nil

}
