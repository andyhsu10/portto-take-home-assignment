package controllers

var (
	controllerInstance *controller
)

type controller struct {
	Block BlockController
	Txn   TxnController
}

func GetController() (instance *controller, err error) {
	if controllerInstance == nil {
		instance, err = newController()
		if err != nil {
			return nil, err
		}
		controllerInstance = instance
	}
	return controllerInstance, nil
}

func newController() (instance *controller, err error) {
	block, err := NewBlockController()
	if err != nil {
		return
	}

	txn, err := NewTxnController()
	if err != nil {
		return
	}

	return &controller{
		Block: block,
		Txn:   txn,
	}, nil
}
