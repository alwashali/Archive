package webhook

import (
	"bytes"
	"fmt"
	"net/http"
)

func DoPost(payload []byte, url string) {

	fmt.Println("sending output to N8N Push webook")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println("status code n8n: ", resp.StatusCode)
	if err != nil {
		fmt.Println("Couldn't send output to N8N", err)

	}
	defer resp.Body.Close()

}

func DoGet(payload []byte, url string) {

}
