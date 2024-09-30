/*
 * @Author: Vincent Yang
 * @Date: 2024-09-30 02:02:45
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2024-09-30 02:23:42
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
