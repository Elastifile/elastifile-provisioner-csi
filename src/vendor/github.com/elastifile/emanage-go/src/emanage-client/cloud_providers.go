package emanage

import (
	"fmt"

	"rest"
	"types"
)

const cloudProviderUri = "api/cloud_providers"

type cloudProviders struct {
	conn *rest.Session
}

func (clp *cloudProviders) GetById(id int) (*types.CloudProvider, error) {
	var result types.CloudProvider
	fullUri := fmt.Sprintf("%s/%d", cloudProviderUri, id)
	err := clp.conn.Request(rest.MethodGet, fullUri, nil, &result)

	return &result, err
}

func (clp *cloudProviders) GetAll() (*[]types.CloudProvider, error) {
	result := []types.CloudProvider{}
	err := clp.conn.Request(rest.MethodGet, cloudProviderUri, nil, &result)
	return &result, err
}

func (clp *cloudProviders) Update(id int, opt *types.CloudProviderUpdateOpts) (types.CloudProvider, error) {
	Name := "cloud_provider"
	params := struct {
		Name string `json:"name,omitempty"`
		types.CloudProviderUpdateOpts
	}{Name, *opt}

	fullUri := fmt.Sprintf("%s/%d", cloudProviderUri, id)
	var result types.CloudProvider
	err := clp.conn.Request(rest.MethodPut, fullUri, params, &result)
	return result, err
}
