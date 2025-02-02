package topo

import (
	"bufio"
	"errors"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"strings"
)

const KERNEL_DEBUG_INNO_GPUTOPO = "/sys/kernel/debug/inno/gpu_topo"

type Topo struct {
	GPUId        string
	NUMAAffinity string
	CPUAffinity  string
}

func parseGPUTopo() ([]Topo, error) {
	file, err := os.Open(KERNEL_DEBUG_INNO_GPUTOPO)
	if err != nil {
		return nil, errors.New("无法获取信息，请确认使用 root 权限，或 sudo 运行")
	}
	defer file.Close()

	var topos []Topo

	scanner := bufio.NewScanner(file)
	skip := true
	for scanner.Scan() {
		if skip {
			skip = false
			continue
		}
		line := scanner.Text()
		temp := strings.Fields(line)
		topos = append(topos, Topo{
			GPUId:        temp[0],
			NUMAAffinity: temp[len(temp)-2],
			CPUAffinity:  temp[len(temp)-1],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return topos, nil

}

func GetGPUTopoToTable() string {
	topos, err := parseGPUTopo()
	if err != nil {
		return "ERROR：无法获取 GPU 拓扑信息（ " + err.Error() + " )"
	}
	t := table.Table{}
	t.SetTitle("Inno GPU Topo")
	header := table.Row{""}
	for _, topo := range topos {
		header = append(header, topo.GPUId)
	}
	header = append(header, "NUMA Affinity", "CPU Affinity")
	t.AppendHeader(header)

	for _, topo := range topos {
		// TODO
		t.AppendRow(table.Row{topo.GPUId, "X", topo.NUMAAffinity, topo.CPUAffinity})
	}

	legend := `
Legend:
	X: Self
`
	return t.Render() + legend
}
