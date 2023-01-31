package task

import (
	"bufio"
	"fmt"

	"github.com/bitly/go-simplejson"
)

func otherf(t *task) string {
	js := simplejson.New()

	r := bufio.NewScanner(&t.Output)

	js.Set("id", t.Id)
	js.Set("timestamp", t.Timestamp)
	js.Set("task_id", t.Task_Id)
	js.Set("owner", t.Owner)
	js.Set("report_id", t.Report_id)
	js.Set("command", t.Command)
	js.Set("tool", t.Tool)
	js.Set("task_type", t.task_type)
	js.Set("output", r.Text())

	js_resp, err := js.MarshalJSON()

	if err != nil {
		fmt.Println(err)
	}
	return string(js_resp)
}
