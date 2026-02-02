package network

import (
	"github.com/stretchr/testify/assert"
	"goblockchain/core"
	"math/rand"
	"strconv"
	"testing"
)

func TestTxPool(t *testing.T) {
	pool := NewTxPool()
	assert.Equal(t, pool.Len(), 0)
}
func TestTxPool_Add(t *testing.T) {
	p := NewTxPool()
	tx := core.NewTransaction([]byte("hello world"))
	assert.Nil(t, p.Add(tx))
	assert.Equal(t, p.Len(), 1)
	p.Flush()
	assert.Equal(t, p.Len(), 0)
}
func TestSortTransactions(t *testing.T) {
	p := NewTxPool()
	txLen := 1000
	for i := 0; i < txLen; i++ {
		tx := core.NewTransaction([]byte(strconv.FormatInt(int64(i), 10)))
		tx.SetFirstSeen(int64(i * rand.Intn(1000)))
		assert.Nil(t, p.Add(tx))
	}
	assert.Equal(t, p.Len(), txLen)
	txx := p.Transactions()
	for i := 0; i < len(txx)-1; i++ {
		assert.True(t, txx[i].FirstSeen() <= txx[i+1].FirstSeen())
	}
}
