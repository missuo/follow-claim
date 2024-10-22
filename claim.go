/*
 * @Author: Vincent Yang
 * @Date: 2024-09-30 02:01:59
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2024-10-12 01:43:24
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
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

func signFollow(cookie string, barkURL string, barkEnable bool, telegramEnable bool, telegramBotToken string, telegramChatID string) string {
	retryDelay := time.Minute
	attempt := 0

	for {
		result, err := attemptSignFollow(cookie, barkURL, barkEnable, telegramEnable, telegramBotToken, telegramChatID)

		if err == nil {
			return result // Success
		}

		if strings.Contains(result, "Already claimed") {
			return result // Stop retrying if "Already claimed"
		}

		attempt++
		message := fmt.Sprintf("Follow: Retry attempt %d after error: %v. Will retry in 1 minute.", attempt, err)
		sendToBark(message, barkURL, barkEnable)
		sendToTelegram(message, telegramEnable, telegramBotToken, telegramChatID)

		time.Sleep(retryDelay)
	}
}

func attemptSignFollow(cookie string, barkURL string, barkEnable bool, telegramEnable bool, telegramBotToken string, telegramChatID string) (string, error) {
	url := "https://api.follow.is/wallets/transactions/claim_daily"
	payload := map[string]string{"csrfToken": extractCSRFToken(cookie)}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "Follow: Error", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "Follow: Error", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.38(0x1800262c) NetType/4G Language/zh_CN")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "Follow: Error", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		message := fmt.Sprintf("Claim Points Failed: %v", result["message"])
		sendToBark(message, barkURL, barkEnable)
		sendToTelegram(message, telegramEnable, telegramBotToken, telegramChatID)
		return message, fmt.Errorf(message)
	}

	successMessage := "Claim Points Success"
	sendToBark(successMessage, barkURL, barkEnable)
	sendToTelegram(successMessage, telegramEnable, telegramBotToken, telegramChatID)
	return successMessage, nil
}

func main() {
	cookiesStr := os.Getenv("COOKIE")
	barkURL := os.Getenv("BARK_URL")
	scheduledTime := os.Getenv("SCHEDULED_TIME")
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")

	if cookiesStr == "" {
		log.Fatal("COOKIE must be set in environment variables")
	}

	var cookies []string
	if strings.Contains(cookiesStr, ",") {
		cookies = strings.Split(cookiesStr, ",")
	} else {
		cookies = []string{cookiesStr}
	}

	if scheduledTime == "" {
		scheduledTime = "00:05"
	}

	hour, minute, err := parseTime(scheduledTime)
	if err != nil {
		log.Fatalf("Invalid SCHEDULED_TIME format: %v", err)
	}

	barkEnable := barkURL != ""
	telegramEnable := telegramBotToken != "" && telegramChatID != ""
	c := cron.New(cron.WithLocation(time.UTC))
	_, err = c.AddFunc(fmt.Sprintf("%s %s * * *", minute, hour), func() {
		for i, cookie := range cookies {
			result := signFollow(strings.TrimSpace(cookie), barkURL, barkEnable, telegramEnable, telegramBotToken, telegramChatID)
			fmt.Printf("Account %d: %s\n", i+1, result)
		}
	})
	if err != nil {
		log.Fatal("Error scheduling task: ", err)
	}

	c.Start()
	fmt.Printf("Scheduler started. Will run daily at %s:%s UTC for %d accounts.\n", hour, minute, len(cookies))
	select {}
}
