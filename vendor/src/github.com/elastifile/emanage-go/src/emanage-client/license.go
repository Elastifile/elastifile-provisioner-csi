package emanage

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type License struct {
	Hosts          int              `json:"hosts"`
	RawCapacity    RawCapacityBytes `json:"raw_capacity"`
	ExpirationDate string           `json:"expiration_date"`
	Comply         bool             `json:"comply"`
	Errors         []string         `json:"errors"`
}

type RawCapacityBytes struct {
	Bytes float64 `json:"bytes"`
}

// Example
// {
// 	"hosts":3,
// 	"raw_capacity":54975581388800.0,
// 	"expiration_date":"2018-04-04",
// 	"comply":false,
// 	"errors":["Exceeded license number of hosts. Limit is 3, using 4"]
// }

type LicenseOpts struct {
	License string `json:"license"`
}

func (opts *LicenseOpts) LoadFromFile(filepath string) error {
	if filepath == "" {
		return errors.New("No file was specified")
	}

	f, err := os.Open(filepath)
	if os.IsNotExist(err) {
		return fmt.Errorf("File does not exist, file='%v'", filepath)
	}
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	opts.License = string(data)

	return nil
}
