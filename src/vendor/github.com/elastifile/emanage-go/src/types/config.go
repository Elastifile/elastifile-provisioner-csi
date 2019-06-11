package types

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-errors/errors"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/koding/multiconfig"

	"helputils"
	logging_config "logging/config"
)

const (
	DataStorePathEnv     = "TESLA_DATASTORE_PATH"
	ContainerJobId       = "TESLA_JOB_ID"
	ConfElabSystemBucket = "conf-elab-system"
	ConfTestConfigs      = "test-configs"
	ConfProductConfig    = "product-config"
	FullConfig           = "full-config"
)

type StartupEnvVars struct {
	FilestoreHost string
}

const EnvPrefix = "ENVIRONMENT"

func (env *StartupEnvVars) GetFromEnvironment() error {
	envLoader := &multiconfig.EnvironmentLoader{
		Prefix:    EnvPrefix,
		CamelCase: true,
	}
	err := envLoader.Load(env)
	fmt.Printf("Environment: %+v\n", env)
	return err
}

type Config struct {
	Tesla        TeslaConfig           `doc:"tesla global configs"`
	Elab         Elab                  `doc:"elab global configs"`
	Tests        Tests                 `doc:"test configs"`
	Logging      logging_config.Config `doc:"logging configs"`
	//Tester       teslatestconfig.Config
	CloudConnect CloudConnect
	// primary system struct remains static so it would be accessed by '-o' flags (pointers cannot access!)
	System System `doc:"entire system data, tesla & elab"`
	// system slice is dynamic so the first would point to the primary, above
	Systems []*System `doc:"slice of system data pointers"`
}

func NewConfig() *Config {
	var conf Config
	//conf.Erun = &erun_config.Config{}
	tagLoader := &multiconfig.TagLoader{}

	if err := tagLoader.Load(&conf); err != nil {
		panic(err)
	}
	conf.Tesla.Jenkins.BuildNumber = os.Getenv("BUILD_NUMBER")
	conf.Tesla.Jenkins.BuildUrl = os.Getenv("BUILD_URL")

	return &conf
}

type Disks struct {
	Gigabytes int `json:"gigabytes"`
}

type CloudProvider struct {
	CloudProviderData
	LocalDiskSize      Disks `json:"local_disk_size,omitempty"`
	PersistentDiskSize Disks `json:"persistent_disk_size,omitempty"`
}

type CloudProviderUpdateOpts struct {
	CloudProviderData
	LocalDiskSize      int `json:"local_disk_size,omitempty"`
	PersistentDiskSize int `json:"persistent_disk_size,omitempty"`
}

type CloudProviderData struct {
	ID                     int    `json:"id,omitempty"`
	Project                string `json:"project,omitempty"`
	Image                  string `json:"image,omitempty"`
	ImageProject           string `json:"image_project,omitempty"`
	Zone                   string `json:"zone,omitempty"`
	Region                 string `json:"region,omitempty"`
	StorageType            string `json:"storage_type,omitempty"`
	LocalNumOfDisks        int    `json:"local_num_of_disks,omitempty"`
	LocalInstanceType      string `json:"local_instance_type,omitempty"`
	PersistentNumOfDisks   int    `json:"persistent_num_of_disks,omitempty"`
	PersistentInstanceType string `json:"persistent_instance_type,omitempty"`
	AwsAccessKeyID         string `json:"aws_access_key_id,omitempty"`
	AwsSecretAccessKey     string `json:"aws_secret_access_key,omitempty"`
	NameServer             string `json:"dns,omitempty"`
	DeploymentModel        string `json:"deployment model,omitempty" "default:cloud_amazon"`
	Instances              int    `json:"number of instances for Cloud,omitempty" default:"3"`
	CloudConfigurationID   int    `json:"cloud_configuration_id"`
	CloudConfigurations    CloudConfigurationCreateOpts
}

