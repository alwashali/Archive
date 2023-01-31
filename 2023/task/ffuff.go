package task

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bitly/go-simplejson"
)

type ffuf_resutl struct {
	Id          string `json:"id"`
	Task_Id     string `json:"task_id"`
	Timestamp   string `json:"timestamp"`
	Owner       string `json:"owner"`
	Report_id   string `json:"report_id"`
	Command     string `json:"command"`
	Tool        string `json:"tool"`
	Task_type   string `json:"task_type"`
	Status      string `json:"status"`
	Host        string `json:"host"`
	URL         string `json:"url"`
	Status_Code int    `json:"status_code"`
}

func ffuff(t *task) string {

	//sb := strings.Builder{}
	buf := bytes.NewBuffer(t.Output.Bytes())
	r := bufio.NewScanner(buf)
	r.Split(bufio.ScanLines)
	//fmt.Println(t.Output)

	ffuf_result := []ffuf_resutl{}

	for r.Scan() {

		ffuf := ffuf_resutl{}

		ffuf_json_result, err := simplejson.NewJson([]byte(r.Text()))
		if err != nil {
			fmt.Println("Error decoding ffuf json", err)
		}
		url, err := ffuf_json_result.Get("url").String()
		if err != nil {
			fmt.Println(err)
		}
		host, err := ffuf_json_result.Get("host").String()
		if err != nil {
			fmt.Println(err)
		}
		status, err := ffuf_json_result.Get("status").Int()
		if err != nil {
			fmt.Println(err)
		}

		ffuf.Timestamp = t.Timestamp
		ffuf.Id = t.Id
		ffuf.Command = t.Command
		ffuf.Owner = t.Owner
		ffuf.Report_id = t.Report_id
		ffuf.Status = t.Status
		ffuf.Tool = t.Tool
		ffuf.Task_type = t.task_type
		ffuf.Status_Code = status
		ffuf.URL = url
		ffuf.Host = host

		ffuf_result = append(ffuf_result, ffuf)

	}

	result, err := json.Marshal(ffuf_result)
	if err != nil {
		fmt.Println(err)
	}

	return string(result)
}
