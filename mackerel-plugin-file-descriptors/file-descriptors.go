package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

var graphdef = map[string](mp.Graphs){
	"file-descriptors.handles": mp.Graphs{
		Label: "File handles",
		Unit:  "integer",
		Metrics: [](mp.Metrics){
			mp.Metrics{Name: "allocated", Label: "Allocated"},
			mp.Metrics{Name: "unused", Label: "Unused"},
			mp.Metrics{Name: "maximum", Label: "Max"},
		},
	},
}

type FileDescriptorsPlugin struct {
	Filepath string
}

var metricVarDef = []string{
	"allocated",
	"unused",
	"maximum",
}

func (m FileDescriptorsPlugin) FetchMetrics() (map[string]interface{}, error) {
	var err error
	strp, err := getDataWithCommand(m.Filepath)
	if err != nil {
		return nil, err
	}

	stat := make(map[string]interface{})
	parseVars(strp, &stat)

	return stat, nil
}

func parseVars(text *string, statp *map[string]interface{}) error {
	stat := *statp

	lines := strings.Split(strings.Trim(string(*text), "\n"), "\t")
	for key, _ := range lines {
		stat[metricVarDef[key]], _ = strconv.ParseUint(lines[key], 10, 64)
	}

	return nil
}

func getDataWithCommand(filepath string) (*string, error) {
	cmd := exec.Command("cat", filepath)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	str := out.String()
	return &str, nil
}

func (m FileDescriptorsPlugin) GraphDefinition() map[string](mp.Graphs) {
	return graphdef
}

var stderrLogger *log.Logger

func getStderrLogger() *log.Logger {
	if stderrLogger == nil {
		stderrLogger = log.New(os.Stderr, "", log.LstdFlags)
	}
	return stderrLogger
}

func main() {
	optFilepath := flag.String("filepath", "/proc/sys/fs/file-nr", "file-nr path")
	flag.Parse()

	var fds FileDescriptorsPlugin
	fds.Filepath = *optFilepath

	helper := mp.NewMackerelPlugin(fds)

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
