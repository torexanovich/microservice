package repo

type RedisRepo interface {
	SetWithTTL(key, value string, seconds int) error
	Get(key string) (interface{}, error)
}