package version

import (
	"bufio"
	"errors"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"regexp"
)

const KERNEL_DEBUG_INNO_VERSION = "/sys/kernel/debug/inno/version"

type InnoVersion struct {
	Key   string
	Value string
}

func parseInnoVersion() ([]InnoVersion, error) {
	file, err := os.Open(KERNEL_DEBUG_INNO_VERSION)
	if err != nil {
		return nil, errors.New("无法获取信息，请确认使用 root 权限，或 sudo 运行")
	}
	defer file.Close()

	var versions []InnoVersion

	driverRegex := regexp.MustCompile(`^Driver Version: (.+)$`)
	deviceNameRegex := regexp.MustCompile(`^Device Name: (.+)$`)
	deviceIDRegex := regexp.MustCompile(`^Device ID: (.+)$`)
	deviceVersionRegex := regexp.MustCompile(`^Device Version: (.+)$`)
	gpuVariantRegex := regexp.MustCompile(`^GPU variant BVNC: (.+)$`)
	firmwareRegex := regexp.MustCompile(`^Firmware Version: (.+)$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := driverRegex.FindStringSubmatch(line); matches != nil {
			versions = append(versions, InnoVersion{
				Key:   "Driver Version",
				Value: matches[1],
			})
		} else if matches := deviceNameRegex.FindStringSubmatch(line); matches != nil {
			versions = append(versions, InnoVersion{
				Key:   "Device Name",
				Value: matches[1],
			})
		} else if matches := deviceIDRegex.FindStringSubmatch(line); matches != nil {
			versions = append(versions, InnoVersion{
				Key:   "Device ID",
				Value: matches[1],
			})
		} else if matches := deviceVersionRegex.FindStringSubmatch(line); matches != nil {
			versions = append(versions, InnoVersion{
				Key:   "Device Version",
				Value: matches[1],
			})
		} else if matches := gpuVariantRegex.FindStringSubmatch(line); matches != nil {
			versions = append(versions, InnoVersion{
				Key:   "GPU Variant BVNC",
				Value: matches[1],
			})
		} else if matches := firmwareRegex.FindStringSubmatch(line); matches != nil {
			versions = append(versions, InnoVersion{
				Key:   "Firmware Version",
				Value: matches[1],
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return versions, nil
}

func GetVersionToTable() string {
	vers, err := parseInnoVersion()
	if err != nil {
		return "ERROR：无法获取 GPU 版本信息（ " + err.Error() + " )"
	}
	t := table.Table{}
	t.SetTitle("Inno GPU Version")

	for _, version := range vers {
		t.AppendRow(table.Row{version.Key, version.Value})
	}
	return t.Render()
}
