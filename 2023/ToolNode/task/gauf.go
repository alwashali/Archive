package task

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/bitly/go-simplejson"
)

func gauf(t *task) string {

	//fmt.Println(t.Output)
	js := simplejson.New()

	js.Set("timestamp", t.Timestamp)
	js.Set("id", t.Id)
	js.Set("owner", t.Owner)
	js.Set("report_id", t.Report_id)
	js.Set("command", t.Command)
	js.Set("tool", t.Tool)
	js.Set("task_type", t.task_type)

	gau_urls := []string{}

	gau_scanner := bufio.NewScanner(&t.Output)
	gau_scanner.Split(bufio.ScanLines)

	for gau_scanner.Scan() {
		url := gau_scanner.Text()

		if strings.HasSuffix(url, ".js") {
			gau_urls = append(gau_urls, url)
		}
	}

	js.Set("url", gau_urls)

	out, err := js.MarshalJSON()
	if err != nil {
		fmt.Println(err)
	}
	return string(out)

}
