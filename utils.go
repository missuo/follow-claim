/*
 * @Author: Vincent Yang
 * @Date: 2024-09-30 02:02:45
 * @LastEditors: Vincent Young
 * @LastEditTime: 2024-09-30 10:25:45
 * @FilePath: /follow-claim/utils.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2024 by Vincent, All Rights Reserved.
 */
package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// extractCSRFToken extracts the value of authjs.csrf-token from a cookie string
func extractCSRFToken(cookieString string) string {
	// Split the cookie string by semicolons
	cookies := strings.Split(cookieString, ";")

	// Iterate through all cookies
	for _, cookie := range cookies {
		// Trim leading and trailing whitespace
		cookie = strings.TrimSpace(cookie)

		// Check if this is the csrf-token
		if strings.HasPrefix(cookie, "authjs.csrf-token=") {
			// Extract and return the value
			return strings.TrimPrefix(cookie, "authjs.csrf-token=")
		}
	}

	// Return an empty string if not found
	return ""
}

// SendToBark sends a message to the Bark service
func sendToBark(message string, barkURL string, barkEnable bool) {
	if !barkEnable {
		return
	}
	// Encode the message to ensure it's URL-safe
	encodedMessage := url.PathEscape(message)

	// Construct the full URL
	fullURL := fmt.Sprintf("%s/Follow Claim/%s?icon=https://dc.missuo.ru/file/1290196428324601907&group=Follow", barkURL, encodedMessage)

	// Send the GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Printf("Failed to send notification: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to send notification: HTTP status %d\n", resp.StatusCode)
		return
	}

	fmt.Printf("Successfully sent notification: %s\n", message)
}

func parseTime(timeStr string) (hour, minute string, err error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("time should be in HH:MM format")
	}

	h, err := strconv.Atoi(parts[0])
	if err != nil || h < 0 || h > 23 {
		return "", "", fmt.Errorf("invalid hour")
	}

	m, err := strconv.Atoi(parts[1])
	if err != nil || m < 0 || m > 59 {
		return "", "", fmt.Errorf("invalid minute")
	}

	return fmt.Sprintf("%02d", h), fmt.Sprintf("%02d", m), nil
}
