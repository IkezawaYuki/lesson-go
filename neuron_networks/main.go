package neuron_networks

type NeuronInterface interface {
	Iter() []*Neuron
}

type Neuron struct {
	In, Out []*Neuron
}

func (n *Neuron) Iter() []*Neuron {
	return []*Neuron{n}
}

type NeuronLayer struct {
	Neuron []Neuron
}

func (n *NeuronLayer) Iter() []*Neuron {
	// todo
}

func (n *Neuron) ConnectTo(other *Neuron) {
	n.Out = append(n.Out, other)
	other.In = append(other.In, n)
}
