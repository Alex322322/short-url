package redirect_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Alex322322/short-url/internal/http/server/handlers/url/redirect"
	"github.com/Alex322322/short-url/internal/lib/logger/discardslog"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
	fixedTime := time.Now().Truncate(time.Second) // обрезаем до секунд для стабильности
	cases := []struct {
		name         string
		alias        string
		url          string
		timestamp    time.Time
		respError    string
		mockError    error
		expectedCode int
		expectedURL  string
	}{
		{
			name:         "Success redirect",
			alias:        "test_alias",
			url:          "https://google.com",
			timestamp:    fixedTime,
			expectedCode: http.StatusFound,
			expectedURL:  "https://google.com",
		},
		{
			name:         "Empty alias",
			alias:        "",
			url:          "https://google.com",
			expectedCode: http.StatusBadRequest,
			timestamp:    fixedTime,
		},
		{
			name:         "Empty URL",
			url:          "",
			alias:        "some_alias",
			expectedCode: http.StatusBadRequest,
			timestamp:    fixedTime,
		},
		{
			name:         "Invalid URL",
			url:          "some invalid URL",
			alias:        "some_alias",
			respError:    "URL is not a valid URL",
			expectedCode: http.StatusBadRequest,
			timestamp:    fixedTime,
		},
		{
			name:         "URL Not Found",
			alias:        "non_existent_alias",
			respError:    "URL not found",
			mockError:    errors.New("URL not found"),
			timestamp:    fixedTime,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "GetURL Error",
			alias:        "test_alias",
			respError:    "internal server error",
			mockError:    errors.New("unexpected error"),
			timestamp:    fixedTime,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlGetterMock := redirect.NewMockURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).
					Once()
			}

			handler := redirect.New(discardslog.NewDiscardLogger(), urlGetterMock)

			req, err := http.NewRequest(http.MethodGet, "/"+tc.alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.expectedCode)

			if tc.expectedCode == http.StatusFound {
				require.Equal(t, tc.expectedURL, rr.Header().Get("Location"))
			}
			if tc.respError != "" {
				var resp struct {
					Error string `json:"error"`
				}
				body := rr.Body.String()
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				require.Equal(t, tc.respError, resp.Error)
			}
		})
	}
}
