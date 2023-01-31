package task

//nmap -p80,443 -oX - --open avleonov.com
import (
	"bytes"
	"fmt"

	xj "github.com/basgys/goxml2json"
	"github.com/bitly/go-simplejson"
)

func nmapf(t *task) string {

	js := simplejson.New()
	js.Set("timestamp", t.Timestamp)
	js.Set("id", t.Id)
	js.Set("owner", t.Owner)
	js.Set("report_id", t.Report_id)
	js.Set("command", t.Command)
	js.Set("tool", t.Tool)
	js.Set("task_type", t.task_type)

	buf := bytes.NewBuffer(t.Output.Bytes())

	// nmap output in xml to json string
	json, err := xj.Convert(buf)
	if err != nil {
		fmt.Println(err)
	}

	//json string to json object
	sj, err := simplejson.NewJson(json.Bytes())
	if err != nil {
		fmt.Println(err)
	}

	address := sj.GetPath("nmaprun", "hosthint", "address", "-addr")
	addr, err := address.String()
	fmt.Println("address nmap", addr)
	if err != nil {
		fmt.Println(err)
	}
	js.Set("address", addr)

	js_ports := sj.GetPath("nmaprun", "host", "ports", "port")

	ps, _ := js_ports.EncodePretty()
	fmt.Println("ports:")

	js.Set("address", addr)
	js.Set("ports", string(ps))

	output, err := js.Encode()
	if err != nil {
		fmt.Println(err)
	}
	return string(output)
}
