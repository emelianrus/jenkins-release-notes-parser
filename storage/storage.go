package storage

type Storage interface {
	Get(key string) ([]byte, error)
	Set(key string, value interface{}) error
	Keys(key string) ([]string, error)
}

func SetStorage(db Storage) Storage {
	return db
}
