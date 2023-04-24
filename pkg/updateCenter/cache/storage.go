package storage

type CacheStorage interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
}

type FileCacheStorage struct {
	cacheDir string
}

func (f FileCacheStorage) Get(key string) ([]byte, error) {
	return []byte{}, nil
}
func (f FileCacheStorage) Set(key string, value []byte) error {
	return nil
}
func (f FileCacheStorage) Delete(key string) error {
	return nil
}

func SetCacheType(storage CacheStorage) CacheStorage {
	return storage
}
