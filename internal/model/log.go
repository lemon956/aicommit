package model

import (
	"net/http"
)

// logRequest is intentionally a no-op.
//
// Previous versions printed request details for debugging, but this caused
// unwanted output. All providers still call this function, so keeping it
// as a no-op disables request logging globally without touching each model.
func logRequest(req *http.Request, body []byte) {
}
