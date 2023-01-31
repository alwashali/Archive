package task

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/alwashali/ToolNode/webhook"
	"github.com/rs/xid"
)

type task struct {
	Id             string
	Task_Id        string
	Timestamp      string
	Owner          string
	Report_id      string
	Command        string
	Tool           string
	task_type      string
	Status         string
	Output         bytes.Buffer
	Return_webhook string
}

var tasks []*task

func New_task(id, owner, report_id, command, task_type, return_webhook string) *task {
	t := task{
		Id:             id,
		Owner:          owner,
		Report_id:      report_id,
		Command:        command,
		task_type:      task_type,
		Return_webhook: return_webhook,
	}

	time := time.Now()
	t.Timestamp = fmt.Sprintf("%d/%02d/%02d %02d:%02d:%02d",
		time.Year(), time.Month(), time.Day(),
		time.Hour(), time.Minute(), time.Second())

	t.Task_Id = xid.New().String()
	t.Tool = strings.Split(command, " ")[0]
	t.Status = "Queued"
	tasks = append(tasks, &t)
	return &t
}

func (t *task) TerminateTask() {

	if t.Status == "running" {
		t.Status = "terminated"
		msg := fmt.Sprintf("Task %s was terminated", t.Task_Id)

		panic(msg)
	}

}

// Prepare output json based on each tool output
func (t *task) Formated_output() string {

	if t.Status == "Running" {
		return "Task Has not finished"

	} else if t.Status == "Queued" {
		return "Task Has not started"

	} else if t.Status == "Finished" {
		if t.Tool == "subfinder" {
			return subfinderf(t)

		} else if t.Tool == "masscan" {
			return masscanf(t)

		} else if strings.Contains(t.Command, "httpx") {
			return httpxf(t)

		} else if t.Tool == "ffuf" {
			return ffuff(t)

		} else if t.Tool == "nmap" {
			return nmapf(t)
		} else if strings.Contains(t.Command, "dnsx") {
			return dnsxf(t)
		} else if t.Tool == "gau" {
			return gauf(t)
		} else if strings.Contains(t.Command, "wpscan") {
			return wpscanf(t)
		}

	}
	return otherf(t)
}

// execute system command in a thread
func (t *task) Execute() {

	// replace two or more spaces with one
	m1 := regexp.MustCompile(`\s+`)
	command := m1.ReplaceAllString(t.Command, " ")

	cmd := exec.Command(strings.Split(command, " ")[0], strings.Split(command, " ")[1:]...)
	fmt.Println(strings.Split(t.Command, " ")[0])
	fmt.Println(strings.Split(t.Command, " ")[1:])

	var stderr bytes.Buffer

	cmd.Stderr = &stderr
	cmd.Stdout = &t.Output
	cmd.Stdin = nil

	t.Status = "Running"
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error starting the command", err, stderr.String())
		t.Status = "Error"
		return
	}

	t.Status = "Finished"
	fmt.Println(t.Formated_output())
	if t.Return_webhook != "" {
		webhook.DoPost([]byte(t.Formated_output()), t.Return_webhook)
	}

	webhook.DoPost([]byte(t.Formated_output()), "http://64.227.160.182:5678/webhook-test/push")
	webhook.DoPost([]byte(t.Formated_output()), "http://64.227.160.182:5678/webhook/push")
}

func Get_task(id string) (*task, error) {
	for _, t := range tasks {

		if t.Task_Id == id {
			return t, nil
		}
	}
	return nil, errors.New("task not found")
}

func Tasks() []*task {
	return tasks
}
