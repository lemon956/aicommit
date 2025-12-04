package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// logRequest prints the details of the HTTP request to standard output.
// It masks sensitive header values like API keys.
func logRequest(req *http.Request, body []byte) {
	fmt.Printf("\n[Request Info]\n")
	fmt.Printf("URL: %s\n", req.URL.String())
	fmt.Printf("Method: %s\n", req.Method)

	fmt.Println("Headers:")
	for k, v := range req.Header {
		val := strings.Join(v, ", ")
		// Mask sensitive headers
		lowerK := strings.ToLower(k)
		if lowerK == "authorization" || lowerK == "x-api-key" || lowerK == "api-key" {
			// Show first few chars if long enough, or just [MASKED]
			if len(val) > 10 {
				val = val[:7] + "...[MASKED]"
			} else {
				val = "[MASKED]"
			}
		}
		fmt.Printf("  %s: %s\n", k, val)
	}

	if len(body) > 0 {
		fmt.Println("Body:")
		// Try to pretty print JSON
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, body, "  ", "  "); err == nil {
			fmt.Printf("  %s\n", prettyJSON.String())
		} else {
			// If not JSON or error, print raw string
			fmt.Printf("  %s\n", string(body))
		}
	}
	fmt.Println("----------------------------------------")
}
