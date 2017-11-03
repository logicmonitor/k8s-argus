package distributor

import "encoding/json"

const (
	roundRobinName = "RoundRobin"
)

// Type is an enum for specifying a distributor type.
type Type int

const (
	// RoundRobin is the name of the round robin distribution strategy.
	RoundRobin Type = iota
)

func (t *Type) String() string {
	switch *t {
	case RoundRobin:
		return roundRobinName
	default:
		return ""
	}
}

// MarshalJSON implements json.Marshaler
func (t *Type) MarshalJSON() ([]byte, error) {
	var aux string

	switch *t {
	case RoundRobin:
		aux = roundRobinName
	}

	b, err := json.Marshal(aux)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// UnmarshalJSON implements json.Unmarshaler
func (t *Type) UnmarshalJSON(b []byte) error {
	var aux string
	err := json.Unmarshal(b, &aux)
	if err != nil {
		return err
	}

	switch aux {
	case roundRobinName:
		rr := RoundRobin
		t = &rr // nolint: ineffassign
	}

	return nil
}
