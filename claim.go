/*
 * @Author: Vincent Yang
 * @Date: 2024-09-30 02:01:59
 * @LastEditors: Vincent Young
 * @LastEditTime: 2024-09-30 10:25:36
 * @FilePath: /follow-claim/claim.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2024 by Vincent, All Rights Reserved.
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

func signFollow(cookie string, barkURL string, barkEnable bool) string {
	url := "https://api.follow.is/wallets/transactions/claim_daily"
	payload := map[string]string{"csrfToken": extractCSRFToken(cookie)}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		sendToBark("Follow: Error: "+err.Error(), barkURL, barkEnable)
		return "Follow: Error"
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		sendToBark("Error: "+err.Error(), barkURL, barkEnable)
		return "Error"
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.38(0x1800262c) NetType/4G Language/zh_CN")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		sendToBark("Error: "+err.Error(), barkURL, barkEnable)
		return "Error"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		message := fmt.Sprintf("Claim Points Failed: %v", result["message"])
		sendToBark(message, barkURL, barkEnable)
		return message
	}

	sendToBark("Claim Points Success", barkURL, barkEnable)
	return "Claim Points Success"
}

func main() {
	cookie := os.Getenv("COOKIE")
	barkURL := os.Getenv("BARK_URL")
	scheduledTime := os.Getenv("SCHEDULED_TIME")

	if cookie == "" {
		log.Fatal("COOKIE must be set in environment variables")
	}

	if scheduledTime == "" {
		scheduledTime = "00:05"
	}

	hour, minute, err := parseTime(scheduledTime)
	if err != nil {
		log.Fatalf("Invalid SCHEDULED_TIME format: %v", err)
	}

	barkEnable := barkURL != ""

	c := cron.New(cron.WithLocation(time.UTC))
	_, err = c.AddFunc(fmt.Sprintf("%s %s * * *", minute, hour), func() {
		result := signFollow(cookie, barkURL, barkEnable)
		fmt.Println(result)
	})
	if err != nil {
		log.Fatal("Error scheduling task: ", err)
	}

	c.Start()
	fmt.Printf("Scheduler started. Will run daily at %s:%s UTC.\n", hour, minute)
	select {}
}
