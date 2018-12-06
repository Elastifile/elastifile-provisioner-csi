package emanage

import (
	"fmt"

	"github.com/go-errors/errors"

	"rest"
)

const snapshotsURI = "/api/snapshots"

type Snapshot struct {
	ID                   int    `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	DataContainerID      int    `json:"data_container_id,omitempty"`
	ConsistencyGroupDcID []int  `json:"data_container_ids,omitempty"`
	UUID                 string `json:"uuid,omitempty"`
	IsWriteable          bool   `json:"is_writeable,omitempty"`
	Locked               bool   `json:"locked,omitempty"`
	Status               string `json:"status,omitempty"`
	CreatedAt            string `json:"created_at,omitempty"`
	UpdatedAt            string `json:"updated_at,omitempty"`
	URL                  string `json:"url,omitempty"`
	conn                 *rest.Session
}

type SnapshotStatistics struct {
	ID               int         `json:"id"`
	SnapshotID       int         `json:"snapshot_id"`
	ReadNumEvents    int         `json:"read_num_events"`
	WriteNumEvents   int         `json:"write_num_events"`
	MdReadNumEvents  int         `json:"md_read_num_events"`
	MdWriteNumEvents int         `json:"md_write_num_events"`
	ReadIo           statIo      `json:"read_io"`
	WriteIo          statIo      `json:"write_io"`
	ReadLatency      statLatency `json:"read_latency"`
	MdReadLatency    statLatency `json:"md_read_latency"`
	MdWriteLatency   statLatency `json:"md_write_latency"`
	WriteLatency     statLatency `json:"write_latency"`
	Timestamp        string      `json:"timestamp"`
	conn             *rest.Session
}
type SnapshotLock struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	SnapshotID int    `json:"snapshot_id"`
	UUID       string `json:"uuid"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type SnapshotLockList []*SnapshotLock

/*
GET /api/snapshots
200
[
  {
    "id": 1,
    "name": "snapshot13",
    "data_container_id": 1,
    "uuid": "dd46c39a-4d72-400d-9757-2f08e91dc9b8",
    "is_writeable": null,
    "locked": false,
    "status": "status_adding",
    "created_at": "2016-01-01T10:02:00.000Z",
    "updated_at": "2016-01-01T10:02:00.000Z",
    "url": "http://test.host/api/snapshots/1"
  }
]*/

type SnapshotList []*Snapshot

type snapshots struct {
	conn *rest.Session
}

func (snaps *snapshots) Get() (SnapshotList, error) {
	var snapsList SnapshotList
	err := snaps.conn.Request(rest.MethodGet, snapshotsURI, nil, &snapsList)
	if err != nil {
		return snapsList, err
	}
	for i, _ := range snapsList {
		snapsList[i].conn = snaps.conn
	}
	return snapsList, nil
}

func (snaps *snapshots) GetById(ID int) (*Snapshot, error) {
	var snap Snapshot
	fullURI := fmt.Sprintf("%s/%d", snapshotsURI, ID)
	err := snaps.conn.Request(rest.MethodGet, fullURI, nil, &snap)
	if err != nil {
		return nil, err
	}
	snap.conn = snaps.conn
	return &snap, nil
}

func (snaps *snapshots) GetByDataContainerId(dcId int) ([]*Snapshot, error) {
	var snapshotsList SnapshotList
	snapsList, err := snaps.Get()
	if err != nil {
		return nil, err
	}
	for _, s := range snapsList {
		if s.DataContainerID == dcId {
			snapshotsList = append(snapshotsList, s)
		}
	}
	return snapshotsList, nil
}

func (snaps *snapshots) Update(snap *Snapshot) (*Snapshot, error) {
	if snap == nil {
		return nil, errors.Errorf("Cannot update snapshot with no data")
	}
	result := &Snapshot{}
	fullURI := fmt.Sprintf("%s/%d", snapshotsURI, snap.ID)
	err := snaps.conn.Request(rest.MethodPut, fullURI, snap, result)
	return result, err
}

func (snaps *snapshots) Create(snap *Snapshot) (*Snapshot, error) {
	if snap == nil {
		return nil, errors.Errorf("Cannot create nil snapshot")
	}
	result := &Snapshot{}
	err := snaps.conn.Request(rest.MethodPost, snapshotsURI, snap, result)
	result.conn = snaps.conn
	return result, err
}
func (snaps *snapshots) CreateConsistencyGroup(snap *Snapshot) (*[]Snapshot, error) {
	if snap == nil {
		return nil, errors.Errorf("Cannot create snapshot with no data")
	}
	fmt.Printf("snapshot data %#v", snap)
	result := &[]Snapshot{}
	uri := fmt.Sprintf("%s/consistency_group", snapshotsURI)
	err := snaps.conn.Request(rest.MethodPost, uri, snap, result)
	for _, rr := range *result {
		rr.conn = snaps.conn
	}
	// result.conn = snaps.conn
	return result, err

}

func (snap *Snapshot) Delete() (tasks AsyncTasks) {
	fullURI := fmt.Sprintf("%s/%d", snapshotsURI, snap.ID)
	var err error
	if tasks.taskIDs, err = snap.conn.AsyncRequest(rest.MethodDelete, fullURI, nil); err != nil {
		tasks.err = err
	} else {
		tasks.conn = snap.conn
	}
	return tasks
}

func (snap *Snapshot) Statistics() (*SnapshotStatistics, error) {
	uri := fmt.Sprintf("%s/%d/statistics", snapshotsURI, snap.ID)
	result := &SnapshotStatistics{}
	err := snap.conn.Request(rest.MethodGet, uri, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (snap *Snapshot) CreateLock(lock *SnapshotLock) (*SnapshotLock, error) {
	fullURI := fmt.Sprintf("%s/%d/locks", snapshotsURI, snap.ID)
	result := &SnapshotLock{}

	err := snap.conn.Request(rest.MethodPost, fullURI, lock, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (snap *Snapshot) ListLocks() (*SnapshotLockList, error) {
	fullURI := fmt.Sprintf("%s/%d/locks", snapshotsURI, snap.ID)
	result := &SnapshotLockList{}

	err := snap.conn.Request(rest.MethodGet, fullURI, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (snap *Snapshot) DeleteLock(id int) (*SnapshotLock, error) {
	fullURI := fmt.Sprintf("%s/%d/locks/%d", snapshotsURI, snap.ID, id)
	result := &SnapshotLock{}
	err := snap.conn.Request(rest.MethodDelete, fullURI, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
