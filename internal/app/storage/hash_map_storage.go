package storage

import (
	"github.com/ruslanjo/url_shortener/internal/app/storage/disk"
)

type HashMapStorage struct {
	data        map[string]string
	diskStorage disk.URLStorage
}

func NewHashMapStorage(disk disk.URLStorage) *HashMapStorage {
	d := map[string]string{}
	return &HashMapStorage{data: d, diskStorage: disk}

}

func (s *HashMapStorage) GetURLByShortLink(shortLink string) (string, error) {

	if s.data == nil {
		s.data = make(map[string]string)
	}

	if value, ok := s.data[shortLink]; ok {
		return value, nil
	} else {
		return value, ErrURLMappingNotFound
	}
}

func (s *HashMapStorage) AddShortURL(shortLink string, fullLink string) error {
	if s.data == nil {
		s.data = make(map[string]string)
	}
	s.data[shortLink] = fullLink
	url := disk.URLSchema{ShortLink: shortLink, FullLink: fullLink}

	if s.diskStorage == nil{ // Мб здесь как-то по другому лучше nil-интерфейс обрабатывать?
		return nil
	}
	if err := s.diskStorage.Persist(url); err != nil {
		return err
	}

	return nil
}

func (s *HashMapStorage) InitStorage(data map[string]string) {
	// Is needed for unit-tests
	s.data = data
}

func (s *HashMapStorage) LoadFromDisk() error {
	d := make(map[string]string)
	diskData, err := s.diskStorage.ReadAll()
	if err != nil {
		return err
	}
	for _, ent := range diskData {
		d[ent.ShortLink] = ent.FullLink
	}
	s.data = d
	return nil
}
