package types

import "sync"

type ReportStatus struct {
	sync.Mutex
	Expected     int
	Destinations []string
	Failures     []*RemoteError
	Processed    []string
	Enqued       []string
	Waiting      []string
}

func NewReportStatus() *ReportStatus {
	return &ReportStatus{}
}

func (s *ReportStatus) Merge(other *ReportStatus) {
	s.Failures = append(s.Failures, other.Failures...)
	s.Destinations = append(s.Destinations, other.Destinations...)
	s.Processed = append(s.Processed, other.Processed...)
	s.Enqued = append(s.Enqued, other.Enqued...)
	s.Waiting = append(s.Waiting, other.Waiting...)
}

func (s *ReportStatus) Completed() float64 {
	if len(s.Processed) == 0 {
		return 0
	}
	return float64(s.Expected) / float64(len(s.Processed))
}

func (s *ReportStatus) Expect(n int) {
	s.Lock()
	s.Expected += n
	s.Unlock()
}

func (s *ReportStatus) Fail(err error) {
	s.Lock()
	s.Failures = append(s.Failures, NewRemoteError(err))
	s.Unlock()
}

func (s *ReportStatus) Append(files []string) {
	s.Lock()
	s.Waiting = append(s.Waiting, files...)
	s.Unlock()
}

func (s *ReportStatus) Enque(file string) {
	s.Lock()
	defer s.Unlock()

	var waiting []string
	for _, f := range s.Waiting {
		if f != file {
			waiting = append(waiting, f)
		}
	}
	s.Waiting = waiting
	s.Enqued = append(s.Enqued, file)
}

func (s *ReportStatus) Finish(file string, err error) {
	s.Lock()
	defer s.Unlock()

	var enqued []string
	for _, f := range s.Enqued {
		if f != file {
			enqued = append(enqued, f)
		}
	}
	s.Enqued = enqued
	s.Processed = append(s.Processed, file)
	if err != nil {
		s.Failures = append(s.Failures, NewRemoteError(err))
	}
}
