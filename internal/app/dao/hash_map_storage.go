package dao


import "fmt"


type HashMapDAO struct {
	storage map[string]string
}

func (dao *HashMapDAO) GetURLByShortLink(shortLink string) (string, error) {
	fmt.Println(dao.storage)

	if dao.storage == nil {
		dao.storage = make(map[string]string)
	}

	if value, ok := dao.storage[shortLink]; ok {
		fmt.Println(value, ok)
		return value, nil
	} else {
		fmt.Println(value, ok)
		return value, ErrURLMappingNotFound
	}
}

func (dao *HashMapDAO) AddShortURL(shortLink string, fullLink string) error {
	fmt.Println("Before", dao.storage)
	if dao.storage == nil {
		dao.storage = make(map[string]string)
	}
	dao.storage[shortLink] = fullLink
	fmt.Println("After", dao.storage)
	return nil
}
