package processor

import (
	"errors"
	"sync"
)

type ProcRegistry struct {
	processors map[string]IProcessor
	mu         sync.RWMutex
}

var DefaultProcRegistry = createDefaultProcRegistry()

func createDefaultProcRegistry() *ProcRegistry {
	registry := &ProcRegistry{
		processors: make(map[string]IProcessor),
	}

	registry.RegisterProcessor(NewIbmEBCDIC(IbmEBCDICDecoder, Decode))
	registry.RegisterProcessor(NewIbmEBCDIC(IbmEBCDICEncoder, Encode))

	return registry
}

func (r *ProcRegistry) RegisterProcessor(processor IProcessor) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.processors[processor.GetName()] = processor
}

func (r *ProcRegistry) GetProcessor(name string) (IProcessor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	processor, exists := r.processors[name]
	if !exists {
		return nil, errors.New("processor not found, name: " + name)
	}
	return processor, nil
}
