package emanage

import (
	"fmt"
	"time"

	"eurl"
	"optional"
	"rest"
)

const policiesUri = "api/policies"

type policies struct {
	conn *rest.Session
}

type Policy struct {
	Id          int              `json:"id"`
	Name        string           `json:"name"`
	Dedup       DedupLevel       `json:"dedup"`
	Compression CompressionLevel `json:"compression"`
	Replication ReplicationLevel `json:"replication"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Url         eurl.URL         `json:"url"`
	IsTemplate  bool        `json:"is_template"`
	IsDefault   bool        `json:"is_default"`
	TenantID    int         `json:"tenant_id"`
}

func (p *policies) GetAll(opt *GetAllOpts) (result []Policy, err error) {
	if opt == nil {
		opt = &GetAllOpts{}
	}
	err = p.conn.Request(rest.MethodGet, policiesUri, opt, &result)
	return
}

func (p *policies) GetFull(policyId int) (Policy, error) {
	uri := fmt.Sprintf("%s/%d", policiesUri, policyId)
	var result Policy
	err := p.conn.Request(rest.MethodGet, uri, nil, &result)
	return result, err
}

type PolicyCreateOpts struct {
	Dedup       *DedupLevel       `json:"dedup,omitempty"`
	Compression *CompressionLevel `json:"compression,omitempty"`
	Replication *ReplicationLevel `json:"replication,omitempty"`
	TenantId    optional.Int      `json:"tenant_id,omitempty"`
}

func (p *policies) Create(name string, opt *PolicyCreateOpts) (Policy, error) {
	if opt == nil {
		opt = &PolicyCreateOpts{}
	}

	params := struct {
		Name string `json:"name"`
		PolicyCreateOpts
	}{name, *opt}

	var result Policy
	err := p.conn.Request(rest.MethodPost, policiesUri, params, &result)
	return result, err
}
