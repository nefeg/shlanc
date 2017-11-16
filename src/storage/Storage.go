package storage


type Storage interface {

	Connect() (isConnected bool)
	Disconnect()

	Get(index string) (record string)
	Add(index string, record string, force bool) (result bool, err error)
	Rm(index string) (result bool, err error)
	List() (data map[string]string)
	Flush()
}