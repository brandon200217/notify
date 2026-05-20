package channel

import "fmt"

type Registry struct {
	channels map[string]Channel
}

func NewRegistry() *Registry {
	return &Registry{
		channels: make(map[string]Channel),
	}
}

func (r *Registry) Register(channelType string, ch Channel) {
	r.channels[channelType] = ch
}

func (r *Registry) Get(channelType string) (Channel, error) {
	ch, ok := r.channels[channelType]
	if !ok {
		return nil, fmt.Errorf("canal no registrado: %s", channelType)
	}
	return ch, nil
}
