package disk

import (
	"io/ioutil"
	"os"
	"sync"
)

type URLSchema struct {
	FullLink  string `json:"full_link"`
	ShortLink string `json:"short_link"`
}

type URLStorage interface {
	Persist(URLSchema) error
	ReadAll() ([]URLSchema, error)
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
	data, err = ioutil.ReadAll(f)
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
	p_data := make([]byte, len(data), len(data)+1)
	_ = copy(p_data, data)
	p_data = append(p_data, '\n')

	if _, err = f.Write(p_data); err != nil {
		return err
	}
	return nil

}
