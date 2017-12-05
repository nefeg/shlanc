package hrentabd


type Storage interface {

	Connect() (isConnected bool)
	Disconnect()

	Pull(index string) (record string)
	Push(index string, record string, force bool) (result bool, err error)
	List() (data map[string]string)
	Version() (version string)
	Flush()
}