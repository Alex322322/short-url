package redirect_test

import (
	"net/http/httptest"
	"testing"

	"github.com/Alex322322/short-url/internal/http/server/handlers/url/redirect"
	"github.com/Alex322322/short-url/internal/lib/api"
	"github.com/Alex322322/short-url/internal/lib/logger/discardslog"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://www.google.com/",
		},
		/*
			{
				name:      "Empty alias",
				alias:     "",
				url:       "https://www.google.com/",
				respError: "Not Found",
			},
			{
				name:      "Empty URL",
				url:       "",
				alias:     "some_alias",
				respError: "URL is a required field",
				mockError: errors.New("unexpected error"),
			},
			{
				name:      "Invalid URL",
				url:       "some invalid URL",
				alias:     "some_alias",
				respError: "URL is not a valid URL",
			},
		*/
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			//t.Parallel()

			urlGetterMock := redirect.NewMockURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).Once()
			}

			handler := redirect.New(discardslog.NewDiscardLogger(), urlGetterMock)

			r := chi.NewRouter()
			r.Get("/{alias}", handler)

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			require.NoError(t, err)

			// check the final URL after redirection
			assert.Equal(t, tc.url, redirectedToURL)
		})
	}
}
