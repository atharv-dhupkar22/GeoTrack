package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"
)

func Timing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// WebSocket connections must bypass HTTP response buffering.
		if strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)

		elapsed := time.Since(start).Nanoseconds()

		for key, values := range recorder.Header() {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		body := recorder.Body.Bytes()
		contentType := recorder.Header().Get("Content-Type")

		if strings.Contains(contentType, "application/json") &&
			recorder.Code != http.StatusNoContent {

			var response map[string]any

			if err := json.Unmarshal(body, &response); err == nil && response != nil {
				response["time_ns"] = strconv.FormatInt(elapsed, 10)

				if updatedBody, err := json.Marshal(response); err == nil {
					body = updatedBody
				}
			}
		}

		w.WriteHeader(recorder.Code)

		if recorder.Code != http.StatusNoContent {
			_, _ = w.Write(bytes.Clone(body))
		}
	})
}
