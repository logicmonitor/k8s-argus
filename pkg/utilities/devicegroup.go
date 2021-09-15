package utilities

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/lm-sdk-go/models"
)

// BuildResourceGroup build
func BuildResourceGroup(lctx *lmctx.LMContext, d *models.DeviceGroup, options ...types.ResourceGroupOption) (*models.DeviceGroup, error) {
	if d == nil {
		// Adding auto prop on resource created by argus
		key := constants.DGCustomPropCreatedBy
		val := fmt.Sprintf("%s%s", constants.CreatedByPrefix, constants.Version)
		d = &models.DeviceGroup{ // nolint: exhaustivestruct
			CustomProperties: []*models.NameAndValue{
				{
					Name:  &key,
					Value: &val,
				},
			},
		}

		for _, option := range options {
			option(d)
		}
	} else {
		for _, option := range options {
			option(d)
		}
	}

	return d, nil
}