type CloudConfiguration struct {
	ID                int    `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	StorageType       string `json:"storage_type,omitempty"`
	NumOfDisks        int    `json:"num_of_disks,omitempty" default:"1"`
	DiskSize          Disks  `json:"disk_size,omitempty" default:"50"`
	MinNumOfInstances int    `json:"min_num_of_instances,omitempty" default:"3"`
	InstanceType      string `json:"instance_type,omitempty" default:"r4.2xlarge"`
	Cores             int    `json:"cores,omitempty"`
	Memory            int    `json:"memory,omitempty"`
}

type CloudConfigurationCreateOpts struct {
	CloudConfiguration
	Name              string `json:"name,omitempty"`
	StorageType       string `json:"storage_type,omitempty" default:"persistent"`
	NumOfDisks        int    `json:"num_of_disks,omitempty" default:"1"`
	DiskSize          int    `json:"disk_size,omitempty" default:"50"`
	InstanceType      string `json:"instance_type,omitempty" default:"r4.2xlarge"`
	Cores             int    `json:"cores,omitempty" default:"8"`
	Memory            int    `json:"memory,omitempty" default:"61"`
	MinNumOfInstances int    `json:"min_num_of_instances,omitempty" default:"3"`
}

type CloudConfigurationUpdateOpts struct {
	CloudConfigurationCreateOpts
	ID int `json:"id,omitempty"`
}

type TeslaConfig struct {
	Vcenter      VCenter
	Emanage      Emanage
	Docker       Docker
	Elfs         Elfs
	Jenkins      Jenkins
	Requirements Requirements
	Reporting    Reporting
	Monitoring   Monitoring
	Limits       Limits
	Outputs      Outputs
	ShellUser    string `doc:"Shell user for login to emanage/loader"`

	Auth struct {
		Uid int `doc:"Uid used to access files. If empty, use current user's uid"`
		Gid int `doc:"Gid used to access files. If empty, use current user's gid"`
	}

	// Known issues that should be ignored when running tests.
	Issues struct {
		Known           []string
		SkipWorkarounds []string
	}

	// Sanity
	Sanity struct {
		IsSanity       bool
		CleanOnly      bool
		SkipTests      bool
		Fail           bool
		DryRun         bool
		KeepContainers bool
		CleanOnFailure bool
		PollElfsHealth int64 `default:"10" doc:"Poll interval to check if elfs not in panic in seconds"`
	}

	// Development aids
	Devel Devel
}

type Elab struct {
	ServerUrl   string `default:"elab.il.elastifile.com"`
	BaseUri     string `default:"api/v1/system/cluster"`
	Timeout     int    `default:"180" doc:"timeout (seconds) for elab get system request + validating system info"`
	GetAttempts int    `default:"3"`
	Diagnostics struct {
		Lenient bool `default:"true" doc:"Warn if issues present instead of Error"`
	}
}

type Docker struct {
	Port        int             `default:"2375" doc:"Port on which the Docker daemon on each host listens. When not set, will connect to local Docker daemon via Unix domain socket"`
	PushPull    PushPullOptions ``
	NoUpdate    bool            `default:"true" doc:"Do not auto-update the tesla client executable"` // TODO: Take care when renaming due to hardcoded env var name (TESLA_DOCKER_NO_UPDATE)
	Local       bool                                                                                  // Run all docker commands using the local socket // TODO: Add a 'hidden' tag
	VolumesFrom []string                                                                              // Add volumes from these containers to newly created containers
}

type Outputs struct {
	TestRail  TestRail
	SqlServer SqlServer
}

type TestRail struct {
	Enabled   bool   `default:"false"`
	URL       string
	Username  string
	Password  string
	RunName   string
	PlanName  string `default:"Reported By Tesla"`
	SuiteName string `default:"Elastifile v2.0 Timeline"`
	ProjectID int    `default:"1"`
	Verbose   bool   `doc:"Add prints and traces"`
}

type SqlServer struct {
	Host        string `default:"35.205.81.217" doc:"SQL server address"`
	User        string `default:"dfr_user"`
	Password    string `default:"dfr"`
	DbName      string `default:"dfr"`
	SslMode     string `default:"" doc:"leaving unset, ssl is disabled."`
	ConnTimeout time.Duration
}

type Telegraf struct {
	Enabled         bool   `default:"false"`
	InfluxDBAddress string `default:"elfs-monitoring.lab.il.elastifile.com" doc:"InfluxDB server address"`
	Port            int    `default:"8086"`
	DBName          string `default:"telegraf"`
}

func (t Telegraf) DBServer() string {
	return fmt.Sprintf("%v:%v", t.InfluxDBAddress, t.Port)
}

type Jenkins struct {
	BuildNumber string `default:"" doc:"Jenkins BUILD_NUMBER, set automatically in Jenkins"`
	BuildUrl    string `default:"" doc:"Jenkins BUILD_URL, set automatically in Jenkins"`
}

type Elfs struct {
	Export       string `                   doc:"ELFS export to use, e.g. 'my_fs0/root'"`
	MountOptions string `default:"-o soft"  doc:"ELFS export mount options, e.g. '-o soft'"`
	ECS          Ecs
	EPA struct {
		Port int `default:"10015" doc:"EPA Port"`
	}
}

type Limits struct {
	FileCount int `default:"10" doc:"File count"`
	JobCount  int `default:"10" doc:"Job count"`

	TotalEntries          int `default:"200" doc:"Total amount of files and directories to be created"`
	FilesPerDirectory     int `default:"30"  doc:"Max amount of files per single directory"`
	FileToDirRatio        int `default:"5"   doc:"ratio of files to directories to create"`
	DirectoryNestingDepth int `default:"20"  doc:"max depth of nested directories"`
	MaxConcurrency        int `default:"3"   doc:"max NFS operations to perform concurrently"`
	WritesPerFile         int `default:"1"   doc:"number of writes to perform per created file"`
	ReadsPerFile          int `default:"1"   doc:"number of reads to perform per created file"`

	MaxFileSizeKb int `default:"128" doc:"Max file size in KiB"`
}

type Devel struct {
	Fake struct {
		Emanage     bool // fake emanage server
		SystemModel bool // fake system model
		Messaging   bool // fake messaging server
		ToolSim     bool // use toolsim instead of actual tools
		ToolTimeout int  // instruct toolsim to exit after n seconds
		// instruct toolsim to use this exit code if exiting
		// on timeout
		ToolExitCode int
		SystemRedeploy struct {
			Enabled       bool `default:"false"  doc:"Enables fake redeploy mode."`
			ShouldSucceed bool `default:"false"  doc:"In case 'True', simulate successful redeploy. In case 'False', simulate redeploy failure."`
		} `doc:"For the experts. Simulate redeploy failures/success."`
	}

	SkipResultsCollection bool // Skip collection of results at end of test suite (see INF-417)
	DumpHTTP              bool
	NatsTraces            bool `default:"true" doc:"Set Nats traces on"`
}

type Requirements struct {
	FreeLoaderSpace  int  `default:"4096"   doc:"Minimum available space (in MiB)"`
	ErrorWhenNoSpace bool `default:"false"  doc:"If false, will try to free disk space by removing Docker images"`
	ForceDisable     bool `default:"false"  doc:"If true, will not perform any requirements validation"`
}

type EMRI struct {
	Destination  string        `default:"/mnt/emri-reports" doc:"Directory where NFS share will be mounted (in EManage VM)"`
	NFSShare     string        `default:"file5.il.elastifile.com" doc:"DNS name of NFS share to store reports on"`
	NFSExport    string        `default:"/mnt/ssd-6T/emri-reports" doc:"Export to use when storing EMRI reports"`
	Kind         string        `default:"ci" doc:"Directory in which to store the reports"`
	Full         bool          `default:"true" doc:"Include ELFS traces in report"`
	MountTimeout time.Duration `default:"5m" doc:"retry mount till this period have passed"`
}

type Reporting struct {
	Enabled   bool   `default:"true" doc:"Master switch for reporting."`
	EMRI      bool   `default:"false" doc:"Collect EMRI reports on test failure"`
	ELFS      bool   `default:"true" doc:"Collect ELFS service status on test failure"`
	ECS       bool   `default:"true" doc:"Collect ECS logs on test failure"`
	Emanage   bool   `default:"true" doc:"Collect EManage logs on test failure"`
	Loaders   bool   `default:"true" doc:"Collect loaders status on test failure"`
	ECSCLI    bool   `default:"true" doc:"Collect ecs-cli output on test failure"`
	EMRIOpts  EMRI
	MountKind string `default:"privileged" doc:"Recognized values are privileged, nfs and none"`
}

func (rep Reporting) String() string {
	var parts []string
	if rep.EMRI {
		parts = append(parts, "EMRI")
	}
	if rep.ELFS {
		parts = append(parts, "ELFS")
	}
	if rep.ECS {
		parts = append(parts, "ECS")
	}
	if rep.Emanage {
		parts = append(parts, "Emanage")
	}
	if rep.Loaders {
		parts = append(parts, "Loaders")
	}
	if rep.ECSCLI {
		parts = append(parts, "ECSCLI")
	}
	return "[" + strings.Join(parts, " ") + "]"
}

type ReportingOpts struct {
	Reporting       Reporting
	EmanageServer   string
	EmanageUser     string
	EmanagePassword string
	VcenterHost     string
	FileStoreHost   string
	Loaders         []Host
}

type Monitoring struct {
	Enabled      bool `default:"false"`
	Tcpdump      bool `default:"true"`
	IpNeighbours bool `default:"true"`
	Telegraf     Telegraf
}

type Tests struct {
	ProductConfigFile string                                                     // `default:"tesla/config/defultProductConfiguration.yaml"`
	TestConfigsFile   string `default:"../config/defultTestsConfiguration.yaml"` // relative to this bin path
	ColdUpgrade       ColdUpgrade

	ExpectedFailures     int  `doc:"number of expected failures"`
	MaxFailures          int  `doc:"number of failures before suite abort"`
	FailOnToolValidation bool `doc:"Fail test if validation failed"`
	OnFailure struct {
		MustReport bool `doc:"Abort cycle if report generation fails"`
		// Redeploy must remain disabled by default, since Redeploy is very heavy process -
		// enablement should be by external regression manager, only!
		SystemRedeploy                bool          `default:"false"  doc:"Redeploy the system to same RPMs, after test failure."`
		SystemRedeployRetries         int           `default:"1"      doc:"system redeploy retries count"`
		SystemRedeployRetriesInterval time.Duration `default:"30s"    doc:"duration to wait between redeployments"`
	}
}

func (t *Tests) TestConfigsPath() string {
	if strings.HasPrefix(t.TestConfigsFile, ".") {
		return filepath.Join(helputils.StartupFolder(), t.TestConfigsFile)
	}
	return t.TestConfigsFile
}

type ColdUpgrade struct {
	EmanageMachineIp string // IP
	UpgradeTarLink   string // url
}

type LocalInstall struct {
	ELFS            string `doc:"path to ELFS rpm to install (relative to tesla client)"`
	ELFSTraceViewer string `doc:"path to ELFS traceviewer rpm to install (relative to tesla client)"`
	ELFSSource      string `doc:"path to ELFS source rpm to install (relative to tesla client)"`
	ELFSTools       string `doc:"path to ELFS tools rpm to install (relative to tesla client)"`
	ELFSTop         string `doc:"path to ELFS top rpm to install (relative to tesla client)"`
	ECS             string `doc:"path to ECS rpm to install (relative to tesla client)"`
	EPA             string `doc:"path to EPA rpm to install (relative to tesla client)"`
	EVP             string `doc:"path to EVP rpm to install (relative to tesla client)"`
	EManage         string `doc:"path to EManage rpm to install (relative to tesla client)"`
	EManageAssets   string `doc:"path to EManage assets rpm to install (relative to tesla client)"`
	ELFSCLI         string `doc:"path to ELFS cli rpm to install (relative to tesla client)"`
	LLVM            string `doc:"path to LLVM libraries rpm to install (relative to tesla client)"`
	EAdmin          string `doc:"path to EAdmin rpm to install (relative to tesla client)"`
	NFSUtils        string `doc:"path to nfs-utils rpm to install (relative to tesla client)"`
	Esync           string `doc:"path to elfs-esync rpm to install (relative to tesla client)"`
	RsyncDR         string `doc:"path to rsync-dr rpm to install (relative to tesla client)"`
}

type Deploy struct {
	Type                int        `doc:"The deployment type to use"`
	ReplicationLevel    int        `default:"2" doc:"The system replication level for all policies and data containers"`
	ExitIfPanic         bool       `doc:"Check and exit if the system is in panic state before the deployment"`
	SkipEmanageTests    bool       `default:"true" doc:"Skip eManage tests upon deploy"`
	ForceInstall        bool       `doc:"Force a clean install regardless the pkg version"`
	EnableClusterKiller bool       `default:"false" doc:"Whether allowing the option to enable system cluster killer"`
	SaveElfsDevParams   bool       `default:"false" doc:"Whether to save elfs dev params"`
	ReplicationAgent    bool       `default:"false" doc:"Whether to deploy a Replication Agent on system deploy."`
	ConfigFile          string     `doc:"JSON config file path, for func setups (empty for VC setups)"`
	Data                DeployData `doc:"if config file is empty then data inside this struct should be provided."`
	TimeZone            string     `default:"Asia/Jerusalem" doc:"Time zone for loaders and vHeads"`
	NTPServers          string     `default:"time1.google.com,time2.google.com,time3.google.com" doc:"NTP servers to synchronize loaders"`
	YumRepository       string     `default:"http://centrepo.il.elastifile.com/el-centos7.1.repo" doc:"Yum repository url"`
	Ova                 OvaDeployment
	LocalInstall        LocalInstall
	LinuxTrace          bool
	CloudProviders      CloudProviderUpdateOpts
	CloudConfigurations CloudConfigurationCreateOpts
}

type System struct {
	Elab         ElabSystem    `json:"elab"         doc:"elab system data"`
	Builds       Builds        `json:"builds"       doc:"System component build versions"`
	Deploy       Deploy        `json:"deploy"       doc:"deployment configuration`
	LicenseFile  string        `json:"license_file" doc:"System license file path"`
	FormatMethod string        `default:"Reset"     doc:"System format method (Redeploy, Reset)"`
	Frontend     string        `                    doc:"ELFS frontend address, e.g. '192.168.0.1'"`
	Site         ElabSite
	Provision    ProvisionData `doc:"Parameters for cloud provision"`
	Setup struct {
		CurrentLoader string `doc:"The host that we're currently running on (must be one of the loaders)"`
		Filestore     string `doc:"host for minio"`
	}
}

