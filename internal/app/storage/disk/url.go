package disk

import (
	"encoding/json"
	"fmt"
	"strings"
)

type JSONURLDiskStorage struct {
	diskStorage DiskStorage
}

func NewURLDiskStorage(ds DiskStorage) *JSONURLDiskStorage {
	return &JSONURLDiskStorage{diskStorage: ds}
}

func (m JSONURLDiskStorage) Persist(entity URLSchema) error {
	data, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	if err = m.diskStorage.persist(data); err != nil {
		return err
	}
	return nil
}

func (m JSONURLDiskStorage) ReadAll() ([]URLSchema, error) {
	data, err := m.diskStorage.readAll()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []URLSchema{}, nil
	}
	dataStr := strings.Trim(string(data), "\n")
	chunks := strings.Split(dataStr, "\n")
	out := make([]URLSchema, 0, len(chunks))
	for _, partition := range chunks {
		schema := new(URLSchema)
		err = json.Unmarshal([]byte(partition), schema)
		if err != nil {
			return nil, fmt.Errorf("unable to load URL data from disk: %s", err.Error())
		}
		out = append(out, *schema)
	}
	return out, nil
}
