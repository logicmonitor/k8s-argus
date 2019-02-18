package node

import (
	"testing"

	"k8s.io/api/core/v1"
)

func TestGetInternalAddress(t *testing.T) {
	var addresses []v1.NodeAddress
	address := getInternalAddress(addresses)
	if address != nil {
		t.Errorf("invalid address: %v", address)
	}
	addresses = append(addresses, v1.NodeAddress{
		Type:    v1.NodeHostName,
		Address: "test",
	})

	address = getInternalAddress(addresses)
	if address == nil || address.Address != "test" {
		t.Errorf("invalid address: %v", address)
	}

	addresses = append(addresses, v1.NodeAddress{
		Type:    v1.NodeInternalIP,
		Address: "127.0.0.1",
	})

	address = getInternalAddress(addresses)
	if address == nil || address.Address != "127.0.0.1" {
		t.Errorf("invalid address: %v", address)
	}
}