type CloudConnect struct {
	Server       string `doc:"ElCC  server address, e.g. '10.11.x.x'"`
	User         string `default:"admin" doc:"username to login to ccweb "`
	Password     string `default:"changeme" doc:"password to login to ccweb "`
	AccessKey    string `default:"" doc:"AccessKey to login to Cloud "`
	SecretKey    string `default:"" doc:"SecretKey to login to Cloud "`
	CloudService string `default:"aws" doc:"object_store"`
	TarFileLink  string `doc:"link to Cloud Connect tar file"`
}

func (sys System) IsFunctional() bool {
	return sys.Deploy.ConfigFile != ""
}

const (
	RPMs      int = iota
	Binaries
	OVA
	ReuseRPMs
	LocalRPMs
	Cloud

	LastDeployType = Cloud
)

var deploymentTypeNames = []string{
	"RPMs",
	"Binaries",
	"OVA",
	"ReuseRPMs",
	"LocalRPMs",
	"Cloud",
}

func DeploymentTypeValid(dt int) bool {
	return 0 <= dt && dt <= LastDeployType
}

func DeploymentTypeString(dt int) string {
	if DeploymentTypeValid(dt) {
		return deploymentTypeNames[dt]
	} else {
		return fmt.Sprintf("invalid deployment type code: %d", dt)
	}
}

