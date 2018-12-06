package emanage

import (
	"rest"
)

const sysReportUri = "api/cluster_reports/recent"

type clusterReports struct {
	conn *rest.Session
}

type ClusterReport struct {
	Id                     int    `json:"id"`
	SystemID               int    `json:"system_id"`
	Timestamp              string `json:"timestamp"`
	RocTransitionTotal     int    `json:"roc_transition_total"`
	RocTransitionDone      int    `json:"roc_transition_done"`
	OwnerShipRecoveryTotal int    `json:"ownership_recovery_total"`
	OwnerShipRecoveryDone  int    `json:"ownership_recovery_done"`
	RocTransitionProgress  int    `json:"roc_transition_progress"`
	RocTransitionID        int    `json:"roc_transition_id"`
	OrcTransitionTotal     int    `json:"orc_transition_total"`
	OrcTransitionDone      int    `json:"orc_transition_done"`
	OrcTransitionProgress  int    `json:"orc_transition_progress"`
	OrcTransitionID        int    `json:"orc_transition_id"`
	EcdbTransitionTotal    int    `json:"ecdb_transition_total"`
	EcdbTransitionDone     int    `json:"ecdb_transition_done"`
	EcdbTransitionProgress int    `json:"ecdb_transition_progress"`
}

func (cr *clusterReports) GetAll() (result []ClusterReport, err error) {
	if err = cr.conn.Request(rest.MethodGet, sysReportUri, nil, &result); err != nil {
		logger.Error("GetAll Error", "err", err)
		return nil, err
	}
	return result, nil
}
