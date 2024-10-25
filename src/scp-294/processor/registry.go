package processor

import (
	"errors"
	"sync"
)

type IProcessor interface {
	GetName() string
	GetType() ProcType
	Process(input []byte) ([]byte, error)
}

const (
	IbmEBCDICEncoder = "IbmEBCDICEncoder"
	IbmEBCDICDecoder = "IbmEBCDICDecoder"
)

type ProcType int

const (
	Encode ProcType = iota
	Decode
	Compress
	Decompress
	Encrypt
	Decrypt
)

type Registry struct {
	processors map[string]IProcessor
	mu         sync.RWMutex
}

func NewProcessorRegistry() *Registry {
	registry := &Registry{
		processors: make(map[string]IProcessor),
	}

	registry.registerProcessor(NewIbmEBCDIC(IbmEBCDICDecoder, Decode))
	registry.registerProcessor(NewIbmEBCDIC(IbmEBCDICEncoder, Encode))

	return registry
}

func (r *Registry) registerProcessor(processor IProcessor) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.processors[processor.GetName()] = processor
}

func (r *Registry) GetProcessor(name string) (IProcessor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	processor, exists := r.processors[name]
	if !exists {
		return nil, errors.New("processor not found, name: " + name)
	}
	return processor, nil
}