type VCDestination struct {
	EsxHost   string
	Datastore string
}

type OvaDeployment struct {
	InstallerPath                string `doc:"Installer script path"`
	EmanageManagementNetworkName string `default:"VM Network" doc:"eManage external/management network. e.g. 'VM Network'"`
	NetworkAutoDetect            bool   `default:"false"`
	DryRun                       bool   `default:"false" doc:"Run deployment as dry-run"`
	AbortAfterStage              int    `default:"0" doc:"Abourt deployment stage"`
}

type Emanage struct {
	Username   string `default:"admin"         doc:"Emanage username"`
	Password   string `default:"changeme"      doc:"Emanage password"`
	SystemId   int    `default:"1"             doc:"Emanage system ID (usually 1)"`
	VIPNetmask string `default:"255.255.255.0" doc:"Emanage virtaul server netmask"`
	Address    string
}

type VCenter struct {
	HostDomain string `default:"lab.il.elastifile.com"`
	Username   string `default:"root"`
	Password   string `default:"vmware"`
	VmRevive   bool   `default:"false"`
	IsFunc     bool   `default:"false"`
}

func (vc *VCenter) URL(cid string) *url.URL {
	// Format: "https://root:vmware@vc7.lab.il.elastifile.com/sdk"
	return &url.URL{
		Scheme: "https",
		User:   url.UserPassword(vc.Username, vc.Password),
		Host:   vc.Host(cid),
		Path:   "/sdk",
	}
}

