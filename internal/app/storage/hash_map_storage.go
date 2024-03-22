package storage

import (
	"context"

	"github.com/ruslanjo/url_shortener/internal/app/storage/disk"
	"github.com/ruslanjo/url_shortener/internal/app/storage/models"
	"github.com/ruslanjo/url_shortener/internal/config"
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
	url := disk.URLModel{ShortLink: shortLink, FullLink: fullLink}

	if s.diskStorage == nil { // Мб здесь как-то по другому лучше nil-интерфейс обрабатывать?
		return nil
	}
	if err := s.diskStorage.Persist(url); err != nil {
		return err
	}

	return nil
}

func (s *HashMapStorage) SaveURLBatched(ctx context.Context, data []models.URLBatch) error {
	// context not used in hashmap storage

	for curPtr := 0; curPtr < len(data); curPtr += config.URLBatchSize {
		upperBound := curPtr+config.URLBatchSize
		if upperBound > len(data){
			upperBound = len(data)
		}
		batch := data[curPtr : upperBound]

		if err := s.diskStorage.PersistBatch(batch); err != nil {
			return err
		}

		for i := 0; i < len(batch); i++ {
			short, long := batch[i].ShortURL, batch[i].OriginalURL
			s.data[short] = long
		}
	}
	return nil
}


func (s *HashMapStorage) PingContext (ctx context.Context) error {
	return nil
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
