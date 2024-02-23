package storage


type HashMapStorage struct {
	storage map[string]string
}

func (dao *HashMapStorage) GetURLByShortLink(shortLink string) (string, error) {

	if dao.storage == nil {
		dao.storage = make(map[string]string)
	}

	if value, ok := dao.storage[shortLink]; ok {
		return value, nil
	} else {
		return value, ErrURLMappingNotFound
	}
}

func (dao *HashMapStorage) AddShortURL(shortLink string, fullLink string) error {
	if dao.storage == nil {
		dao.storage = make(map[string]string)
	}
	dao.storage[shortLink] = fullLink
	return nil
}


func (dao *HashMapStorage) InitStorage(storage map[string]string){
	// Is needed for unit-tests
	dao.storage = storage
}