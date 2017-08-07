package etcd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/utilities"

	"github.com/coreos/etcd/client"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	lm "github.com/logicmonitor/lm-sdk-go"
	log "github.com/sirupsen/logrus"
)

// Controller is the etcd controller for discovering etcd nodes.
type Controller struct {
	*types.Base
}

// Member is a discovered etcd memeber.
type Member struct {
	Name string
	URL  *url.URL
}

// DiscoverByToken discovers the etcd node IP addresses using the etcd discovery service.
func (c *Controller) DiscoverByToken() ([]*Member, error) {
	members := []*Member{}
	response, err := http.Get(c.Config.EtcdDiscoveryToken)
	if err != nil {
		return nil, err
	}
	n := client.Response{}
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf, &n)
	if err != nil {
		return nil, err
	}
	log.Infof("Discovered %d etcd members", len(n.Node.Nodes))
	for _, member := range n.Node.Nodes {
		s := strings.Split(member.Value, "=")
		if len(s) == 2 {
			u, err := url.Parse(s[1])
			if err != nil {
				return nil, err
			}
			m := &Member{
				Name: s[0],
				URL:  u,
			}
			members = append(members, m)
			c.addDevice(m)
		}
	}

	return members, nil
}

func (c *Controller) addDevice(member *Member) {
	device := c.makeDeviceObject(member)
	restResponse, apiResponse, err := c.LMClient.AddDevice(device, false)
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		log.Errorf("Failed to add etcd host: %v", _err)
	}
	log.Infof("Added etcd member %q", member.Name)
}

func (c *Controller) makeDeviceObject(member *Member) (device lm.RestDevice) {
	categories := constants.EtcdCategory

	device = lm.RestDevice{
		Name:                 member.URL.Hostname(),
		DisplayName:          member.Name + "-" + c.Config.ClusterName,
		DisableAlerting:      c.Config.DisableAlerting,
		HostGroupIds:         "1",
		PreferredCollectorId: c.Config.PreferredCollector,
		CustomProperties: []lm.NameAndValue{
			{
				Name:  "system.categories",
				Value: categories,
			},
			{
				Name:  "auto.clustername",
				Value: c.Config.ClusterName,
			},
			{
				Name:  "auto.clientport",
				Value: member.URL.Port(),
			},
		},
	}

	return
}
