package main

import (
	"encoding/json"
	"fmt"
	"log"
    "os"
	"os/exec"
    "path/filepath"
	"strings"
    "strconv"
)

func usage() {
    fmt.Println("Usage:", filepath.Base(os.Args[0]),"[options]")
    fmt.Println()
    fmt.Println("Options:")
    fmt.Println("\t--device (-d) <device>\t", "The name of the device (case sensitive)")
    fmt.Println("\t--udid (-u, -udid)\t", "Only display the device UDID")
    fmt.Println()
}

type device struct {
	Name                 string `json:"deviceName"`
	Class                string `json:"deviceClass"`
	UDID                 string `json:"deviceIdentifier"`
	IMEI                 string `json:"deviceIMEI"`
	SerialNumber         string `json:"deviceSerialNumber"`
	SoftwareVersion      string `json:"deviceSoftwareVersion"`
	ProductVersion       string `json:"productVersion"`
	BuildVersion         string `json:"buildVersion"`
	Architecture         string `json:"deviceArchitecture"`
	DevelopmentStatus    string `json:"deviceDevelopmentStatus"`
	Capacity             string `json:"deviceCapacity"`
	WifiMAC              string `json:"deviceWiFiMAC"`
	BluetoothMAC         string `json:"deviceBluetoothMAC"`
	Type                 string `json:"deviceType"`
	PlatformID           string `json:"platformIdentifier"`
	ColorString          string `json:"deviceColorString"`
	EnclosureColorString string `json:"deviceEnclosureColorString"`
}

func (d *device) String() (s string) {
    cap, _ := strconv.Atoi(d.Capacity)
    s =  fmt.Sprintf("\033[1;33mName:                  \033[0m%s\n", d.Name)
    s += fmt.Sprintf("\033[1;33mClass:                 \033[0m%s\n", d.Class)
    s += fmt.Sprintf("\033[1;33mUDID:                  \033[0m%s\n", d.UDID)
    s += fmt.Sprintf("\033[1;33mIMEI:                  \033[0m%s\n", d.IMEI)
    s += fmt.Sprintf("\033[1;33mSerial Number:         \033[0m%s\n", d.SerialNumber)
    s += fmt.Sprintf("\033[1;33mSoftware Version:      \033[0m%s\n", d.SoftwareVersion)
    s += fmt.Sprintf("\033[1;33mArchitecture:          \033[0m%s\n", d.Architecture)
    s += fmt.Sprintf("\033[1;33mDevelopment Status:    \033[0m%s\n", d.DevelopmentStatus)
    s += fmt.Sprintf("\033[1;33mCapacity:              \033[0m%.2f GB\n", float64(cap) / (1024.0 * 1024.0 * 1000.0))
    s += fmt.Sprintf("\033[1;33mWiFi MAC Address:      \033[0m%s\n", d.WifiMAC)
    s += fmt.Sprintf("\033[1;33mBluetooth MAC Address: \033[0m%s\n", d.BluetoothMAC)
    s += fmt.Sprintf("\033[1;33mColor:                 \033[0m%s\n", d.ColorString)
    s += fmt.Sprintf("\033[1;33mEnclosure Color:       \033[0m%s\n", d.EnclosureColorString)
    s += fmt.Sprintf("\033[1;33mType:                  \033[0m%s\n", d.Type)
    s += fmt.Sprintf("\033[1;33mPlatform Identifier:   \033[0m%s\n", d.PlatformID)
    return
}

func parseSavedDevices() []device {
	output, err := exec.Command("defaults", "read", "com.apple.dt.Xcode", "DVTSavediPhoneDevices").Output()
	if err != nil {
		log.Fatal(err)
	}

	defaults := strings.Replace(string(output), "        {\n", "{", -1)
	defaults = strings.Replace(defaults, "        ", "\"", -1)
	defaults = strings.Replace(defaults, " = ", "\": \"", -1)
	defaults = strings.Replace(defaults, "\";\n", "\",\n", -1)
	defaults = strings.Replace(defaults, ";\n", "\",\n", -1)
	defaults = strings.Replace(defaults, ",\n    }", "\"}", -1)
	defaults = strings.Replace(defaults, "\"\"", "\"", -1)
	defaults = strings.Replace(defaults, "    }\n", "}", -1)
	defaults = strings.Replace(defaults, "    },\n", "},", -1)
	defaults = strings.Replace(defaults, "(\n", "[", -1)
	defaults = strings.Replace(defaults, ")\n", "]\n", -1)

	var devices []device
	err = json.Unmarshal([]byte(defaults), &devices)
	if err != nil {
		log.Fatal(err)
	}

	return devices
}

func main() {
	devices := parseSavedDevices()

    var deviceIDOnly bool
    var deviceName string

    if len(os.Args) == 2 {
        if os.Args[1] == "-h" || os.Args[1] == "--help" {
            usage()
            os.Exit(0)
        } else if os.Args[1] == "-u" || os.Args[1] == "-udid" || os.Args[1] == "--udid" {
            deviceIDOnly = true
        }
    } else if len(os.Args) == 3 || len(os.Args) == 4 {
        for i, arg := range os.Args {
            if arg == "-d" || arg == "--device" {
                deviceName = os.Args[i+1]
            } else if arg == "-u" || arg == "-udid" || arg == "--udid" {
                deviceIDOnly = true
            }
        }
    }

    for _, device := range devices {
        if len(deviceName) > 0 {
            if device.Name == deviceName {
                if deviceIDOnly == true {
/*                    fmt.Printf("\033[1;33mName: \033[0m\"%s\"\n", device.Name)
                    fmt.Printf("\033[1;33mUDID: \033[0m%s\n", device.UDID)*/
                    fmt.Println(device.UDID)
                } else {
                    fmt.Println(device.String())
                }
            }
        } else {
            if deviceIDOnly == true {
                fmt.Printf("\033[1;33mName: \033[0m%s\n", device.Name)
                fmt.Printf("\033[1;33mUDID: \033[0m%s\n", device.UDID)
                fmt.Println()
            } else {
                fmt.Println(device.String())
            }
        }
    }
}
