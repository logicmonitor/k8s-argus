package rest

import (
	"fmt"
	"net/url"
)

func deviceGroupName(cluster string) string {
	groupName := fmt.Sprintf("Kubernetes Cluster: %s", cluster)
	groupName = url.QueryEscape(groupName)
	return groupName
}

func collectorGroupName(cluster string) string {
	groupName := fmt.Sprintf("Kubernetes Cluster: %s", cluster)
	groupName = url.QueryEscape(groupName)
	return groupName
}

func dashboardGroupName(cluster string) string {
	return url.QueryEscape(fmt.Sprintf("Kubernetes Cluster: %s Dashboards", cluster))
}

func serviceGroupName(cluster string) string {
	return url.QueryEscape(fmt.Sprintf("Kubernetes Cluster: %s Services", cluster))
}
