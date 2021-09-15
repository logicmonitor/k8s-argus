package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
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

func TestGetNodesMap(t *testing.T) {

	nodeTestCases := []struct {
		name               string
		clientSet          kubernetes.Interface
		input              string
		expectedOutput     string
		expectedNodesCount int
		err                string
	}{
		{
			name:               "No Node",
			clientSet:          fake.NewSimpleClientset(),
			expectedNodesCount: 0,
		},
		{
			name: "2 worker-nodes",
			clientSet: fake.NewSimpleClientset(&v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "worker-node1",
				},
				Spec: v1.NodeSpec{},
				Status: v1.NodeStatus{
					Addresses: []v1.NodeAddress{
						{
							Type:    v1.NodeInternalIP,
							Address: "172.31.18.239",
						},
					},
				},
			}, &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "worker-node2",
				},
				Spec: v1.NodeSpec{},
				Status: v1.NodeStatus{
					Addresses: []v1.NodeAddress{
						{
							Type:    v1.NodeInternalIP,
							Address: "172.31.18.240",
						},
					},
				},
			}),
			input:              "",
			expectedOutput:     "",
			expectedNodesCount: 2,
		},
	}

	assert := assert.New(t)
	for _, testCase := range nodeTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			nodesMap, err := GetNodesMap(testCase.clientSet)

			// check if err not nil
			if err != nil {
				assert.EqualError(err, testCase.err, "TestCase: \"%s\" \nResult: Expected error \"%s\" but got \"%s\"", testCase.name, testCase.err, err.Error())
			}

			// check expected nodes count
			assert.Equal(testCase.expectedNodesCount, len(nodesMap), "TestCase: \"%s\" \nResult: Expected pod count \"%d\" but got \"%d\"", testCase.name, testCase.expectedNodesCount, len(nodesMap))
		})
	}
}
