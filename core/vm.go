package core

import (
	"goblockchain/util"
)

type Instruction byte

const (
	InstrPushInt  Instruction = 0x0a //10
	InstrAdd      Instruction = 0x0b
	InstrPushByte Instruction = 0x0c
	InstrPack     Instruction = 0x0d
	InstrSub      Instruction = 0x0e
	InstrStore    Instruction = 0x0f
	InstrGet      Instruction = 0xae
)

type Stack struct {
	data []any //interface{}
	sp   int
}

func NewStack(size int) *Stack {
	return &Stack{
		data: make([]any, size),
		sp:   0,
	}
}
func (s *Stack) Push(v any) {
	s.data[s.sp] = v
	s.sp++
}
func (s *Stack) Pop() any {
	value := s.data[s.sp-1]
	s.sp--
	return value
}

type VM struct {
	data          []byte
	ip            int
	stack         *Stack
	contractState *State
}

func NewVM(data []byte, state *State) *VM {
	return &VM{
		contractState: state,
		data:          data,
		ip:            0, //instruction pointer
		stack:         NewStack(128),
	}
}
func (vm *VM) Run() error {
	for {
		instr := Instruction(vm.data[vm.ip])
		//fmt.Println("instr", instr)
		if err := vm.Exec(instr); err != nil {
			return err
		}
		vm.ip++
		if vm.ip >= len(vm.data) {
			break
		}
	}
	return nil
}
func (vm *VM) Exec(instr Instruction) error {
	switch instr {
	case InstrGet:
		var (
			key = vm.stack.Pop().([]byte)
		)
		value, err := vm.contractState.Get(key)
		if err != nil {
			return err
		}
		//fmt.Println("value", value)
		vm.stack.Push(value)
	case InstrStore:
		var serializedValue []byte
		value := vm.stack.Pop()
		key := vm.stack.Pop().([]byte)
		switch v := value.(type) {
		case int:
			serializedValue = util.SerialzeInt64(int64(v))
		default:
			panic("TODO:unknown type")
		}
		vm.contractState.Put(key, serializedValue)
		//fmt.Println("key", string(key), "value", serializedValue)
	case InstrPushInt:
		vm.stack.Push(int(vm.data[vm.ip-1]))
	case InstrPushByte:
		vm.stack.Push(byte(vm.data[vm.ip-1]))
	case InstrPack:
		n := vm.stack.Pop().(int)
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			b[n-i-1] = vm.stack.Pop().(byte)
		}

		vm.stack.Push(b)
	case InstrAdd:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		c := a + b
		vm.stack.Push(c)
	case InstrSub:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		c := b - a
		vm.stack.Push(c)
	}
	return nil
}
