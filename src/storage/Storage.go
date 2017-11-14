package storage

type Index string
type Record string

type Storage interface {

	Connect() (isConnected bool)
	Disconnect()

	Get(index Index) (record Record)
	Add(index Index, record Record, force bool) (result bool, err error)
	Rm(index Index) (result bool, err error)
	List() (data map[Index]Record)
	Flush()
}