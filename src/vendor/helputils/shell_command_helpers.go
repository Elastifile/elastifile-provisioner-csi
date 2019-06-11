package helputils

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-errors/errors"
)

func CurrentUser() (*user.User, error) {
	u, err := user.Current()
	if err != nil {
		return nil, errors.New(err)
	}
	return u, nil
}

func IsRootUser() (bool, error) {
	currentUser, err := CurrentUser()
	if err != nil {
		return false, err
	}
	return currentUser.Uid == "0", nil
}

func Grep(text string, regex string) []string {
	re := regexp.MustCompile(regex)
	return FilterStr(
		strings.Split(text, "\n"),
		func(line string) bool { return re.MatchString(line) },
	)
}

func ExecuteShellString(script string) ([]byte, error) {
	var err error
	tmpFile, err := TmpRandomFileName("tesla-script")
	if err != nil {
		return nil, err
	}
	teslaScriptPath := tmpFile + ".sh"

	if !strings.HasSuffix(script, "\n") {
		script += "\n"
	}

	if err := ioutil.WriteFile(teslaScriptPath, []byte(script), os.ModePerm); err != nil {
		return []byte{}, err
	}

	out, err := exec.Command("/bin/bash", teslaScriptPath).CombinedOutput()
	if err == nil {
		os.Remove(teslaScriptPath)
	}

	return out, err
}

func CreateLink(link string, target string, mode os.FileMode) error {
	if err := MkdirAll(filepath.Dir(link), mode); err != nil {
		return err
	}

	if out, err := exec.Command("ln", "-sf", target, link).CombinedOutput(); err != nil {
		return fmt.Errorf("failed linking target: %s, link: %s, err: %s, out: %s", target, link, err, out)
	}

	return nil
}

func MkdirAll(path string, mode os.FileMode) error {
	if err := os.MkdirAll(path, mode); err != nil {
		return errors.Wrap(err, 0)
	}

	if err := ForceFileMode(path, mode); err != nil {
		return err
	}

	return nil
}

func TouchFile(targetPath string, mode os.FileMode) error {
	if file, err := os.OpenFile(targetPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, mode); err != nil {
		return err
	} else {
		if err = file.Close(); err != nil {
			return err
		}
		if err = ForceFileMode(targetPath, mode); err != nil {
			return err
		}
	}
	return nil
}

func AppendFile(targetPath string, sourcePath string, offset int64, mode os.FileMode) (int, error) {
	var size int

	if data, err := ReadFrom(sourcePath, offset); err != nil {
		return 0, err
	} else if file, err := os.OpenFile(targetPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, mode); err != nil {
		return 0, err
	} else {
		_, err := file.Write(data)
		if err != nil {
			return 0, err
		}
		if err := file.Close(); err != nil {
			return 0, err
		}
		size = len(data)
	}
	if err := ForceFileMode(targetPath, mode); err != nil {
		return size, err
	}

	return size, nil
}

// workaround Go's bug in os.OpenFile() which does not change file's mode as requested
func ForceFileMode(path string, mode os.FileMode) error {
	if fstat, err := os.Stat(path); err != nil {
		return err
	} else if fstat.Mode() != mode {
		cmd := fmt.Sprintf("chmod %o %s", mode, path)
		isRootUser, err := IsRootUser()
		if err != nil {
			return err
		}
		if !isRootUser {
			cmd = "/usr/bin/sudo " + cmd
		}
		if out, err := ExecuteShellString(cmd); err != nil {
			return fmt.Errorf("%s: %s %s", err, out, path)
		}
	}
	return nil
}

func TmpRandomFileName(prefix string) (string, error) {
	currentUser, err := CurrentUser()
	if err != nil {
		return "", err
	}
	currUser := currentUser.Username
	if prefix != "" {
		prefix += "-"
	}
	return filepath.Join("/tmp", currUser+"_"+prefix+time.Now().Format("20060102-150405.000")), nil
}

func CpuCount() (result int) {
	cpusByt, err := ExecuteShellString("nproc")
	if err != nil {
		println(err)
		return 0
	}
	cpusByt = cpusByt[:len(cpusByt)-1]

	result, err = strconv.Atoi(string(cpusByt))
	if err != nil {
		return 0
	}

	return
}

func ResolvedHost(host string) string {
	if net.ParseIP(host) == nil {
		addrs, err := net.LookupHost(host)
		if err == nil && len(addrs) > 0 {
			return addrs[0]
		}
		return ""
	}
	return host
}
