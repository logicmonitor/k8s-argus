package enums

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// String returns string name
func (resourceType ResourceType) String() string {
	bytes, err := json.Marshal(resourceType)
	if err != nil {
		return ""
	}
	if str, err := strconv.Unquote(fmt.Sprintf(`%s`, bytes)); err == nil {
		return str
	}

	return ""
}

// FQName returns string name
func (resourceType *ResourceType) FQName(name string) string {
	if apiGroup := resourceType.APIGroup(); apiGroup != "" {
		return fmt.Sprintf("%s.%s/%s", resourceType.Singular(), apiGroup, name)
	}

	return fmt.Sprintf("%s/%s", resourceType.Singular(), name)
}

// Singular returns string name
func (resourceType *ResourceType) Singular() string {
	bytes, err := json.Marshal(resourceType)
	if err != nil {
		return ""
	}
	if str, err := strconv.Unquote(fmt.Sprintf(`%s`, bytes)); err == nil {
		return strings.TrimSuffix(str, "s")
	}

	return ""
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (resourceType *ResourceType) UnmarshalText(text []byte) error {
	l, err := ParseResourceType(string(text))
	if err != nil {
		return err
	}

	*resourceType = l

	return nil
}

// LMName lmname
func (resourceType ResourceType) LMName(meta *metav1.PartialObjectMetadata) string {
	s := ShortResourceType(resourceType)
	if resourceType != Namespaces && resourceType.IsNamespaceScopedResource() {
		return fmt.Sprintf("%s-%s-%s", meta.Name, s.String(), meta.Namespace)
	}

	return fmt.Sprintf("%s-%s", meta.Name, s.String())
}

// ShortResourceType to specifically use as short resource type
type ShortResourceType ResourceType

// UnmarshalText implements encoding.TextUnmarshaler.
func (resourceType *ShortResourceType) UnmarshalText(text []byte) error {
	l, err := ParseShortResourceType(string(text))
	if err != nil {
		return err
	}
	*resourceType = l

	return nil
}

// String returns string name
func (resourceType *ShortResourceType) String() string {
	bytes, err := json.Marshal(resourceType)
	if err != nil {
		return ""
	}
	if str, err := strconv.Unquote(fmt.Sprintf(`%s`, bytes)); err == nil {
		return str
	}

	return ""
}

// GetConflictsCategory returns category name for conflicts group
func (resourceType *ResourceType) GetConflictsCategory() string {
	return fmt.Sprintf("%s%s", resourceType.GetCategory(), "Conflict")
}

// GetDeletedCategory returns category name for conflicts group
func (resourceType *ResourceType) GetDeletedCategory() string {
	return fmt.Sprintf("%s%s", resourceType.GetCategory(), "Deleted")
}

// GetCategory returns category name for group
func (resourceType *ResourceType) GetCategory() string {
	return fmt.Sprintf("%s%s", "Kubernetes", resourceType.Title())
}

// ObjectMeta returns object meta from interface object
func (resourceType *ResourceType) ObjectMeta(obj interface{}) (*metav1.PartialObjectMetadata, error) {
	if *resourceType == ETCD || *resourceType == Unknown {
		return nil, fmt.Errorf("no metadata for resource: %s", *resourceType)
	}
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}
	return meta.AsPartialObjectMetadata(accessor), nil
}
