package keyvalue

//go:generate mockgen -package=mock -destination=mock/interfaces.go -source=interfaces.go

type Storage interface {
	Load(key string) (interface{}, bool)
	Remove(key string)
	Exist(key string) bool
	Save(key string, value interface{}) error
	LoadAll() (map[string]interface{}, error)
}
