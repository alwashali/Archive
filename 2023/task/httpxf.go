package task

import (
	"fmt"

	"github.com/bitly/go-simplejson"
)

func httpxf(t *task) string {

	//fmt.Println(t.Output)
	js := simplejson.New()

	js.Set("timestamp", t.Timestamp)
	js.Set("id", t.Id)
	js.Set("owner", t.Owner)
	js.Set("report_id", t.Report_id)
	js.Set("command", t.Command)
	js.Set("tool", t.Tool)
	js.Set("task_type", t.task_type)

	httpxJS, err := simplejson.NewJson(t.Output.Bytes())
	if err != nil {
		fmt.Println("Error decoding httpx json", err)
	}
	port, err := httpxJS.Get("port").String()
	if err != nil {
		fmt.Println(err)
	}
	url, err := httpxJS.Get("url").String()
	if err != nil {
		fmt.Println(err)
	}
	title, err := httpxJS.Get("title").String()
	if err != nil {
		fmt.Println(err)
	}
	scheme, err := httpxJS.Get("scheme").String()
	if err != nil {
		fmt.Println(err)
	}
	webserver, err := httpxJS.Get("webserver").String()
	if err != nil {
		fmt.Println(err)
	}
	method, err := httpxJS.Get("method").String()
	if err != nil {
		fmt.Println(err)
	}
	host, err := httpxJS.Get("host").String()
	if err != nil {
		fmt.Println(err)
	}
	status_code, err := httpxJS.Get("status_code").Int()
	if err != nil {
		fmt.Println(err)
	}

	content_length, err := httpxJS.Get("content_length").Int64()
	if err != nil {
		fmt.Println(err)
	}

	a, err := httpxJS.Get("a").StringArray()
	if err != nil {
		fmt.Println(err)
	}

	raw_headers, err := httpxJS.Get("raw_header").String()
	if err != nil {
		fmt.Println(err)
	}

	body, err := httpxJS.Get("body").String()
	if err != nil {
		fmt.Println(err)
	}

	js.Set("url", url)
	js.Set("port", port)
	js.Set("title", title)
	js.Set("webserver", webserver)
	js.Set("scheme", scheme)
	js.Set("method", method)
	js.Set("host", host)
	js.Set("status_code", status_code)
	js.Set("content_length", content_length)
	js.Set("a", a)
	js.Set("raw_headers", raw_headers)
	js.Set("body", body)
	out, err := js.MarshalJSON()
	if err != nil {
		fmt.Println(err)
	}
	return string(out)

}
