package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/bitly/go-simplejson"
)

func wpscanf(t *task) string {
	var wpjs map[string]interface{}
	buf := bytes.NewBuffer(t.Output.Bytes())
	if err := json.Unmarshal(buf.Bytes(), &wpjs); err != nil {
		log.Fatal("Failed to marchal wpscan output", err)
	}

	js := simplejson.New()
	js.Set("timestamp", t.Timestamp)
	js.Set("id", t.Id)
	js.Set("owner", t.Owner)
	js.Set("report_id", t.Report_id)
	js.Set("command", t.Command)
	js.Set("tool", t.Tool)
	js.Set("task_type", t.task_type)
	js.Set("interesting_findings_", wpjs["interesting_findings"])

	output, err := js.Encode()
	if err != nil {
		fmt.Println(err)
	}
	return string(output)

}
