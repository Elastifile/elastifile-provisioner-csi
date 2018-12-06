package issues

import (
	"fmt"

	log "gopkg.in/inconshreveable/log15.v2"
)

var issues = make(map[string]error)

func ReportIssue(issue string, summary string, err error) {
	description := fmt.Sprintf("%v: %v", issue, summary)
	log.Info("Reported issue", "description", description, "err", err)
	issues[description] = err
}

func Issues() map[string]error {
	return issues
}