func (vc *VCenter) Host(cid string) string {
	if vc.IsFunc {
		cid = "7a"
	}
	return fmt.Sprintf("vc%s.%s", cid, vc.HostDomain)
}

type InterfaceNet struct {
	Address string
	Subnet  int
	Name    string
}

type EnodeUpgradeState string

const (
	EnodeUpgradeStateIdle     EnodeUpgradeState = "status_idle"
	EnodeUpgradeStateUpgrade  EnodeUpgradeState = "status_upgrading"
	EnodeUpgradeStateFail     EnodeUpgradeState = "status_upgrade_failed"
	EnodeUpgradeStateStop     EnodeUpgradeState = "stopped"
	EnodeUpgradeStateStopping EnodeUpgradeState = "stopping"
)

type VHeadRole string

const (
	Frontend  VHeadRole = "frontent"  // HCI Frontend
	Converged VHeadRole = "converged" // HCI Frontend+Backend
	DataStore VHeadRole = "datastore" // BM Frontent+Backend
)

func (vhr VHeadRole) String() string {
	return string(vhr)
}

type DeployData struct {
	DataNet         InterfaceNet
	DataNet2        InterfaceNet
	NfsNet          InterfaceNet
	DisksPrefix     []string
	Enodes          []DeployEnode
	ExternalNetName string `default:"VM Network"   doc:"eManage external/management network. e.g. 'VM Network'"`
	ExternalUseDhcp bool
	Datastore       string
	VLANID          int
	VLANID2         int
	BroadcastNic    string
	BroadcastNic2   string
	LoaderNFSNic    string
}

type DeployEnode struct {
	Name     string
	Type     string
	Role     string
	DataMac  string
	DataMac2 string
	NfsMac   string
}

func (sys *System) DefaultVip(net *Network) string {
	return sys.Vip(net, 1)
}

func (sys *System) Vip(net *Network, n int) string {
	if sys.Elab.IsCloud() {
		return sys.Frontend
	} else {
		if net.NetworkId == "" {
			return ""
		} else if strings.HasSuffix(net.NetworkId, ".0") {
			return strings.Join(strings.Split(net.NetworkId, ".")[:3], ".") + "." + strconv.Itoa(n)
		}
		return net.NetworkId
	}
}

func (sys *System) IsPhysicalEnode(i int) bool {
	if 0 <= i && i < len(sys.Deploy.Data.Enodes) {
		return sys.Deploy.Data.Enodes[i].Type == "physical"
	}
	return false
}

func (sys *System) EnodeRole(i int, def VHeadRole) string {
	enodes := &sys.Deploy.Data.Enodes
	if 0 <= i && i < len(*enodes) && (*enodes)[i].Role != "" {
		return (*enodes)[i].Role
	}
	return string(def)
}

func (sys *System) EnodeRoles() (roles []string) {
	for _, enode := range sys.Deploy.Data.Enodes {
		roles = append(roles, enode.Role)
	}
	return roles
}

// CurrentLoaderInternal returns the current loader host on internal network
func (sys *System) CurrentLoaderInternal() (loaderInternal string) {
	loaderInternal = sys.Setup.CurrentLoader
	for i, loaderExternal := range sys.Elab.LoaderIpsExternal() {
		if loaderInternal == loaderExternal {
			loaderInternal = sys.Elab.LoaderIpsInternal()[i]
		}
	}
	return
}

// masterLoaderNode returns the master loader node
func (sys *System) masterLoaderNode() *LoaderNode {
	loaders := sys.Elab.UsableLoaders()
	if len(loaders) > 0 {
		return &loaders[0]
	}
	return nil
}

// MasterLoader returns the master loader host
func (sys *System) MasterLoader() Host {
	loaderNode := sys.masterLoaderNode()
	if loaderNode == nil {
		return ""
	}
	return Host(loaderNode.IpAddress)
}

// masterLoaderByNetworkType returns the master loader host on the specified network
func (sys *System) masterLoaderByNetworkType(networkType int) Host {
	if sys.Elab.IsCloud() {
		loaderNode := sys.masterLoaderNode()
		if loaderNode != nil {
			switch networkType {
			case InternalNetworkId:
				return Host(loaderNode.InternalAddr())
			case ExternalNetworkId:
				return Host(loaderNode.ExternalAddr())
			default:
				panic(fmt.Sprintf("Unsupported network type id: %v", networkType))
			}
		} else {
			return ""
		}
	} else {
		return sys.MasterLoader()
	}
}

// MasterLoaderInternal returns the master loader host on internal network
func (sys *System) MasterLoaderInternal() Host {
	return sys.masterLoaderByNetworkType(InternalNetworkId)
}

