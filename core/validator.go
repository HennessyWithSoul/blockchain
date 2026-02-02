package core

import "fmt"

type Validator interface {
	ValidateBlock(*Block) error
}
type BlockValidator struct {
	bc *BlockChain
}

func NewBlockValidator(bc *BlockChain) *BlockValidator {
	return &BlockValidator{bc}
}
func (v *BlockValidator) ValidateBlock(b *Block) error {
	if v.bc.HasBlock(b.Height) {
		return fmt.Errorf("blockchain %d already exists", b.Height)
	}
	if b.Height != v.bc.Height()+1 {
		//fmt.Println("block too high-------------")
		return fmt.Errorf("block (%s) too high", b.Hash(BlockHasher{}))
	}
	prevHeader, err := v.bc.GetHeader(b.Height - 1)
	if err != nil {
		return err
	}
	hash := BlockHasher{}.Hash(prevHeader)
	if hash != b.PrevBlockHash {
		return fmt.Errorf("the hash of the previous block (%s) is invalid", b.PrevBlockHash)
	}
	if err := b.Verify(); err != nil {
		return err
	}
	return nil
}
