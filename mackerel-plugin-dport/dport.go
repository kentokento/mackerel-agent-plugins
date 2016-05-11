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

type DportPlugin struct {
	Option    string
	Target    string
	NameLabel string
}

func (m DportPlugin) FetchMetrics() (map[string]interface{}, error) {
	strp, _ := getDataWithCommand(m.Target)

	stat := make(map[string]interface{})
	stat[m.NameLabel], _ = strconv.ParseUint(strings.Trim(strp, "\n"), 10, 64)

	return stat, nil
}

func getDataWithCommand(target string) (string, error) {
	out, err := pipeline.Output(
		[]string{"ss", target},
		[]string{"wc", "-l"},
	)
	if err != nil {
		return "0", err
	}

	return string(out), nil
}

func (m DportPlugin) GraphDefinition() map[string](mp.Graphs) {
	return map[string](mp.Graphs){
		"dport.count": mp.Graphs{
			Label: "Dport process count",
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
	flag.Parse()

	var target, nameLabel string
	if *optPort != "" {
		target = "( dport = :" + *optPort + " )"
		nameLabel = "port_" + *optPort
	} else {
		log.Fatalln("Bad request.")
	}

	var op DportPlugin
	op.Target = target
	op.NameLabel = nameLabel

	helper := mp.NewMackerelPlugin(op)

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		helper.OutputDefinitions()
	} else {
		helper.OutputValues()
	}
}
