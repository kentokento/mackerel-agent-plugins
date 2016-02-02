package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	pipeline "github.com/mattn/go-pipeline"
)

type OpenProcPlugin struct {
	Option    string
	Target    string
	NameLabel string
}

func (m OpenProcPlugin) FetchMetrics() (map[string]interface{}, error) {
	strp, _ := getDataWithCommand(m.Option, m.Target)

	stat := make(map[string]interface{})
	stat[m.NameLabel], _ = strconv.ParseUint(strings.Trim(strp, "\n"), 10, 64)

	return stat, nil
}

func getDataWithCommand(option string, target string) (string, error) {
	out, err := pipeline.Output(
		[]string{"lsof", option, target},
		[]string{"wc", "-l"},
	)
	if err != nil {
		return "0", err
	}

	return string(out), nil
}

func (m OpenProcPlugin) GraphDefinition() map[string](mp.Graphs) {
	return map[string](mp.Graphs){
		"open-process.count": mp.Graphs{
			Label: "Open process count",
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: m.NameLabel, Label: m.NameLabel},
			},
		},
	}

}

var stderrLogger *log.Logger

func getStderrLogger() *log.Logger {
	if stderrLogger == nil {
		stderrLogger = log.New(os.Stderr, "", log.LstdFlags)
	}
	return stderrLogger
}

func main() {
	optPort := flag.String("port", "", "port")
	optName := flag.String("name", "", "name")
	optPid := flag.String("pid", "", "pid")
	optUser := flag.String("user", "", "pid")
	flag.Parse()

	var option, target, nameLabel string
	if *optPort != "" {
		option = "-i"
		target = ":" + *optPort
		nameLabel = "port_" + *optPort
	} else if *optName != "" {
		option = "-c"
		target = *optName
		nameLabel = "name_" + *optName
	} else if *optPid != "" {
		option = "-p"
		target = *optPid
		nameLabel = "pid_" + *optPid
	} else if *optUser != "" {
		option = "-u"
		target = *optUser
		nameLabel = "user_" + *optUser
	} else {
		log.Fatalln("Bad request.")
	}

	var op OpenProcPlugin
	op.Option = option
	op.Target = target
	op.NameLabel = nameLabel

	helper := mp.NewMackerelPlugin(op)

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
