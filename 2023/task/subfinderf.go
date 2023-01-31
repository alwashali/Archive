package task

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bitly/go-simplejson"
)

type subfinder_result struct {
	Id        string `json:"id"`
	Task_Id   string `json:"task_id"`
	Timestamp string `json:"timestamp"`
	Owner     string `json:"owner"`
	Report_id string `json:"report_id"`
	Command   string `json:"command"`
	Tool      string `json:"tool"`
	Task_type string `json:"task_type"`
	Status    string `json:"status"`
	Host      string `json:"host"`
	Source    string `json:"source"`
}

func subfinderf(t *task) string {

	buf := bytes.NewBuffer(t.Output.Bytes())
	r := bufio.NewScanner(buf)
	r.Split(bufio.ScanLines)

	var subfinderf_array []subfinder_result

	for r.Scan() {
		sff := subfinder_result{}
		subfinderL, err := simplejson.NewJson(r.Bytes())
		if err != nil {
			fmt.Println(err)
		}
		host, err := subfinderL.Get("host").String()
		if err != nil {
			fmt.Println(err)
		}
		source, err := subfinderL.Get("source").String()
		if err != nil {
			fmt.Println(err)
		}

		sff.Timestamp = t.Timestamp
		sff.Id = t.Id
		sff.Command = t.Command
		sff.Owner = t.Owner
		sff.Report_id = t.Report_id
		sff.Status = t.Status
		sff.Tool = t.Tool
		sff.Task_type = t.task_type
		sff.Source = source
		sff.Host = host

		subfinderf_array = append(subfinderf_array, sff)

	}

	result, err := json.Marshal(subfinderf_array)
	if err != nil {
		fmt.Println(err)
	}

	return string(result)
}