// MasterLoaderExternal returns the master loader host on external network
func (sys *System) MasterLoaderExternal() Host {
	return sys.masterLoaderByNetworkType(ExternalNetworkId)
}

// TeslaManager returns the Tesla Manager's IP address
func (sys *System) TeslaManager() Host {
	return sys.MasterLoader()
}

// TeslaManagerExternal returns the Tesla Manager's external IP address
func (sys *System) TeslaManagerExternal() Host {
	return sys.masterLoaderByNetworkType(ExternalNetworkId)
}

// TeslaManagerInternal returns the Tesla Manager's internal IP address
func (sys *System) TeslaManagerInternal() Host {
	return sys.masterLoaderByNetworkType(InternalNetworkId)
}

func (sys *System) MessagingHost() Host {
	master := sys.TeslaManager()
	if master != "" {
		return master
	}
	return "localhost"
}

func (sys *System) messagingHostByNetworkType(networkType int) Host {
	master := Host("localhost")
	switch networkType {
	case InternalNetworkId:
		master = sys.TeslaManagerInternal()
	case ExternalNetworkId:
		master = sys.TeslaManagerExternal()
	default:
		panic(fmt.Sprintf("Unsupported network type id: %v", networkType))
	}

	return master
}

func (sys *System) MessagingHostInternal() Host {
	return sys.messagingHostByNetworkType(InternalNetworkId)
}

func (sys *System) MessagingHostExternal() Host {
	return sys.messagingHostByNetworkType(ExternalNetworkId)
}

func (sys *System) filestoreHostByNetworkType(networkType int) Host {
	var address string

	switch networkType {
	case InternalNetworkId:
		address = string(sys.TeslaManagerInternal())
	case ExternalNetworkId:
		address = string(sys.TeslaManagerExternal())
	default:
		panic(fmt.Sprintf("Unsupported network type id: %v", networkType))
	}

	if address != "" {
		sys.Setup.Filestore = address
	}

	if sys.Setup.Filestore == "" {
		panic(fmt.Sprintf("Filestore on system %v is not initialized. sys: %+v", sys.Elab.Data.Id, sys))
	}

	return Host(sys.Setup.Filestore)
}

// FilestoreHostInternal returns filestore host on internal network
func (sys *System) FilestoreHostInternal() Host {
	return sys.filestoreHostByNetworkType(InternalNetworkId)
}

// FilestoreHostExternal returns filestore host on external network
func (sys *System) FilestoreHostExternal() Host {
	return sys.filestoreHostByNetworkType(ExternalNetworkId)
}

// TODO: this function (and its usages) can be removed once eLab does this
// TODO: Consider removing this function altogether, presuming internal/external approach makes it redundant
func (sys *System) SwitchToInternalAddresses() {
	// Update Loaders
	for i := range sys.Elab.Data.Loaders {
		sys.Elab.Data.Loaders[i].IpAddress = sys.Elab.Data.Loaders[i].InternalAddr()
	}

	// Update eManages
	for i := range sys.Elab.Data.Emanage {
		sys.Elab.Data.Emanage[i].IpAddress = sys.Elab.Data.Emanage[i].InternalAddr()
	}

	// Update filestore
	sys.Setup.Filestore = string(sys.FilestoreHostInternal())
}

type Builds struct {
	Ident string `default:"lastStableBuild"`
}

type Usage struct {
	Percent float64 `json:"percent"`
}

type Capacity struct {
	Bytes int `json:"bytes"`
}

type ToolName string

const (
	Erun      ToolName = "erun"
	Sfs2008            = "sfs2008"
	Sfs2014            = "sfs2014"
	Teslatest          = "teslatest"
	Migration          = "migration"
	Cthon              = "cthon"
	Vdbench            = "vdbench"
	FsTool             = "fstool"
	Fio                = "fio"
)

const (
	SystemDeployCommand      = "system_deploy"
	ReplicationDeployCommand = "replication_deploy"
	SystemProvisionCommand   = "system_provision"
	SystemDestroyCommand     = "system_destroy"
)

type Status string

const (
	StatusEnodeFailed     Status = "enode_failed"
	StatusEnodeInit       Status = "enode_init"
	StatusEnodeConfigured Status = "enode_configured"
	StatusEnodeRunning    Status = "enode_running"
	StatusEnodeStopped    Status = "enode_stopped"
	StatusEnodeRemoving   Status = "enode_removing"
	StatusEnodeRemoved    Status = "enode_removed"
)

type SetupStatus string

const (
	SetupStatusUninitialized  SetupStatus = "uninitialized"
	SetupStatusDeployed       SetupStatus = "deployed"
	SetupStatusConfigured     SetupStatus = "configured"
	SetupStatusActive         SetupStatus = "active"
	SetupStatusPendingRemoval SetupStatus = "pending_removal"
)

type StatusName string

const (
	StatusNamePendingRemoval StatusName = "pending_removal"
	StatusNameInit           StatusName = "init"
	StatusNameUnreachable    StatusName = "unreachable"
	StatusNameNeedRecovery   StatusName = "need_recovery"
	StatusNameDegraded       StatusName = "degraded"
	StatusNameActive         StatusName = "active"
	StatusNameInitializing   StatusName = "initializing"
	StatusNameFenced         StatusName = "fenced"
)

