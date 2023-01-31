package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/bitly/go-simplejson"
)

type Dnsx_result struct {
	Host        string
	Resolver    string
	A           []string
	AAAA        []string
	MX          []string
	SOA         []string
	NS          []string
	TXT         []string
	Status_code string
}

func dnsxf(t *task) string {
	//fmt.Println(t.Output)
	js := simplejson.New()

	js.Set("timestamp", t.Timestamp)
	js.Set("id", t.Id)
	js.Set("owner", t.Owner)
	js.Set("report_id", t.Report_id)
	js.Set("command", t.Command)
	js.Set("tool", t.Tool)
	js.Set("task_type", t.task_type)

	buf := bytes.NewBuffer(t.Output.Bytes())

	var dnsx map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &dnsx); err != nil {
		log.Fatal("Failed to marchal wpscan output", err)
	}
	js.Set("dnsx", dnsx)

	out, err := js.MarshalJSON()
	if err != nil {
		fmt.Println(err)
	}
	return string(out)
}
