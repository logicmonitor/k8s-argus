package utilities

import (
	"strings"
	"testing"
)

func TestGetBatchDisplayNames(t *testing.T) {
	startIndex := 0
	names := GetBatchDisplayNames("displayName", "clusterName", 3, &startIndex)
	t.Logf("names: %v", strings.Join(names, " | "))
	names = GetBatchDisplayNames("displayName", " cluster Name ", 3, &startIndex)
	t.Logf("names: %v", strings.Join(names, " | "))
	names = GetBatchDisplayNames("displayName", "   cluster       Name    ", 3, &startIndex)
	t.Logf("names: %v", strings.Join(names, " | "))
}
