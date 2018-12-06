package emanage

import (
	"fmt"
	"path"
	"time"

	tm "github.com/buger/goterm"
	"github.com/go-errors/errors"

	"eurl"
	"helputils"
	"rest"
	"types"
)

var Timeout time.Duration = 10 * time.Minute

const (
	systemsUri = "api/systems"
	sysID      = 1
)

type systems struct {
	session *rest.Session
}

type SystemState string

const (
	StateSystemInit    SystemState = "system_init"
	StateConfigured    SystemState = "configured"
	StateMapped        SystemState = "mapped"
	StateInService     SystemState = "in_service"
	StateClosingGates  SystemState = "closing_gates"
	StateUpClosedGates SystemState = "up_closed_gates"
	StateShuttingDown  SystemState = "shutting_down"
	StateDown          SystemState = "down"
	StateLockdown      SystemState = "lockdown"
	StateUnknown       SystemState = "unknown"
)

type SystemDetails struct {
	Name             string             `json:"name"`
	Id               int                `json:"id"`
	Status           SystemState        `json:"status,omitempty"`
	ConnectionStatus string             `json:"connection_status,omitempty"` // TODO: have a type // "connection_ok"
	Uptime           string             `json:"uptime,omitempty"`
	Version          string             `json:"version,omitempty"`
	ReplicationLevel int                `json:"replication_level,omitempty"` // 2
	ControlAddress   string             `json:"control_address,omitempty"`   // TODO: use Host // "localhost"
	ControlPort      int                `json:"control_port,omitempty"`      // TODO: have type for ports // 10016
	NfsAddress       string             `json:"nfs_address,omitempty"`       // FIXME: use Host // "192.168.0.1"
	NfsIpRange       int                `json:"nfs_ip_range,omitempty"`
	DataAddress      string             `json:"data_address,omitempty"`      // "10.0.0.1"
	DataIpRange      int                `json:"data_ip_range,omitempty"`     // 16,
	DataVlan         int                `json:"data_vlan,omitempty"`         // null
	DataAddress2     string             `json:"data_address2,omitempty"`     // "10.0.0.1"
	DataIpRange2     int                `json:"data_ip_range2,omitempty"`    // 16,
	DataVlan2        int                `json:"data_vlan2,omitempty"`        // 5
	DataMtu          int                `json:"data_mtu,omitempty"`          // 9000
	DataMtu2         int                `json:"data_mtu2,omitempty"`         // 9000
	DeploymentModel  string             `json:"deployment_model,omitempty"`  // "hci"
	ExternalUseDhcp  bool               `json:"external_use_dhcp,omitempty"` // null
	ExternalAddress  string             `json:"external_address,omitempty"`  // null
	ExternalIpRange  string             `json:"external_ip_range,omitempty"` // null
	ExternalGateway  string             `json:"external_gateway,omitempty"`  // null
	ExternalNetwork  string             `json:"external_network,omitempty"`  // null
	CreatedAt        time.Time          `json:"created_at,omitempty"`
	UpdatedAt        time.Time          `json:"updated_at,omitempty"`
	Url              *eurl.URL          `json:"url,omitempty"` // Does not always appear -- only in "full" results
	TimeZone         string             `json:"time_zone,omitempty"`
	NTPServers       string             `json:"ntp_servers,omitempty"`
	UpgradeState     SystemUpgradeState `json:"upgrade_state"`
	UpgradePhase     SystemUpgradePhase `json:"upgrade_phase"`
	ShowWizard       bool               `json:"show_wizard,omitempty"`
	NameServer       string             `json:"name_server,omitempty"`
}
type StateError struct {
	Expected SystemState
	Actual   SystemState
}

func (e *StateError) Error() string {
	return fmt.Sprintf("Expected system state '%v', Actual is '%v'", e.Expected, e.Actual)
}

func (ss *systems) GetAll(opt *GetAllOpts) ([]SystemDetails, error) {
	if opt == nil {
		opt = &GetAllOpts{}
	}

	var result []SystemDetails
	return result, ss.session.Request(rest.MethodGet, systemsUri, opt, &result)
}

