package dao


type HashMapDAO struct {
	storage map[string]string
}

func (dao *HashMapDAO) GetURLByShortLink(shortLink string) (string, error) {

	if dao.storage == nil {
		dao.storage = make(map[string]string)
	}

	if value, ok := dao.storage[shortLink]; ok {
		return value, nil
	} else {
		return value, ErrURLMappingNotFound
	}
}

func (dao *HashMapDAO) AddShortURL(shortLink string, fullLink string) error {
	if dao.storage == nil {
		dao.storage = make(map[string]string)
	}
	dao.storage[shortLink] = fullLink
	return nil
}
