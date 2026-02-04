package datastorage

type DataStorage interface {
	Get(uuid string) (int64, error)
	ChangeBalance(sum int64, uuid string) error
}