func (ss *systems) GetById(id int) (*System, *SystemDetails, error) {
	uri := fmt.Sprintf("%s/%d", systemsUri, id)
	var result SystemDetails
	err := ss.session.Request(rest.MethodGet, uri, nil, &result)
	if err != nil {
		return nil, nil, err
	}
	system := System{
		session: ss.session,
		id:      id,
	}
	return &system, &result, nil
}

func (ss *systems) MustGetById(id int) (*System, *SystemDetails) {
	sys, det, err := ss.GetById(id)
	if err != nil {
		panic(err)
	}
	return sys, det
}

func (ss *systems) Update(id int, sysInfo *SystemDetails) (*SystemDetails, error) {
	uri := fmt.Sprintf("%s/%d", systemsUri, id)

	var result SystemDetails
	err := ss.session.Request(rest.MethodPut, uri, sysInfo, &result)
	return &result, err
}

func (ss *systems) SystemDetailsFromDeployment(deployment *types.Deployment) *SystemDetails {
	elabSys := &deployment.System.Elab
	deploy := &deployment.System.Deploy
	nets := &elabSys.Data.Networks
	dataNet := &nets.DataNetwork[0]
	dataNet2 := &nets.DataNetwork[1]
	nfsNet := elabSys.NfsNetwork()
	mtu := 9000

	return &SystemDetails{
		ReplicationLevel: deploy.ReplicationLevel,
		NfsIpRange:       nfsNet.NetworkMask,
		DataAddress:      dataNet.NetworkId,
		DataIpRange:      dataNet.NetworkMask,
		DataVlan:         dataNet.VlanId,
		DataAddress2:     dataNet2.NetworkId,
		DataIpRange2:     dataNet2.NetworkMask,
		DataVlan2:        dataNet2.VlanId,
		DataMtu:          mtu,
		DataMtu2:         mtu,
		TimeZone:         deploy.TimeZone,
		NTPServers:       deploy.NTPServers,
		ExternalUseDhcp:  true,
	}
}

func (ss *systems) Create(id int, details *SystemDetails) (*SystemDetails, error) {
	uri := fmt.Sprintf("%s/%d", systemsUri, id)

	var result SystemDetails
	err := ss.session.Request(rest.MethodPut, uri, details, &result)

	return &result, err
}

type System struct {
	session *rest.Session
	id      int
}

func (s *System) anyRequest(method rest.HttpMethod, endpoint string, body interface{}, async bool, result interface{}) error {
	parts := []string{systemsUri, fmt.Sprintf("%d", s.id)}
	if endpoint != "" {
		parts = append(parts, endpoint)
	}
	uri := path.Join(parts...)

	if async {
		if tIDs, err := s.session.AsyncRequest(method, uri, body); err != nil {
			return err
		} else {
			return s.session.WaitAllTasks(tIDs)
		}
	} else {
		return s.session.Request(method, uri, body, result)
	}
}

func (s *System) anyRequestWithDetailsResponse(method rest.HttpMethod, endpoint string, body interface{}, async bool) (*SystemDetails, error) {
	var result SystemDetails

	if err := s.anyRequest(method, endpoint, body, async, &result); err != nil {
		return nil, err
	}

	if async {
		logger.Debug("Received new system state", "details", result)
	}

	return &result, nil
}

func (s *System) request(method rest.HttpMethod, endpoint string, body interface{}) (*SystemDetails, error) {
	return s.anyRequestWithDetailsResponse(method, endpoint, body, false)
}

func (s *System) asyncRequest(method rest.HttpMethod, endpoint string, body interface{}) (*SystemDetails, error) {
	return s.anyRequestWithDetailsResponse(method, endpoint, body, true)
}

func (s *System) GetDetails() (*SystemDetails, error) {
	return s.request(rest.MethodGet, "", nil)
}

func (s *System) logAction(msg string, kv ...interface{}) error {
	details, e := s.GetDetails()
	if e != nil {
		return e
	}

	kv = append(kv, "id", s.id, "name", details.Name)
	logger.Info(msg, kv...)

	return nil
}

