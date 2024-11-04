package processor

import (
	"sync"
)

type ProcRegistry struct {
	processors map[string]IProcessor
	mu         sync.RWMutex
}

func (r *ProcRegistry) RegisterProcessor(processor IProcessor) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.processors[processor.GetName()] = processor
}

func (r *ProcRegistry) GetProcessor(name string) IProcessor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	processor, exists := r.processors[name]
	if !exists {
		return nil
	}
	return processor
}

func (r *ProcRegistry) GetProcessorNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.processors))
	for name := range r.processors {
		names = append(names, name)
	}
	return names
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
