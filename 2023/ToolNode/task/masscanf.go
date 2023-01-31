package task

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type response struct {
	Id        string `json:"id"`
	Task_Id   string `json:"task_id"`
	Timestamp string `json:"timestamp"`
	Owner     string `json:"owner"`
	Report_id string `json:"report_id"`
	Command   string `json:"command"`
	Tool      string `json:"tool"`
	Task_type string `json:"task_type"`
	Ip        string `json:"ip"`
	Ports     []struct {
		Port   int    `json:"port"`
		Status string `json:"status"`
		Proto  string `json:"proto"`
		Reason string `json:"reason"`
		TTL    int    `json:"ttl"`
	}
}

type masscan struct {
	Ip        string `json:"ip"`
	Timestamp string `json:"timestamp"`
	Ports     []struct {
		Port   int    `json:"port"`
		Status string `json:"status"`
		Proto  string `json:"proto"`
		Reason string `json:"reason"`
		TTL    int    `json:"ttl"`
	}
}

func masscanf(t *task) string {
	masscan_array := []*masscan{}
	resp := []response{}

	buf := bytes.NewBuffer(t.Output.Bytes())

	err := json.Unmarshal(buf.Bytes(), &masscan_array)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println("len:", len(masscan_array))
	if len(masscan_array) > 0 {
		for i := 0; i < len(masscan_array); i++ {
			r := response{}
			r.Timestamp = t.Timestamp
			r.Id = t.Id
			r.Task_Id = t.Task_Id
			r.Owner = t.Owner
			r.Report_id = t.Report_id
			r.Task_type = t.task_type
			r.Tool = t.Tool
			r.Ip = masscan_array[i].Ip
			r.Ports = masscan_array[i].Ports
			resp = append(resp, r)
		}

		resp_json, err := json.Marshal(resp)
		if err != nil {
			fmt.Println(err)
		}

		return string(resp_json)
	}

	return buf.String()

}
