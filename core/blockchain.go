package core

import (
	"fmt"
	"sync"
)

type BlockChain struct {
	store Storage
	//logger    log.Logger
	lock          sync.RWMutex
	accountStates *AccountState
	blocks        []Block
	headers       []*Header
	validator     Validator
	contractState *State
}

func NewBlockChain(genesis *Block) (*BlockChain, error) {
	bc := &BlockChain{
		contractState: NewState(),
		headers:       []*Header{},
		accountStates: NewAccountState(),
	}
	bc.validator = NewBlockValidator(bc)
	bc.store = NewMemoryStorage()
	err := bc.addBlockWithoutValidation(genesis)
	return bc, err
}
func (bc *BlockChain) SetValidator(validator Validator) {
	bc.validator = validator
}
func (bc *BlockChain) AddBlock(b *Block) error {

	if err := bc.validator.ValidateBlock(b); err != nil {
		//fmt.Println("ok", bc.Height())
		return err
	}
	bc.lock.Lock()
	defer bc.lock.Unlock()
	//run vm code
	//for _, tx := range b.Transactions {
	//vm := bc.NewVM(tx.Data)
	//vm.Run()
	//}
	for _, tx := range b.Transactions {
		if len(tx.Data) > 0 {
			vm := NewVM(tx.Data, bc.contractState)
			if err := vm.Run(); err != nil {
				return err
			}
		}
		if tx.TxInner != nil {
			if err := bc.handleNativeNFT(tx); err != nil {
				return err
			}
		}
		if tx.Value > 0 {
			fmt.Printf("--------------%+v\n", tx)
			if err := bc.handleNativeTransfer(tx); err != nil {
				fmt.Println("---------", err)
				return err
			}
		}
	}
	fmt.Println(" adding new block", bc.Height()+1)
	if err := bc.addBlockWithoutValidation(b); err != nil {
		return err
	}
	return nil

	//return nil
}
func (bc *BlockChain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}
func (bc *BlockChain) Height() uint32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return uint32(len(bc.headers) - 1)
}
func (bc *BlockChain) addBlockWithoutValidation(b *Block) error {
	//bc.lock.Lock()
	//defer bc.lock.Unlock()
	bc.headers = append(bc.headers, b.Header)
	bc.blocks = append(bc.blocks, *b)
	//log.Println("msg", "new block", "hash", b.Hash(BlockHasher{}), "height", b.Height, "transactions", len(b.Transactions))
	//bc.logger.Log("msg","new block","hash",b.Hash(BlockHasher{}),"height",b.Height,"transactions",len(b.Transactions))
	return bc.store.Put(b)

}
func (bc *BlockChain) GetHeader(height uint32) (*Header, error) {
	//fmt.Println(height, bc.Height())
	if height > bc.Height() {
		//log.Println("height is out of range", height)
		return nil, fmt.Errorf("given block header (%d) too high", height)
	}
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return bc.headers[height], nil
}
func (bc *BlockChain) GetBlock(height uint32) (*Block, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given block (%d) too high", height)
	}
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return &bc.blocks[height], nil
}
func (bc *BlockChain) AddGenesisBlock(b *Block) error {
	return nil
}
func (bc *BlockChain) handleNativeTransfer(tx *Transaction) error {
	if err := bc.accountStates.RemoveBalance(tx.From.Address(), tx.Value); err != nil {
		return err
	}
	return bc.accountStates.AddBalance(tx.To.Address(), tx.Value)
}
func (bc *BlockChain) handleNativeNFT(tx *Transaction) error {
	switch t := tx.TxInner.(type) {
	case CollectionTx:
		hash := tx.Hash(TxHasher{})
		bc.store.PutCollection(hash, &t)
	case MintTx:
		//coll, ok := bc.store.GetCollection(t.Collection)
		//if ok != nil {
		break
		//}
	default:
		fmt.Printf("unsupproted tx type %v", t)
	}
	return nil
}
