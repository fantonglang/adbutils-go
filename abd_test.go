package adbutilsgo

import (
	"fmt"
	"testing"
)

func TestListDevices(t *testing.T) {
	devices := ListDevices()
	if devices == nil {
		t.Error("no devices")
		return
	}
	for _, d := range devices {
		t.Log(fmt.Sprintln("device: " + d))
		return
	}
}

func TestFindFreePort(t *testing.T) {
	port := findFreePort()
	if port == "" {
		t.Error("find free port failed")
		return
	}
	t.Log(fmt.Sprintln(port))
}

func TestPortForward(t *testing.T) {
	deviceId := "c574dd45"
	devicePort := "7912"
	hostPort, err := PortForward(deviceId, devicePort, PortForwardOptions{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(fmt.Sprintln(hostPort))
}

func TestOpenAtxPortForward(t *testing.T) {
	deviceId := "c574dd45"
	hostPort, err := OpenAtxPortForward(deviceId)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(fmt.Sprintln(hostPort))
}

// func TestPull(t *testing.T) {
// 	deviceId := "c574dd45"
// 	err := Pull(deviceId, "/sdcard/myscreenrec.apk", "tmp/")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// }

func TestShell(t *testing.T) {
	deviceId := "c574dd45"
	out, err := Shell(deviceId, "am start -n com.ss.android.ugc.aweme/.splash.SplashActivity")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(out)
}
