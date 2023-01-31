package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alwashali/ToolNode/task"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

func run(resp http.ResponseWriter, r *http.Request) {
	//wg := waiting
	command := r.PostFormValue("command")
	report_id := r.PostFormValue("report_id")
	owner := r.PostFormValue("owner")
	task_type := r.PostFormValue("type")
	id := r.PostFormValue("id")
	return_webhook := r.PostFormValue("return_webhook")

	if command == "" || report_id == "" || owner == "" || id == "" {
		json.NewEncoder(resp).Encode(map[string]string{"error": "one or more fields are empty"})
		return
	}

	task := task.New_task(id, owner, report_id, command, task_type, return_webhook)

	//response
	json := simplejson.New()
	json.Set("id", task.Id)
	json.Set("task_id", task.Task_Id)
	json.Set("timestamp", task.Timestamp)
	json.Set("command", command)

	response, err := json.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	go task.Execute()

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(response)
}

func kill(resp http.ResponseWriter, r *http.Request) {
	// if called will terminate the tool task
	go func() {
		time.Sleep(time.Second * 2)
		os.Exit(1)
	}()
	json.NewEncoder(resp).Encode(map[string]string{"Status": "killed Successfully"})

}

func kill_task(resp http.ResponseWriter, r *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	id := r.PostFormValue("id")
	if id == "" {
		json.NewEncoder(resp).Encode(map[string]string{"error": "task ID not found in request"})
		return
	}

	t, err := task.Get_task(id)
	fmt.Println("killing", t.Task_Id)

	if err != nil {
		json.NewEncoder(resp).Encode(map[string]string{"Status": "Task not found"})
		return
	}

	t.TerminateTask()
	json.NewEncoder(resp).Encode(map[string]string{"Status": "Task terminated"})

}

func status(resp http.ResponseWriter, r *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	id := r.PostFormValue("id")
	if id == "" {
		json.NewEncoder(resp).Encode(map[string]string{"error": "task ID not found in request"})
		return
	}

	t, err := task.Get_task(id)

	if err != nil {
		json.NewEncoder(resp).Encode(map[string]string{"Status": "Task Not found"})
		return
	}

	fmt.Println("status:", t.Status)
	json := simplejson.New()
	json.Set("id", t.Task_Id)
	json.Set("timestamp", t.Timestamp)
	json.Set("status", t.Status)

	payload, err := json.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	resp.Write(payload)

}

func output(resp http.ResponseWriter, r *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	id := r.PostFormValue("id")
	if id == "" {
		json.NewEncoder(resp).Encode(map[string]string{"error": "task ID not found in request"})
		return
	}

	t, err := task.Get_task(id)
	if err != nil {
		json.NewEncoder(resp).Encode(map[string]string{"Error": err.Error()})
		return
	}

	Json_outpu := t.Formated_output()

	resp.Write([]byte(Json_outpu))
}

// perform couple of health checks
func health(resp http.ResponseWriter, r *http.Request) {
	json.NewEncoder(resp).Encode(map[string]string{"Status": "healthy"})
}

func tasks(resp http.ResponseWriter, r *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	tasks := task.Tasks()
	resp_bytes, err := json.Marshal(tasks)

	if err != nil {
		fmt.Println(err)
	}
	resp.Write(resp_bytes)

}

func main() {

	r := mux.NewRouter()

	// to run a command
	r.HandleFunc("/run", run).Methods("POST")
	// kill a task
	r.HandleFunc("/halt", kill).Methods("GET")
	// status of a task
	r.HandleFunc("/status", status).Methods("POST")
	//output
	r.HandleFunc("/output", output).Methods("POST")
	r.HandleFunc("/kill", kill_task).Methods("POST")
	r.HandleFunc("/tasks", tasks).Methods("GET")
	r.HandleFunc("/health", health).Methods("GET")

	log.Println("Listening 8080...")
	http.ListenAndServe(":8080", r)
}