func SupportedTools() []ToolName {
	return []ToolName{
		Erun,
		Sfs2008,
		Sfs2014,
		Teslatest,
		Migration,
		Cthon,
		Vdbench,
		FsTool,
		Fio,
	}
}

// Get vheads as Hosts
func (conf *Config) VHeads() (vheads []Host) {
	for _, vhd := range conf.System.Elab.Data.Vheads {
		vheads = append(vheads, Host(vhd.IpAddress))
	}
	return vheads
}

// Get loaders as Hosts
func (conf Config) Loaders() (loaders []Host) {
	for _, ip := range conf.LoadersStr() {
		loaders = append(loaders, Host(ip))
	}
	return loaders
}

// LoadersExternal returns a slice of loader IP addresses on the specified network.
// On setups with no external/internal networks, e.g. on-prem, default IP addresses are returned.
func (conf Config) loadersByNetworkType(networkType int) (loaders []Host) {
	if !conf.System.Elab.IsCloud() {
		loaders = conf.Loaders()
	} else {
		for _, loader := range conf.System.Elab.UsableLoaders() {
			loaders = append(loaders, Host(loader.Networks[networkType].IpAddress))
		}
	}
	return loaders
}

func (conf Config) LoadersInternal() (loaders []Host) {
	return conf.loadersByNetworkType(InternalNetworkId)
}

func (conf Config) LoadersExternal() (loaders []Host) {
	return conf.loadersByNetworkType(ExternalNetworkId)
}

// Get loaders of all systems as Hosts
func (conf Config) AllLoaders() (loaders []Host) {
	for _, sys := range conf.Systems {
		for _, ip := range sys.Elab.LoaderIps() {
			loaders = append(loaders, Host(ip))
		}
	}
	return loaders
}

// Get loaders of all systems as Hosts
func (conf Config) AllLoadersInternal() (loaders []Host) {
	for _, sys := range conf.Systems {
		for _, ip := range sys.Elab.LoaderIpsInternal() {
			loaders = append(loaders, Host(ip))
		}
	}
	conf.Logging.PrintWarn("AllLoadersInternal - exit", "loaders", loaders)
	return loaders
}

// For parsing text templates (e.g. sfs2008 config files)
func (conf Config) LoadersStr() (loaders []string) {
	return conf.System.Elab.LoaderIps()
}

func (conf *Config) SetLoaders(ipAddrs ...string) {
	loaders := make([]LoaderNode, 0)
	// filter existing loaders on conf
	for _, loader := range conf.System.Elab.Data.Loaders {
		for _, ip := range ipAddrs {
			if loader.IpAddress == ip {
				loaders = append(loaders, loader)
			}
		}
	}
	// TODO:
	// add the following validation before exit, len(loaders) should equal len(ipAddrs) and
	// return error if not (to ensure we found the exact amount of requested items).
	conf.System.Elab.Data.Loaders = loaders
}

func (conf *Config) MasterLoader() Host {
	return conf.System.MasterLoader()
}

func (conf *Config) MasterLoaderInternal() Host {
	return conf.System.MasterLoaderInternal()
}

func (conf *Config) MasterLoaderExternal() Host {
	return conf.System.MasterLoaderExternal()
}

// Get Disks
func (conf *Config) DisksPrefix() (disks []string) {
	return conf.System.Deploy.Data.DisksPrefix
}

// Get Enodes
func (conf *Config) Enodes() (enodes []string) {
	for _, en := range conf.System.Deploy.Data.Enodes {
		enodes = append(enodes, en.Name)
	}
	return enodes
}

// Get message bus host
func (conf *Config) MessagingHost() Host {
	return conf.System.MessagingHost()
}

// Get message bus host on internal network
func (conf *Config) MessagingHostInternal() Host {
	return conf.System.MessagingHostInternal()
}

// Get message bus host on external network
func (conf *Config) MessagingHostExternal() Host {
	return conf.System.MessagingHostExternal()
}

// FilestoreHost returns filestore host
func (conf *Config) FilestoreHost() Host {
	if conf.System.Setup.Filestore == "" {
		conf.System.Setup.Filestore = string(conf.System.TeslaManager())
		if conf.System.Setup.Filestore == "" {
			panic(fmt.Sprintf("conf.FilestoreHost() is not initialized - %+v", conf))
		}
	}
	return Host(conf.System.Setup.Filestore)
}

//// FilestoreHostInternal returns filestore host on internal network
//func (conf *Config) FilestoreHostInternal() Host {
//	// TODO: Find the root cause of muti-system code being broken - this  is just a hack to make tests run in the mean time
//	if len(conf.Systems) > 0 {
//		return conf.Systems[0].filestoreHostByNetworkType(InternalNetworkId)
//	} else {
//		helputils.PrintErrorWithStack(errors.Errorf("Empty conf.Systems: %+v", conf))
//		return conf.System.filestoreHostByNetworkType(InternalNetworkId)
//	}
//}

// FilestoreHostInternal returns filestore host on internal network
func (conf *Config) FilestoreHostInternal() Host {
	return conf.System.filestoreHostByNetworkType(InternalNetworkId)
}

// FilestoreHostExternal returns filestore host on external network
func (conf *Config) FilestoreHostExternal() Host {
	return conf.System.filestoreHostByNetworkType(ExternalNetworkId)
}

