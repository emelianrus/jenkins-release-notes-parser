package redisStorage

func (r *RedisStorage) Status() error {
	return r.DB.Status()
}