func (s *System) ForceReset(reconfig bool) (*SystemDetails, error) {
	if e := s.logAction("Force reset system"); e != nil {
		return nil, e
	}
	params := struct {
		Async     bool `json:"async"`
		SkipTests bool `json:"skip_tests"`
		Reconfig  bool `json:"reconfig"`
	}{
		Async:     true,
		SkipTests: true,
		Reconfig:  reconfig,
	}
	return s.asyncRequest(rest.MethodPost, "force_reset", &params)
}

func (s *System) GetHealth() (*Health, error) {
	var health Health
	healthUri := "health"
	uri := path.Join(systemsUri, fmt.Sprintf("%d", s.id), healthUri)

	err := s.session.Request(rest.MethodGet, uri, nil, &health)
	if err != nil {
		return nil, err
	}
	return &health, nil
}

func (s *System) AcceptEULA() error {
	acceptEULAUri := "accept_eula"
	uri := path.Join(systemsUri, fmt.Sprintf("%d", s.id), acceptEULAUri)
	return s.session.Request(rest.MethodPost, uri, nil, nil)
}

type SystemStartOpts struct {
	SkipTests       bool
	MeltdownRecovey bool
	SkipSetup       bool
}

func (s *System) Start(opts SystemStartOpts) (*SystemDetails, error) {
	if e := s.logAction(tm.Bold("Start system"), helputils.MustStructToKeyValueInterfaces(opts)...); e != nil {
		return nil, e
	}

	if !opts.SkipSetup {
		// EL-671: We need to call Setup (which is idempotent) before calling Start
		_, err := s.Setup(nil, opts.SkipTests)
		if err != nil {
			return nil, err
		}
	}
	params := struct {
		Async            bool `json:"async"`
		CreateDefaults   bool `json:"create_defaults"`
		MeltdownRecovery bool `json:"meltdown_recovery"`
	}{
		Async:            true,
		CreateDefaults:   false,
		MeltdownRecovery: opts.MeltdownRecovey,
	}
	return s.asyncRequest(rest.MethodPost, "start", &params)
}

func (s *System) Shutdown() (*SystemDetails, error) {
	if e := s.logAction("Shutting down system"); e != nil {
		return nil, e
	}

	details, err := s.asyncRequest(rest.MethodPost, "shutdown", nil)
	if err != nil {
		return nil, err
	}

	return details, nil
}

func (s *System) Setup(answers map[string]interface{}, skipTests bool) (*SystemDetails, error) {
	params := struct {
		Async              bool                   `json:"async"`
		SkipTests          bool                   `json:"skip_tests"`
		ForceResignDevices bool                   `json:"force_resign_devices"`
		AutoStart          bool                   `json:"auto_start"`
		AnswerFile         map[string]interface{} `json:"answer_file,omitempty"`
	}{
		Async:              true,
		SkipTests:          skipTests,
		AutoStart:          true,
		ForceResignDevices: true,

		AnswerFile: answers,
	}
	return s.asyncRequest(rest.MethodPost, "setup", &params)
}

func (s *System) Deploy() (*SystemDetails, error) {
	return s.asyncRequest(rest.MethodPost, "deploy", nil)
}

