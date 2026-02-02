package core

import (
	"fmt"
	"goblockchain/types"
	"sync"
)

type State struct {
	data map[string][]byte
}
type AccountState struct {
	mu   sync.RWMutex
	data map[types.Address]uint64
}

func NewAccountState() *AccountState {
	return &AccountState{
		data: make(map[types.Address]uint64),
	}
}
func (s *AccountState) GetBalance(key types.Address) (uint64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	balance, ok := s.data[key]
	if !ok {
		return 0, fmt.Errorf("Account not found")
	}
	return balance, nil
}

func (s *AccountState) AddBalance(addr types.Address, balance uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[addr]
	if !ok {
		s.data[addr] = balance
	} else {
		s.data[addr] += balance
	}
	return nil
}
func (s *AccountState) RemoveBalance(addr types.Address, balance uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	amount, ok := s.data[addr]
	if !ok {
		return fmt.Errorf("Account not found")
	}
	if amount > balance {
		return fmt.Errorf("Account balance is too low")
	}
	s.data[addr] -= amount
	return nil
}
func NewState() *State {
	return &State{make(map[string][]byte)}
}
func (s *State) Put(k, v []byte) error {
	s.data[string(k)] = v
	return nil
}
func (s *State) Delete(k []byte) error {
	delete(s.data, string(k))
	return nil
}
func (s *State) Get(k []byte) ([]byte, error) {
	if v, ok := s.data[string(k)]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("key %s not found", string(k))
}
