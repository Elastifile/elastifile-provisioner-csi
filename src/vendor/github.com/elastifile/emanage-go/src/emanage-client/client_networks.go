package emanage

import (
	"fmt"

	"rest"
)

const clientNetworksUri = "/api/client_networks"

type clientNetworks struct {
	conn *rest.Session
}

type ClientNetwork struct {
	Id                 int      `json:"id"`
	Name               string   `json:"name"`
	Vlan               int      `json:"vlan"`
	Subnet             string   `json:"subnet" default:"172.16.0.0"`
	Range              int      `json:"range" default:"16"`
	IpAddresses        []string `json:"ip_addresses" doc:"VIPs for NFS access"`
	PrivateIpAddresses []string `json:"private_ip_addresses" doc:"vHead networking, >= #enodes"`
	Mtu                int      `json:"mtu"`
	Url                string   `json:"url"`
}

func (cn *clientNetworks) GetAll() (result []ClientNetwork, err error) {
	err = cn.conn.Request(rest.MethodGet, clientNetworksUri, nil, &result)
	return result, err
}

func (cn *clientNetworks) GetById(id int) (ClientNetwork, error) {
	var result ClientNetwork
	fullUri := fmt.Sprintf("%s/%d", clientNetworksUri, id)
	err := cn.conn.Request(rest.MethodGet, fullUri, nil, &result)
	return result, err
}

func (cn *clientNetworks) Create(opts *ClientNetwork) (ClientNetwork, error) {
	var result ClientNetwork
	err := cn.conn.Request(rest.MethodPost, clientNetworksUri, opts, &result)
	return result, err
}

func (cn *clientNetworks) Update(opts *ClientNetwork) (result ClientNetwork, err error) {
	if opts.Id == 0 {
		err = fmt.Errorf("missing client network id for update")
	} else {
		fullUri := fmt.Sprintf("%s/%d", clientNetworksUri, opts.Id)
		err = cn.conn.Request(rest.MethodPut, fullUri, opts, &result)
	}
	return result, err
}

func (cn *clientNetworks) Delete(id int) (ClientNetwork, error) {
	var result ClientNetwork
	fullUri := fmt.Sprintf("%s/%d", clientNetworksUri, id)
	err := cn.conn.Request(rest.MethodDelete, fullUri, nil, &result)
	return result, err
}

func (cn *clientNetworks) GetDefaultClientNetwork() (cnw ClientNetwork, err error) {
	clientNetworkSlice, err := cn.GetAll()
	if err != nil {
		return cnw, fmt.Errorf("Error retrieving client network slice from client networks api. err=%v", err)
	}
	for _, cnws := range clientNetworkSlice {
		return cnws, nil
	}
	return cnw, fmt.Errorf("No client networks from client networks api")
}

func (cn *ClientNetwork) GetDefaultFrontendIp() (string, error) {
	if len(cn.IpAddresses) == 0 {
		return "", fmt.Errorf("No frontend ip addresses found")
	}
	for _, ip := range cn.IpAddresses {
		// TODO: add extra ip address validations
		if len(ip) > 0 {
			// found valid ip, return
			return ip, nil
		}
	}
	return "", fmt.Errorf("No valid frontend ip addresses found")
}
