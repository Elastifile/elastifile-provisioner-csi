package types

type ElabSites struct {
	Data []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
}

type ElabSite struct {
	Data struct {
		Cloud      ElabSiteCloud      `json:"cloud"`
		CreatedAt  int                `json:"created_at"`
		Docker     ElabSiteDocker     `json:"docker"`
		Domain     string             `json:"domain"`
		Eloader    ElabSiteEloader    `json:"eloader"`
		Emri       ElabSiteEmri       `json:"emri"`
		Host       ElabSiteHost       `json:"host"`
		ID         string             `json:"id"`
		Name       string             `json:"name"`
		Provider   int                `json:"provider"`
		Repository ElabSiteRepository `json:"repository"`
		Role       int                `json:"role"`
		UpdatedAt  int                `json:"updated_at"`
		Vcenter    ElabSiteVcenter    `json:"vcenter"`
	} `json:"data"`
	Meta struct {
		APIVersion  int     `json:"api_version"`
		ElabVersion string  `json:"elab_version"`
		Error       bool    `json:"error"`
		ExecTime    float64 `json:"exec_time"`
		Group       string  `json:"group"`
		Instance    string  `json:"instance"`
		Resource    string  `json:"resource"`
		Status      int     `json:"status"`
		Timestamp   int     `json:"timestamp"`
	} `json:"meta"`
}

type ElabSiteCloud struct {
	Image struct {
		Project string `json:"project"`
	} `json:"image"`
	Password string `json:"password"`
	Project  struct {
		Prefix string `json:"prefix"`
	} `json:"project"`
	Region   string `json:"region"`
	Security struct {
		Profile string `json:"profile"`
	} `json:"security"`
	User string `json:"user"`
	Zone string `json:"zone"`
}

type ElabSiteDocker struct {
	Auth     bool   `json:"auth"`
	Password string `json:"password"`
	Registry string `json:"registry"`
	Secured  bool   `json:"secured"`
	User     string `json:"user"`
}

type ElabSiteEmri struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	Path     string `json:"path"`
	Type     string `json:"type"`
	User     string `json:"user"`
}

type ElabSiteEloader struct {
	Password string `json:"password"`
	Prikey   string `json:"prikey"`
	Pubkey   string `json:"pubkey"`
	User     string `json:"user"`
	Version  string `json:"version"`
}

type ElabSiteHost struct {
	Password string `json:"password"`
	User     string `json:"user"`
}

type ElabSiteRepository struct {
	Git string `json:"git"`
	Rpm string `json:"rpm"`
}

type ElabSiteVcenter struct {
	Password string `json:"password"`
	User     string `json:"user"`
}
