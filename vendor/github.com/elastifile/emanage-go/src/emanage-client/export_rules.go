package emanage

import (
	"fmt"

	"errors"
	"optional"
	"rest"
)

type ExportRule struct {
	VmId        int                  `json:"vm_id,omitempty"`
	HostId      int                  `json:"host_id,omitempty"`
	ClientMac   string               `json:"client_mac,omitempty"`
	ClientIp    string               `json:"client_ip,omitempty"`
	Access      ExportAccessModeType `json:"access_permission,omitempty"`
	UserMapping UserMappingType      `json:"user_mapping,omitempty"`
	Uid         optional.Int         `json:"uid,omitempty"`
	Gid         optional.Int         `json:"gid,omitempty"`
}

func (rule *ExportRule) Match(other *ExportRule) bool {
	return !(other.VmId != 0 && other.VmId != rule.VmId ||
		other.HostId != 0 && other.HostId != rule.HostId ||
		other.ClientMac != "" && other.ClientMac != rule.ClientMac ||
		other.ClientIp != "" && other.ClientIp != rule.ClientIp ||
		other.Access != "" && other.Access != rule.Access ||
		other.UserMapping != "" && other.UserMapping != rule.UserMapping ||
		other.Uid != nil && (rule.Uid == nil || *other.Uid != *rule.Uid) ||
		other.Gid != nil && (rule.Gid == nil || *other.Gid != *rule.Gid))
}

func (rule *ExportRule) Update(opt *ExportRule) {
	if opt.VmId != 0 {
		rule.VmId = opt.VmId
	}
	if opt.HostId != 0 {
		rule.HostId = opt.HostId
	}
	if opt.ClientMac != "" {
		rule.ClientMac = opt.ClientMac
	}
	if opt.ClientIp != "" {
		rule.ClientIp = opt.ClientIp
	}
	if opt.Access != "" {
		rule.Access = opt.Access
	}
	if opt.UserMapping != "" {
		rule.UserMapping = opt.UserMapping
	}
	if opt.Uid != nil {
		rule.Uid = optional.NewInt(*opt.Uid)
	}
	if opt.Gid != nil {
		rule.Gid = optional.NewInt(*opt.Gid)
	}
}

func (ex *exports) RuleCreate(export *Export, rule *ExportRule) (Export, error) {
	if rule == nil {
		rule = &ExportRule{}
	}

	var result Export
	err := ex.session.Request(rest.MethodPut, ex.createUri(export), rule, &result)
	return result, err
}

func (ex *exports) RuleUpdate(export *Export, rule *ExportRule) (Export, error) {
	var result Export
	if rule == nil {
		return result, fmt.Errorf("missing export %s update rule options", export.Name)
	}
	err := ex.session.Request(rest.MethodPut, ex.updateUri(export), rule, &result)
	return result, err
}

func (ex *exports) RuleDelete(export *Export, rule *ExportRule) error {
	//delRule := ExportRule{
	//	ClientIp:	rule.ClientIp,
	//	HostId:		rule.HostId,
	//	VmId:		rule.VmId,
	//	ClientMac:	rule.ClientMac,
	//}
	if rule == nil {
		return errors.New("Can't operate empty rule")
	}
	return ex.session.Request(rest.MethodPut, ex.deleteUri(export), rule, nil)
}

func (ex *exports) RuleSelect(export *Export, rule *ExportRule) (result []*ExportRule) {
	for _, cr := range export.ClientRules {
		if cr.Match(rule) {
			result = append(result, &cr)
		}
	}
	return result
}

func (ex *exports) createUri(export *Export) string {
	return ex.fullUri(export, "set")
}

func (ex *exports) updateUri(export *Export) string {
	return ex.fullUri(export, "update")
}

func (ex *exports) deleteUri(export *Export) string {
	return ex.fullUri(export, "remove")
}

func (ex *exports) fullUri(export *Export, op string) string {
	return fmt.Sprintf("%s/%d/%s_rule", exportsUri, export.Id, op)
}