func (s *System) AnswerFile() (*AnswerFile, error) {
	answerFileUri := "answer_file"
	uri := path.Join(systemsUri, fmt.Sprintf("%d", s.id), answerFileUri)

	var result AnswerFile
	err := s.session.Request(rest.MethodGet, uri, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type AnswerFile struct {
	Nodes []AnswerFileNode `json:"nodes"`
}

func (af *AnswerFile) ByHost(host string) (n *AnswerFileNode, err error) {
	for _, node := range af.Nodes {
		if node.Name == host {
			nd := node
			n = &nd
		}
	}
	if n == nil {
		return nil, errors.Errorf("Failed mathcing node by name: %s", host)
	}
	return
}

type AnswerFileService struct {
	ID struct {
		Type string `json:"type"`
	} `json:"id"`
	Cores    []int `json:"cores"`
	DeviceID struct {
		UUID string `json:"uuid"`
		Type string `json:"type"`
	} `json:"device_id,omitempty"`
}

type AnswerFileServices []AnswerFileService

func (a AnswerFileServices) ByDeviceID(id ...string) (dstores AnswerFileServices) {
	for _, dstore := range a {
		for _, devID := range id {
			if dstore.DeviceID.UUID == devID {
				dstores = append(dstores, dstore)
			}
		}
	}
	return dstores
}

type AnswerFileNode struct {
	Name     string              `json:"name"`
	Services []AnswerFileService `json:"services"`
}

const ServiceDStore = "DATASTORE_SERVICE_INSTANCE"

func (n AnswerFileNode) DStoreServices() (dstores AnswerFileServices) {
	for _, service := range n.Services {
		if service.ID.Type == ServiceDStore {
			dstores = append(dstores, service)
		}
	}
	return
}

type ReportListElement struct {
	Timestamp  time.Time `json:"timestamp"`
	ReportID   string    `json:"report_id"`
	ReportType string    `json:"report_type"`
	EnodeIds   []int     `json:"enode_ids"`
	Size       struct {
		Bytes int `json:"bytes"`
	} `json:"size"`
	IncludeControl bool   `json:"include_control"`
	Description    string `json:"description"`
}

func (s *System) ListReports() ([]ReportListElement, error) {
	var result []ReportListElement
	err := s.anyRequest(rest.MethodGet, "list_reports", nil, false, &result)

	return result, err

}

type CreatedReportDetails struct {
	ID           int       `json:"id"`
	UUID         string    `json:"uuid"`
	LastError    string    `json:"last_error,omitempty"`
	Priority     int       `json:"priority"`
	Attempts     int       `json:"attempts"`
	Queue        string    `json:"queue,omitempty"`
	Name         string    `json:"name"`
	CurrentStep  string    `json:"current_step"`
	StepProgress int       `json:"step_progress"`
	StepTotal    int       `json:"step_total"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Host         string    `json:"host,omitempty"`
	TaskType     string    `json:"task_type"`
	URL          string    `json:"url"`
}

type ReportType string

const (
	ReportTypeFull    ReportType = "full"
	ReportTypeMinimal ReportType = "minimal"
)

func getUUIDsMap(reportList []ReportListElement) map[string]bool {
	m := make(map[string]bool)
	for _, rep := range reportList {
		m[rep.ReportID] = true
	}

	return m
}

// This is a HACK until emanage (EL-3487) provides us with an api to get a newly created reports UUID
func findNewUUIDs(reportListBefore []ReportListElement, reportListAfter []ReportListElement) []string {
	var uuids []string
	beforeUUIDMap := getUUIDsMap(reportListBefore)
	for _, rep := range reportListAfter {
		if _, exists := beforeUUIDMap[rep.ReportID]; !exists {
			uuids = append(uuids, rep.ReportID)
		}
	}

	return uuids
}

func (s *System) CreateReportForNodes(reportType ReportType, ipList []string) ([]string, CreatedReportDetails, error) {
	var result CreatedReportDetails

	reportListBefore, e := s.ListReports()
	if e != nil {
		return nil, result, e
	}

	params := struct {
		Async      bool       `json:"async"`
		ReportType ReportType `json:"report_type"`
		IPList     []string   `json:"ip_list,omitempty"`
	}{
		Async:      true,
		ReportType: reportType,
		IPList:     ipList,
	}
	if err := s.anyRequest(rest.MethodPost, "create_report", &params, true, &result); err != nil {
		return nil, result, err
	}

	reportListAfter, err := s.ListReports()
	if err != nil {
		return nil, result, err
	}

	return findNewUUIDs(reportListBefore, reportListAfter), result, nil
}

func (s *System) CreateReportForAllNodes(reportType ReportType) ([]string, CreatedReportDetails, error) {
	return s.CreateReportForNodes(reportType, nil)
}

func (s *System) DeleteReport(uuid string, ipList []string) (*SystemDetails, error) {
	params := struct {
		ReportID string   `json:"report_id"`
		IPList   []string `json:"ip_list,omitempty"`
	}{
		ReportID: uuid,
		IPList:   ipList,
	}
	return s.request(rest.MethodPost, "delete_report", &params)
}

func (s *System) DeleteReportOnAllNodes(uuid string) (*SystemDetails, error) {
	return s.DeleteReport(uuid, nil)
}

type PreparedReportDetails struct {
	Path string `json:"path,omitempty"`
}

func (s *System) PrepareReport(uuid string, ipList []string) (*PreparedReportDetails, error) {
	var result PreparedReportDetails
	params := struct {
		ReportID   string     `json:"report_id"`
		ReportType ReportType `json:"report_type,omitempty"`
		IPList     []string   `json:"ip_list,omitempty"`
		PathOnly   bool       `json:"path_only,omitempty"`
	}{
		ReportID:   uuid,
		ReportType: ReportTypeFull,
		IPList:     ipList,
		// TODO: if PathOnly == false The report is sent as a type octet-stream which we do not support
		PathOnly: true,
	}
	if err := s.anyRequest(rest.MethodGet, "download_report", &params, false, &result); err != nil {
		return &result, err
	}

	return &result, nil
}

func (s *System) PrepareReportFromAllNodes(uuid string) (*PreparedReportDetails, error) {
	return s.PrepareReport(uuid, nil)
}

type Capacity struct {
	RawUsage           Bytes   `json:"raw_usage"`
	RawCapacity        Bytes   `json:"raw_capacity"`
	EffectiveUsage     Bytes   `json:"effective_usage"`
	EffectiveCapacity  Bytes   `json:"effective_capacity"`
	DataReductionRatio float64 `json:"data_reduction_ratio"`
	DedupRatio         float64 `json:"dedup_ratio"`
	CompressionRatio   float64 `json:"compression_ratio"`
	TopDataContainers  []struct {
		ID             int    `json:"id"`
		Name           string `json:"name"`
		UUID           string `json:"uuid"`
		UsedCapacity   Bytes  `json:"used_capacity"`
		NamespaceScope string `json:"namespace_scope"`
		DataType       string `json:"data_type"`
		Policy         struct {
			ID          int       `json:"id"`
			Name        string    `json:"name"`
			Dedup       int       `json:"dedup"`
			Compression int       `json:"compression"`
			Replication int       `json:"replication"`
			CreatedAt   time.Time `json:"created_at"`
			UpdatedAt   time.Time `json:"updated_at"`
			SoftQuota   Bytes     `json:"soft_quota"`
			HardQuota   Bytes     `json:"hard_quota"`
			IsTemplate  bool      `json:"is_template"`
			IsDefault   bool      `json:"is_default"`
		} `json:"policy"`
		PolicyID       int       `json:"policy_id"`
		Dedup          int       `json:"dedup"`
		Compression    int       `json:"compression"`
		SoftQuota      Bytes     `json:"soft_quota"`
		HardQuota      Bytes     `json:"hard_quota"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		ExportsCount   int       `json:"exports_count"`
		DirPermissions int       `json:"dir_permissions"`
		DirUID         int       `json:"dir_uid"`
		DirGid         int       `json:"dir_gid"`
	} `json:"top_data_containers"`
}

func (s *System) Capacity() (*Capacity, error) {
	var result Capacity
	err := s.anyRequest(rest.MethodGet, "capacity", nil, false, &result)
	return &result, err
}

func (s *System) GetVersion() (*Version, error) {
	var versionList VersionList
	versionUri := "version"
	uri := path.Join(systemsUri, fmt.Sprintf("%d", s.id), versionUri)

	err := s.session.Request(rest.MethodGet, uri, nil, &versionList)
	if err != nil {
		return nil, err
	}

	versionDetails, err := getECSVersion(&versionList)
	if err != nil {
		return nil, err
	}

	version, err := parseVersion(versionDetails.Version)
	if err != nil {
		return nil, err
	}

	return version, nil
}

func (s *System) GetLicense() (*License, error) {
	var license License
	licenseUri := "license"
	uri := path.Join(systemsUri, fmt.Sprintf("%d", s.id), licenseUri)

	err := s.session.Request(rest.MethodGet, uri, nil, &license)
	if err != nil {
		return nil, err
	}
	return &license, nil
}

func (s *System) UploadLicense(opts *LicenseOpts) (*License, error) {
	licenseUri := "upload_license"
	uri := path.Join(systemsUri, fmt.Sprintf("%d", s.id), licenseUri)

	var result License
	err := s.session.Request(rest.MethodPut, uri, opts, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// System  upgrade
type UpgradeOpts struct {
	Type                string `json:"type"`
	DegradedReplication bool   `json:"degraded_replication"`
	AdminPasswd         string `json:"admin_passwd"`
	SkipTest            bool   `json:"skip_version_test"`
	TarLink             string `json:"-"`
	EmanageHost         string `json:"-"`
}
type UpgradeOutput struct {
	ID           int         `json:"id"`
	UUID         string      `json:"uuid"`
	Handler      string      `json:"handler"`
	Priority     int         `json:"priority"`
	Attempts     int         `json:"attempts"`
	LastError    interface{} `json:"last_error"`
	Queue        interface{} `json:"queue"`
	Status       int         `json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	RunSync      bool        `json:"run_sync"`
	CurrentStep  interface{} `json:"current_step"`
	StepProgress interface{} `json:"step_progress"`
	StepTotal    interface{} `json:"step_total"`
	HostID       interface{} `json:"host_id"`
	TaskType     string      `json:"task_type"`
}

func (s *System) UpgradeStart(body UpgradeOpts) (*UpgradeOutput, error) {
	var upgrade UpgradeOutput
	upgUri := "upgrade_start"
	uri := path.Join(systemsUri, fmt.Sprintf("%d", s.id), upgUri)

	err := s.session.Request(rest.MethodPost, uri, body, &upgrade)
	if err != nil {
		return nil, err
	}
	return &upgrade, nil
}

type SystemUpgrade struct {
	ID           int         `json:"id"`
	UUID         string      `json:"uuid"`
	Handler      string      `json:"handler"`
	Priority     int         `json:"priority"`
	Attempts     int         `json:"attempts"`
	LastError    string      `json:"last_error"`
	Queue        interface{} `json:"queue"`
	Status       int         `json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	RunSync      bool        `json:"run_sync"`
	CurrentStep  string      `json:"current_step"`
	StepProgress int         `json:"step_progress"`
	StepTotal    int         `json:"step_total"`
	HostID       interface{} `json:"host_id"`
	TaskType     string      `json:"task_type"`
}
type SystemUpgradeState string
type SystemUpgradePhase string

const (
	SystemUpgradeStateIdle    SystemUpgradeState = "idle"
	SystemUpgradeStateRun     SystemUpgradeState = "running"
	SystemUpgradeStateFail    SystemUpgradeState = "failed"
	SystemUpgradeStatePaused  SystemUpgradeState = "paused"
	SystemUpgradeStatePausing SystemUpgradeState = "pausing"
	SystemUpgradePhaseNone    SystemUpgradePhase = "none"
	SystemUpgradePhaseEmanage SystemUpgradePhase = "emanage"
	SystemUpgradePhaseEnodes  SystemUpgradePhase = "enodes"
)

func (s *System) UpgradePause() error {
	var upgrade SystemUpgrade
	upgUri := "upgrade_pause"
	// body := &UpgradeOpts{}
	uri := path.Join(systemsUri, fmt.Sprintf("%d", s.id), upgUri)

	err := s.session.Request(rest.MethodPost, uri, nil, &upgrade)
	return err
}

func (s *System) UpgradeResume() error {
	var upgrade UpgradeOutput
	upgUri := "upgrade_resume"
	uri := path.Join(systemsUri, fmt.Sprintf("%d", s.id), upgUri)

	err := s.session.Request(rest.MethodPost, uri, nil, &upgrade)
	return err
}

func (s *System) SetupReplicationAgent() (*SystemDetails, error) {
	if e := s.logAction("Setup replication agent"); e != nil {
		return nil, e
	}
	params := struct {
		Async bool `json:"async"`
	}{
		Async: false,
	}
	return s.request(rest.MethodPost, "setup_replication_agent", &params)
}

func (s *System) DeployReplicationAgent() (*SystemDetails, error) {
	if e := s.logAction("Deploy replication agent"); e != nil {
		return nil, e
	}
	params := struct {
		Async bool `json:"async"`
	}{
		Async: false,
	}
	return s.request(rest.MethodPost, "deploy_replication_agent", &params)
}
