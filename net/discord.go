package net

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	log "github.com/wlbr/commons/log"
)

type DiscordSender struct {
	queue   chan string
	webhook string
}

func NewDiscordSender(webhook string) *DiscordSender {
	d := &DiscordSender{queue: make(chan string, 10000), webhook: webhook}
	if webhook != "" {
		go d.loop()
	}
	return d
}

func (d *DiscordSender) loop() {
	for {
		m := <-d.queue
		t := d.pushMessageToDiscord(m)
		if t > 0 {
			time.Sleep(time.Duration(t) * time.Second)
		}
	}
}

func (d *DiscordSender) pushMessageToDiscord(message string) (timeout int) {
	payload := new(bytes.Buffer)
	m := fmt.Sprintf("{\"username\": \"CS:GO\", \"content\": \"%s\"}", message)

	payload.WriteString(m)
	resp, err := http.Post(d.webhook, "application/json", payload)
	if err != nil {
		log.Error("Error sending to discord: '%s' - %v", m, err)
		return -1
	}
	defer resp.Body.Close()

	remaining, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Limit-Remaining"))
	timeout, _ = strconv.Atoi(resp.Header.Get("X-RateLimit-Reset-After"))
	log.Info("Discord rate management:  Remaining: %d   Reset-After: %d - %s", remaining, timeout, resp.Header.Get("X-RateLimit-Reset-After"))
	if remaining > 0 {
		timeout = 0
	}

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusTooManyRequests:
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error("Error %d writing to Discord %s", resp.StatusCode, responseBody)
			return -1
		}
		var result map[string]int
		json.Unmarshal([]byte(responseBody), &result)

		timeout = result["retry_after"]
		break
	default:
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error("Error %d writing to Discord %s", resp.StatusCode, responseBody)
		}
		log.Info("Discord StatusCode %d:   Message: \"%s\"  Response: %s", resp.StatusCode, m, responseBody)
	}
	return timeout
}

func (d *DiscordSender) Send(message string) {
	d.queue <- message
}
