package core

import (
	"errors"
	"sync"
)

type TxnPool struct {
	//TODO
	usedTxn   sync.Map
	unusedTxn sync.Map
}

const (
	TxnUsed   = "Txn already used"
	TxnInPool = "Txn already in pool"
)

func NewTxnPool() *TxnPool {
	return &TxnPool{}
}
func (tp *TxnPool) AddTxn(txn *Transaction) error {
	if _, exists := tp.usedTxn.Load(txn.hash); exists {
		return errors.New(TxnUsed)
	}
	_, loaded := tp.unusedTxn.LoadOrStore(txn.hash, txn)
	if loaded {
		return errors.New(TxnInPool)
	}
	return nil
}

func (tp *TxnPool) ExtractTxn() ([]*Transaction, error) {
	txns := make([]*Transaction, 0)
	ch := make(chan *Transaction)
	var wg sync.WaitGroup
	// 启动收集goroutine
	go func() {
		for txn := range ch {
			txns = append(txns, txn)
		}
	}()
	tp.unusedTxn.Range(func(_, value any) bool {
		wg.Add(1)
		go func(txn *Transaction) {
			defer wg.Done()
			ch <- txn
		}(value.(*Transaction))
		return true
	})

	wg.Wait()
	close(ch)
	return txns, nil
}

func (tp *TxnPool) ApplyBlock(block *Block) error {
	for _, tx := range block.Transactions {
		if _, exists := tp.unusedTxn.Load(tx.hash); exists {
			tp.unusedTxn.Delete(tx.hash)
		}
		tp.usedTxn.Store(tx.hash, tx)
	}
	return nil
}
