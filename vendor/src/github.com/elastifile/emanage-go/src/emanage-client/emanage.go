// The emanage package provides access to the Emanage REST API.
//
// Object names and data types are closely related to the names exported by Emanage.
//
// NOTE: This package is far from final and needs to undergo a lot of changes,
// as outlined in comments throughout the package's code.
package emanage

import (
	"net/url"
	"path"
	"time"

	log "gopkg.in/inconshreveable/log15.v2"

	"logging"
	"rest"
)

var logger = logging.NewLogger("emanage")

type Client struct {
	ClusterReports      *clusterReports
	DataContainers      *dataContainers
	Enodes              *enodes
	ControlTasks        *controlTasks
	Exports             *exports
	Hosts               *hosts
	Events              *events
	NetworkInterfaces   *netInterfaces
	Policies            *policies
	Sessions            *rest.Session
	Devices             *devices
	Snapshots           *snapshots
	Statistics          *statistics
	Systems             *systems
	EmanageVMs          *emanageVms
	Tenants             *tenants
	VMManagers          *vmManagers
	VMs                 *vms
	Clients             *clients
	ClientNetworks      *clientNetworks
	ReplicationAgents   *replicationAgents
	RemoteSites         *remoteSites
	DcPairs             *dcPairs
	CloudProviders      *cloudProviders
	CloudConfigurations *cloudConfigurations

	log.Logger
}

type AsyncTasks struct {
	taskIDs []*rest.TaskID
	conn    *rest.Session
	err     error
}

func (tasks *AsyncTasks) Wait() error {
	if tasks.err != nil {
		return tasks.err
	}
	return tasks.conn.WaitAllTasks(tasks.taskIDs)
}

func (tasks *AsyncTasks) Error() error {
	return tasks.err
}

func EmanageURL(host string) *url.URL {
	baseURL := &url.URL{
		Scheme: "http",
		Host:   host,
	}
	return baseURL
}

func NewClient(baseURL *url.URL) *Client {
	s := rest.NewDefaultSession(baseURL)
	return &Client{
		ClusterReports:      &clusterReports{s},
		DataContainers:      &dataContainers{s},
		Enodes:              &enodes{s},
		Events:              &events{s},
		ControlTasks:        &controlTasks{s},
		Exports:             &exports{s},
		Hosts:               &hosts{s},
		NetworkInterfaces:   &netInterfaces{s},
		Policies:            &policies{s},
		Sessions:            s,
		Devices:             &devices{s},
		Statistics:          &statistics{s},
		Systems:             &systems{s},
		EmanageVMs:          &emanageVms{s},
		Snapshots:           &snapshots{s},
		Tenants:             &tenants{s},
		VMManagers:          &vmManagers{s},
		VMs:                 &vms{s},
		ClientNetworks:      &clientNetworks{s},
		Clients:             &clients{s},
		ReplicationAgents:   &replicationAgents{s},
		RemoteSites:         &remoteSites{s},
		DcPairs:             &dcPairs{s},
		CloudProviders:      &cloudProviders{s},
		CloudConfigurations: &cloudConfigurations{s},

		Logger: log.New("baseURL", baseURL.String()),
	}
}

func (client *Client) RetriedLogin(username string, password string) error {
	return client.RetriedLoginTimeout(username, password, 30*time.Second)
}

func (client *Client) RetriedLoginTimeout(username string, password string, timeout time.Duration) error {
	return client.Sessions.RetriedLoginTimeout(username, password, timeout)
}

func (client *Client) GetEMSMonitor() (EmanageMonitor, error) {
	var monitorResult EmanageMonitor
	uri := path.Join(emanageVMsUri, monitorUri)

	err := client.Sessions.Request(rest.MethodGet, uri, nil, &monitorResult)
	if err != nil {
		return monitorResult, err
	}
	return monitorResult, nil
}