// MatchInterfacesOnSameHost returns true if the IP addresses are identical
// or both belong to the same host according to eLab
// Note: this function doesn't account for on-prem data/client networks
func (conf *Config) MatchInterfacesOnSameHost(addrA Host, addrB Host) bool {
	var allHostNets [][]NodeNetwork
	if addrA == addrB {
		return true
	}

	if conf.System.Elab.IsCloud() {
		for _, sys := range conf.Systems {
			for _, node := range sys.Elab.Data.Loaders {
				allHostNets = append(allHostNets, node.Networks)
			}
			for _, node := range sys.Elab.Data.Emanage {
				allHostNets = append(allHostNets, node.Networks)
			}
			for _, node := range sys.Elab.Data.Vheads {
				allHostNets = append(allHostNets, node.Networks)
			}
		}

		for _, singleHostNets := range allHostNets {
			matchCount := 0
			for _, singleNet := range singleHostNets {
				if singleNet.IpAddress == string(addrA) || singleNet.IpAddress == string(addrB) {
					matchCount++
				}
			}
			if matchCount == 2 { // Both addresses are accounted for
				return true
			}
		}
	}
	return false
}

// Get ELFS frontend as Host
func (conf *Config) Frontend() Host {
	return Host(conf.System.Frontend)
}

// Get emanage server URL
func (conf *Config) EmanageURL() *url.URL {
	baseURL := &url.URL{
		Scheme: "http",
		Host:   conf.EmanageServer(),
	}
	return baseURL
}

// Get emanage server URL
func (conf *Config) CloudConnectURL() *url.URL {
	if conf.CloudConnect.Server == "" {
		return nil
	}
	return &url.URL{
		Scheme: "http",
		Host:   conf.CloudConnect.Server,
	}
}

func (conf *Config) EmanageServer() string {
	return conf.System.Elab.EmanageIp()
}

func (conf *Config) AddSystem() int {
	conf.Systems = append(conf.Systems, &System{})
	last := len(conf.Systems) - 1

	*conf.Systems[last] = conf.System

	lastElab := &conf.Systems[last].Elab
	conf.Logging.PrintDebug("Added system",
		"index", last,
		"id", lastElab.Data.Id,
		"name", lastElab.Data.Name,
		"LoaderIps", lastElab.LoaderIps(),
		"Emanage", lastElab.Data.Emanage,
	)

	return len(conf.Systems)
}

func (conf *Config) MakeFirstSystemReference() {
	if len(conf.Systems) > 0 {
		conf.System = *conf.Systems[0]
		conf.Systems[0] = &conf.System
	}
}

func (conf *Config) SystemIndexByLoader(host Host) int {
	for i, sys := range conf.Systems {
		for _, loader := range sys.Elab.LoaderIps() {
			if loader == string(host) {
				return i
			}
		}
	}
	return -1
}

func (conf *Config) SystemById(id string) *System {
	for _, sys := range conf.Systems {
		if sys.Elab.Data.Id == id {
			return sys
		}
		conf.Logging.PrintDebug("did not match", "sys.Elab.Data.Id", sys.Elab.Data.Id, "id", id)
	}
	return nil
}

func (conf *Config) SystemAttrStrings(attr func(sys *System) interface{}) []string {
	result := make([]string, len(conf.Systems))
	for i, sys := range conf.Systems {
		result[i] = fmt.Sprintf("%v", attr(sys))
	}
	return result
}

func (conf *Config) SystemNames() []string {
	return conf.SystemAttrStrings(
		func(sys *System) interface{} { return sys.Elab.Data.Name },
	)
}

func (conf *Config) SystemIds() []string {
	return conf.SystemAttrStrings(
		func(sys *System) interface{} { return sys.Elab.Data.Id },
	)
}

func (conf *Config) SystemDeployTypes() (result []string, err error) {
	result = conf.SystemAttrStrings(
		func(sys *System) interface{} {
			if !DeploymentTypeValid(sys.Deploy.Type) {
				err = multierror.Append(err, errors.New(DeploymentTypeString(sys.Deploy.Type)))
			}
			return DeploymentTypeString(sys.Deploy.Type)
		},
	)
	return
}

func FilePathInDatastore(filename string) (string, error) {
	datastore, ok := os.LookupEnv(DataStorePathEnv)
	if !ok {
		return filename, errors.Errorf("missing env var %s", DataStorePathEnv)
	}
	return filepath.Join(datastore, filename), nil
}

func FilePathInJenkins(filename string, currentLoader string) string {
	jobID, ok := os.LookupEnv(ContainerJobId)
	if !ok {
		panic(errors.Errorf("Missing teslatest's container job id (%v)", ContainerJobId))
	}

	return filepath.Join("results", jobID, "teslatest", currentLoader, filename)
}

type Ecs struct {
	Server string `default:"0.0.0.0:10016" doc:"ECS server address and port, e.g. 'func11-cm:10016'"`
}

func (ecs *Ecs) SetServer(host string) {
	ecs.Server = strings.Replace(ecs.Server, strings.Split(ecs.Server, ":")[0], host, 1)
}

func (ecs *Ecs) Host() string {
	return strings.Split(ecs.Server, ":")[0]
}
