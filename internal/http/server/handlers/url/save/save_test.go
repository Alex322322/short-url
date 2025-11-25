package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Alex322322/short-url/internal/http/server/handlers/url/save"
	"github.com/Alex322322/short-url/internal/lib/logger/discardslog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	//"github.com/Alex322322/short-url/internal/http/server/handlers/url/save"
)

func TestSaveHandler(t *testing.T) {
	fixedTime := time.Now().Truncate(time.Second) // обрезаем до секунд для стабильности

	cases := []struct {
		name      string
		alias     string
		url       string
		timestamp time.Time
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
			timestamp: fixedTime,
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://google.com",
			timestamp: fixedTime,
		},
		{
			name:      "Empty URL",
			url:       "",
			alias:     "some_alias",
			respError: "URL is a required field",
			timestamp: time.Now(),
		},
		{
			name:      "Invalid URL",
			url:       "some invalid URL",
			alias:     "some_alias",
			respError: "URL is not a valid URL",
			timestamp: time.Now(),
		},
		{
			name:      "SaveURL Error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "internal server error",
			mockError: errors.New("unexpected error"),
			timestamp: time.Now(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := save.NewMockURLSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string"), fixedTime).
					Return(int64(1), tc.mockError).
					Once()
			}

			handler := save.New(discardslog.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s", "timestamp": "%v"}`, tc.url, tc.alias, fixedTime.Format(time.RFC3339))

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
