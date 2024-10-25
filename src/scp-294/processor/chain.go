package processor

import "errors"

type Chain struct {
	processors []IProcessor
}

func NewProcessorChain() *Chain {
	return &Chain{
		processors: []IProcessor{},
	}
}

func (chain *Chain) AddProcessor(processor IProcessor) {
	chain.processors = append(chain.processors, processor)
}

func (chain *Chain) RemoveProcessor(name string) error {
	for i, processor := range chain.processors {
		if processor.GetName() == name {
			chain.processors = append(chain.processors[:i], chain.processors[i+1:]...)
			return nil
		}
	}
	return errors.New("processor not found in chain, name: " + name)
}

func (chain *Chain) Process(input []byte) ([]byte, error) {
	var err error
	for _, processor := range chain.processors {
		input, err = processor.Process(input)
		if err != nil {
			return nil, err
		}
	}
	return input, nil
}
