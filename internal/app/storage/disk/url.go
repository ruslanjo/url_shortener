package disk

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ruslanjo/url_shortener/internal/app/storage/models"
)

type JSONURLDiskStorage struct {
	diskStorage diskStorage
}

func NewURLDiskStorage(path string) *JSONURLDiskStorage {
	ds := diskStorage{path}
	return &JSONURLDiskStorage{diskStorage: ds}
}

func (m JSONURLDiskStorage) Persist(entity URLModel) error {
	data, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	if err = m.diskStorage.persist(data); err != nil {
		return err
	}
	return nil
}

func (m JSONURLDiskStorage) PersistBatch(data []models.URLBatch) error {
	bytesData, err := m.prepareBatch(data)
	if err != nil {
		return err
	}

	if err = m.diskStorage.persist(bytesData); err != nil {
		return err
	}
	return nil
}

func (m JSONURLDiskStorage) prepareBatch(data []models.URLBatch) ([]byte, error) {
	var out []byte

	for i, ent := range data {
		m := URLModel{
			FullLink:  ent.OriginalURL,
			ShortLink: ent.ShortURL,
		}
		marsh, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}

		if i != len(data)-1 {
			marsh = append(marsh, '\n')
		}

		out = append(out, marsh...)
	}
	return out, nil
}

func (m JSONURLDiskStorage) ReadAll() ([]URLModel, error) {
	data, err := m.diskStorage.readAll()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []URLModel{}, nil
	}
	dataStr := strings.Trim(string(data), "\n")
	chunks := strings.Split(dataStr, "\n")
	out := make([]URLModel, 0, len(chunks))
	for _, partition := range chunks {
		schema := new(URLModel)
		err = json.Unmarshal([]byte(partition), schema)
		if err != nil {
			return nil, fmt.Errorf("unable to load URL data from disk: %s", err.Error())
		}
		out = append(out, *schema)
	}
	return out, nil
}
