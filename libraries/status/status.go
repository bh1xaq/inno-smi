package status

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"inno-smi/libraries/mem"
	"os"
	"regexp"
	"time"
)

const KERNEL_DEBUG_INNO_STATUS = "/sys/kernel/debug/inno/status"

type InnoStatus struct {
	DriverStatus      string
	DeviceID          string
	FirmwareStatus    string
	ServerErrors      string
	HWREventCount     string
	CRREventCount     string
	SLREventCount     string
	WGPErrorCount     string
	TRPErrorCount     string
	FWFEventCount     string
	APMEventCount     string
	GPUUtilisation    string
	TDMUtilisation    string
	GEOMUtilisation   string
	ThreeDUtilisation string
	CDMUtilisation    string
}

func parseInnoStatus() (*InnoStatus, error) {
	file, err := os.Open(KERNEL_DEBUG_INNO_STATUS)
	if err != nil {
		return nil, errors.New("无法获取信息，请确认使用 root 权限，或 sudo 运行")
	}
	defer file.Close()

	var status InnoStatus

	driverStatusRegex := regexp.MustCompile(`^Driver Status:\s+(.+)$`)
	deviceIDRegex := regexp.MustCompile(`^Device ID:\s+(.+)$`)
	firmwareStatusRegex := regexp.MustCompile(`^Firmware Status:\s+(.+)$`)
	serverErrorsRegex := regexp.MustCompile(`^Server Errors:\s+(.+)$`)
	hwrEventRegex := regexp.MustCompile(`^HWR Event Count:\s+(.+)$`)
	crrEventRegex := regexp.MustCompile(`^CRR Event Count:\s+(.+)$`)
	slrEventRegex := regexp.MustCompile(`^SLR Event Count:\s+(.+)$`)
	wgpErrorRegex := regexp.MustCompile(`^WGP Error Count:\s+(.+)$`)
	trpErrorRegex := regexp.MustCompile(`^TRP Error Count:\s+(.+)$`)
	fwfEventRegex := regexp.MustCompile(`^FWF Event Count:\s+(.+)$`)
	apmEventRegex := regexp.MustCompile(`^APM Event Count:\s+(.+)$`)
	gpuUtilisationRegex := regexp.MustCompile(`^GPU Utilisation:\s+(.+)$`)
	tdmUtilisationRegex := regexp.MustCompile(`^TDM\s+Utilisation:\s+(.+)$`)
	geomUtilisationRegex := regexp.MustCompile(`^GEOM\s+Utilisation:\s+(.+)$`)
	threeDUtilisationRegex := regexp.MustCompile(`^3D\s+Utilisation:\s+(.+)$`)
	cdmUtilisationRegex := regexp.MustCompile(`^CDM\s+Utilisation:\s+(.+)$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := driverStatusRegex.FindStringSubmatch(line); matches != nil {
			status.DriverStatus = matches[1]
		} else if matches := deviceIDRegex.FindStringSubmatch(line); matches != nil {
			status.DeviceID = matches[1]
		} else if matches := firmwareStatusRegex.FindStringSubmatch(line); matches != nil {
			status.FirmwareStatus = matches[1]
		} else if matches := serverErrorsRegex.FindStringSubmatch(line); matches != nil {
			status.ServerErrors = matches[1]
		} else if matches := hwrEventRegex.FindStringSubmatch(line); matches != nil {
			status.HWREventCount = matches[1]
		} else if matches := crrEventRegex.FindStringSubmatch(line); matches != nil {
			status.CRREventCount = matches[1]
		} else if matches := slrEventRegex.FindStringSubmatch(line); matches != nil {
			status.SLREventCount = matches[1]
		} else if matches := wgpErrorRegex.FindStringSubmatch(line); matches != nil {
			status.WGPErrorCount = matches[1]
		} else if matches := trpErrorRegex.FindStringSubmatch(line); matches != nil {
			status.TRPErrorCount = matches[1]
		} else if matches := fwfEventRegex.FindStringSubmatch(line); matches != nil {
			status.FWFEventCount = matches[1]
		} else if matches := apmEventRegex.FindStringSubmatch(line); matches != nil {
			status.APMEventCount = matches[1]
		} else if matches := gpuUtilisationRegex.FindStringSubmatch(line); matches != nil {
			status.GPUUtilisation = matches[1]
		} else if matches := tdmUtilisationRegex.FindStringSubmatch(line); matches != nil {
			status.TDMUtilisation = matches[1]
		} else if matches := geomUtilisationRegex.FindStringSubmatch(line); matches != nil {
			status.GEOMUtilisation = matches[1]
		} else if matches := threeDUtilisationRegex.FindStringSubmatch(line); matches != nil {
			status.ThreeDUtilisation = matches[1]
		} else if matches := cdmUtilisationRegex.FindStringSubmatch(line); matches != nil {
			status.CDMUtilisation = matches[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &status, nil
}

func GetMainInnoStatusToTable() string {
	status, err := parseInnoStatus()
	if err != nil {
		return "ERROR：无法获取 GPU 信息（ " + err.Error() + " )"
	}
	ttime := time.Now()

	countTable := table.Table{}
	countTable.SetTitle("GPU Count")
	countTable.Style().Size.WidthMax = 95
	countTable.Style().Size.WidthMin = 95
	countTable.AppendHeader(table.Row{"Key", "Value"})
	countTable.AppendRows([]table.Row{
		table.Row{"HWR Event", status.HWREventCount},
		table.Row{"CRR Event", status.CRREventCount},
		table.Row{"SLR Event", status.SLREventCount},
		table.Row{"WGP Error", status.WGPErrorCount},
		table.Row{"TRP Error", status.TRPErrorCount},
		table.Row{"FWF Error", status.FWFEventCount},
		table.Row{"APM Error", status.APMEventCount},
		table.Row{"Server Errors", status.ServerErrors},
	})
	utilisationTable := table.Table{}
	utilisationTable.SetTitle(fmt.Sprintf("GPU Utilisation: %s%%", status.GPUUtilisation))
	utilisationTable.Style().Size.WidthMax = 95
	utilisationTable.Style().Size.WidthMin = 95
	utilisationTable.AppendHeader(table.Row{"Key", "Value"})
	utilisationTable.AppendRows([]table.Row{
		table.Row{"TDM", status.TDMUtilisation},
		table.Row{"GEOM", status.GEOMUtilisation},
		table.Row{"3D", status.ThreeDUtilisation},
		table.Row{"CDM", status.CDMUtilisation},
	})
	mainTable := table.Table{}
	mainTable.SetTitle("INNO-SMI For Fantasy II-M     " + "Driver Status:" + status.DriverStatus + "     Firmware Status:" + status.FirmwareStatus)
	mainTable.Style().Size.WidthMax = 100
	mainTable.Style().Size.WidthMin = 100
	mainTable.AppendRows([]table.Row{
		table.Row{utilisationTable.Render()},
		table.Row{countTable.Render()},
	})
	progress := "Progress:"
	meminfo := mem.GetInfoToTable()

	return ttime.String() + "\n" + mainTable.Render() + "\n" + progress + "\n" + meminfo
}
