package emanage

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJsonUnmarshal(t *testing.T) {
	jsonContent := `[
  {
    "report_id": "59d131e9-bf97-4caf-8b52-50813924b2a4",
    "name": "minimal",
    "description": "",
    "ips": [
      "1.1.1.4",
      "1.1.1.3",
      "1.1.1.2"
    ],
    "time": "2016-04-20T11:33:52.919Z"
  },
  {
    "report_id": "59d131d9-bf97-4caf-8b52-50813924b2a4",
    "name": "minimal",
    "description": "",
    "ips": [
      "1.1.1.4",
      "1.1.1.3",
      "1.1.1.2"
    ],
    "time": "2016-04-20T11:33:52.919Z"
  }
]`
	var reportList []ReportListElement

	if err := json.Unmarshal([]byte(jsonContent), &reportList); err != nil {
		t.Fatalf("error returned %+v", err)
	}

	fmt.Printf("%+v", reportList)
}
