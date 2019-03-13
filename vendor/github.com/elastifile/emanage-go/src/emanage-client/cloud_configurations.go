package emanage

import (
	"fmt"

	"rest"
	"types"
)

const cloudConfigurationUri = "api/cloud_configurations"

type cloudConfigurations struct {
	conn *rest.Session
}

func (clc *cloudConfigurations) GetById(id int) (*types.CloudConfiguration, error) {
	var result types.CloudConfiguration
	fullUri := fmt.Sprintf("%s/%d", cloudConfigurationUri, id)
	err := clc.conn.Request(rest.MethodGet, fullUri, nil, &result)
	return &result, err
}

func (clc *cloudConfigurations) GetAll() (*[]types.CloudConfiguration, error) {
	result := []types.CloudConfiguration{}
	err := clc.conn.Request(rest.MethodGet, cloudConfigurationUri, nil, &result)
	return &result, err
}

func (clc *cloudConfigurations) Create(opt *types.CloudConfigurationCreateOpts) (types.CloudConfiguration, error) {
	var result types.CloudConfiguration
	err := clc.conn.Request(rest.MethodPost, cloudConfigurationUri, opt, &result)
	return result, err
}

func (clc *cloudConfigurations) Update(opt *types.CloudConfigurationUpdateOpts) (types.CloudConfiguration, error) {
	params := struct {
		types.CloudConfigurationUpdateOpts
	}{*opt}
	fullUri := fmt.Sprintf("%s/%d", cloudConfigurationUri, params.ID)
	var result types.CloudConfiguration
	err := clc.conn.Request(rest.MethodPut, fullUri, params, &result)
	return result, err
}
