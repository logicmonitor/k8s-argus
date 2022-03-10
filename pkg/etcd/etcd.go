package etcd

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/etcd/client"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Controller is the etcd controller for discovering etcd nodes.
type Controller struct {
	types.ResourceManager
}

// Member is a discovered etcd member.
type Member struct {
	Name string
	URL  *url.URL
}

// DiscoverByToken discovers the etcd node IP addresses using the etcd discovery service.
func (c *Controller) DiscoverByToken() ([]*Member, error) {
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"name": "etcd-discovery"}))
	log := lmlog.Logger(lctx)
	members := make([]*Member, 0)
	conf, err := config.GetConfig(lctx)
	if err != nil {
		log.Errorf("Failed to get config")
		return nil, err
	}
	response, err := http.Get(conf.EtcdDiscoveryToken) // nolint: bodyclose,noctx
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("Failed to close response body: %s", err)
		}
	}(response.Body)
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
		if len(s) == 2 { // nolint: gomnd
			u, err := url.Parse(s[1])
			if err != nil {
				return nil, err
			}
			m := &Member{
				Name: s[0],
				URL:  u,
			}
			members = append(members, m)
			c.addResource(lctx, m)
		}
	}

	return members, nil
}

func (c *Controller) addResource(lctx *lmctx.LMContext, member *Member) {
	log := lmlog.Logger(lctx)
	log.Tracef("Adding ETCD member: %s", member.Name)
	// Check if the etcd member has already been added.
	d, err := c.FindByDisplayName(lctx, enums.ETCD, fmtMemberDisplayName(member))
	if err != nil {
		log.Errorf("Failed to find etcd member %q: %v", member.Name, err)

		return
	}

	if d != nil {
		log.Warnf("ETCD member with name %s already exists", member.Name)
		return
	}

	rt := enums.ETCD
	// Add the etcd member.
	// nolint: exhaustivestruct
	if _, err := c.AddFunc()(lctx, enums.ETCD, metav1.ObjectMeta{},
		c.args(member, rt.GetCategory())...,
	); err != nil {
		log.Errorf("Failed to add etcd member %q: %v", member.URL.Hostname(), err)

		return
	}

	log.Infof("Added etcd member %q", member.Name)
}

// nolint: unparam
func (c *Controller) args(member *Member, category string) []types.ResourceOption {
	return []types.ResourceOption{
		c.Name(member.URL.Hostname()),
		c.DisplayName(fmtMemberDisplayName(member)),
		c.SystemCategory(category, enums.Add),
		c.Auto("clientport", member.URL.Port()),
	}
}

func fmtMemberDisplayName(member *Member) string {
	return member.Name + "-" + member.URL.Hostname()
}
