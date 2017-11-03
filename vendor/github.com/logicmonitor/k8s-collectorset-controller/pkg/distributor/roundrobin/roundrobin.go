package roundrobin

import (
	"fmt"
	"sync"

	"github.com/logicmonitor/k8s-collectorset-controller/api"
)

// Distributor is a struct that implements the following interfaces:
// -   distributor.DistributionStratgey
type Distributor struct {
	ids []int32
	sync.Mutex
}

// ID implements distributor.DistributionStratgey
func (d *Distributor) ID(req *api.CollectorIDRequest) (*api.CollectorIDReply, error) {
	d.Lock()
	defer d.Unlock()

	if len(d.ids) == 0 {
		return nil, fmt.Errorf("No available collectors")
	}

	i, j := d.ids[0], d.ids[1:]
	d.ids = append(j, i)
	reply := &api.CollectorIDReply{Id: i}

	return reply, nil
}

// SetIDs implements distributor.DistributionStratgey
func (d *Distributor) SetIDs(ids []int32) error {
	d.ids = ids
	return nil
}
