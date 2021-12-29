package adbutilsgo

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strings"
)

func init() {
	pft = make([]*PortForwardRecord, 0)
}

func ListDevices() []string {
	out, err := exec.Command("adb", "devices").Output()
	if err != nil {
		return nil
	}
	lines := strings.Split(string(out), "\n")
	started := false
	devices := make([]string, 0)
	for _, _line := range lines {
		if strings.HasPrefix(_line, "List of devices attached") {
			started = true
			continue
		}
		if started {
			deviceLine := strings.Split(_line, "\t")
			if len(deviceLine) != 2 {
				continue
			}
			devices = append(devices, deviceLine[0])
		}
	}
	return devices
}

func findFreePort() string {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return ""
	}
	defer l.Close()
	port := l.Addr().String()[5:]
	return port
}

type PortForwardRecord struct {
	DeviceId   string
	DevicePort string
	HostPort   string
}

var pft []*PortForwardRecord

func InitializePortForwardTable(tbl []*PortForwardRecord) {
	if tbl != nil {
		pft = tbl
	}
}

type PortForwardOptions struct {
	NoCache bool
}

func PortForward(deviceId string, devicePort string, options PortForwardOptions) (string, error) {
	if !options.NoCache {
		for _, item := range pft {
			if item.DeviceId == deviceId && item.DevicePort == devicePort {
				return item.HostPort, nil
			}
		}
	}

	freePort := findFreePort()
	if freePort == "" {
		return "", errors.New("no free port")
	}
	_, err := exec.Command("adb", "-s", deviceId, "forward", "tcp:"+freePort, "tcp:"+devicePort).Output()
	if err != nil {
		return "", errors.New("adb usb port forward failed")
	}
	pft = append(pft, &PortForwardRecord{
		DeviceId:   deviceId,
		DevicePort: devicePort,
		HostPort:   freePort,
	})
	return freePort, nil
}

func OpenAtxPortForward(deviceId string) (string, error) {
	hostPort, err := PortForward(deviceId, "7912", PortForwardOptions{})
	if err != nil {
		return hostPort, err
	}
	retryTimes := 5
	for ; retryTimes > 0; retryTimes-- {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/version", hostPort))
		if err == nil && resp.StatusCode == 200 {
			return hostPort, nil
		}
		hostPort, err = PortForward(deviceId, "7912", PortForwardOptions{NoCache: true})
		if err != nil {
			return hostPort, err
		}
	}
	return "", errors.New("atx-agent version failed")
}

func Install(deviceId string, apkPath string) error {
	_, err := exec.Command("adb", "-s", deviceId, "install", apkPath).Output()
	if err != nil {
		return err
	}
	return nil
}

func Pull(deviceId string, devicePath string, hostPath string) error {
	var err error
	if hostPath == "" {
		_, err = exec.Command("adb", "-s", deviceId, "pull", devicePath).Output()
	} else {
		_, err = exec.Command("adb", "-s", deviceId, "pull", devicePath, hostPath).Output()
	}

	if err != nil {
		return err
	}
	return nil
}

func Push(deviceId string, hostPath string, devicePath string) error {
	_, err := exec.Command("adb", "-s", deviceId, "push", hostPath, devicePath).Output()

	if err != nil {
		return err
	}
	return nil
}

func Shell(deviceId string, cmd string) (string, error) {
	out, err := exec.Command("adb", "-s", deviceId, "shell", cmd).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
