package processor

import "errors"

type ProcChain struct {
	processors []IProcessor
}

func NewProcessorChain(procNames []string) (*ProcChain, error) {
	chain := &ProcChain{
		processors: []IProcessor{},
	}

	for _, name := range procNames {
		processor, err := DefaultProcRegistry.GetProcessor(name)
		if err != nil {
			return nil, err
		}
		chain.AddProcessor(processor)
	}
	return chain, nil
}

func (chain *ProcChain) AddProcessor(processor IProcessor) {
	chain.processors = append(chain.processors, processor)
}

func (chain *ProcChain) RemoveProcessor(name string) error {
	for i, processor := range chain.processors {
		if processor.GetName() == name {
			chain.processors = append(chain.processors[:i], chain.processors[i+1:]...)
			return nil
		}
	}
	return errors.New("processor not found in chain, name: " + name)
}

func (chain *ProcChain) Process(arr []byte) ([]byte, error) {
	var err error
	for _, processor := range chain.processors {
		arr, err = processor.Process(arr)
		if err != nil {
			return nil, err
		}
	}
	return arr, nil
}
