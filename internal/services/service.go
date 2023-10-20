package services

var (
	serviceInstance *service
)

func GetService() (instance *service, err error) {
	if serviceInstance == nil {
		instance, err = newService()
		if err != nil {
			return nil, err
		}
		serviceInstance = instance
	}
	return serviceInstance, nil
}

type service struct {
	Block BlockService
	Txn   TxnService
}

func newService() (instance *service, err error) {
	block, err := NewBlockService()
	if err != nil {
		return
	}

	txn, err := NewTxnService()
	if err != nil {
		return
	}

	return &service{
		Block: block,
		Txn:   txn,
	}, nil
}
