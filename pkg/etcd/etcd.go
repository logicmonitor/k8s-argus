package etcd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/etcd/client"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	log "github.com/sirupsen/logrus"
)

// Controller is the etcd controller for discovering etcd nodes.
type Controller struct {
	types.DeviceManager
}

// Member is a discovered etcd member.
type Member struct {
	Name string
	URL  *url.URL
}

// DiscoverByToken discovers the etcd node IP addresses using the etcd discovery service.
func (c *Controller) DiscoverByToken() ([]*Member, error) {
	members := []*Member{}
	response, err := http.Get(c.Config().EtcdDiscoveryToken)
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
	// Check if the etcd member has already been added.
	d, err := c.FindByDisplayName(fmtMemberDisplayName(member))
	if err != nil {
		log.Errorf("Failed to find etcd member %q: %v", member.Name, err)
		return
	}

	if d != nil {
		return
	}

	// Add the etcd member.
	if _, err := c.Add(
		c.args(member, constants.EtcdCategory)...,
	); err != nil {
		log.Errorf("Failed to add etcd member %q: %v", member.URL.Hostname(), err)
		return
	}

	log.Infof("Added etcd member %q", member.Name)
}

// nolint: unparam
func (c *Controller) args(member *Member, category string) []types.DeviceOption {
	categories := utilities.BuildSystemCategoriesFromLabels(category, nil)
	return []types.DeviceOption{
		c.Name(member.URL.Hostname()),
		c.DisplayName(fmtMemberDisplayName(member)),
		c.SystemCategories(categories),
		c.Auto("clientport", member.URL.Port()),
	}
}

func fmtMemberDisplayName(member *Member) string {
	return member.Name + "-" + member.URL.Hostname()
}
