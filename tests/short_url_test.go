package tests

import (
	"net/http"
	"net/url"
	"os"
	"path"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"github.com/Alex322322/short-url/internal/http/server/handlers/url/save"
	"github.com/Alex322322/short-url/internal/lib/api"
	"github.com/Alex322322/short-url/internal/lib/random"
)

const (
	host = "localhost:8084"
)

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	// setup httpexpect client
	e := httpexpect.Default(t, u.String())

	err := godotenv.Load("../.env")
	require.NoError(t, err, "Failed to load .env file")

	password := os.Getenv("HTTP_SERVER_PASSWORD")
	require.NotEmpty(t, password, "HTTP_SERVER_PASSWORD must be set in .env file")

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.GenerateAlias(10),
		}).
		WithBasicAuth("admin", password).
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("alias")
}


func TestURLShortener_SaveRedirectRemove(t *testing.T) {
	err := godotenv.Load("../.env")
	require.NoError(t, err, "Failed to load .env file")

	password := os.Getenv("HTTP_SERVER_PASSWORD")
	require.NotEmpty(t, password, "HTTP_SERVER_PASSWORD must be set in .env file")

	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   "invalid_url",
			alias: gofakeit.Word(),
			error: "URL is not a valid URL",
		},
		{
			name:  "Empty Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			// Save
			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("admin", password).
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			// Redirect
			testRedirect(t, alias, tc.url)

			// Remove
			reqDel := e.DELETE("/"+path.Join("url", alias)).
				WithBasicAuth("admin", password).
				Expect().Status(http.StatusOK).
				JSON().Object()
			reqDel.Value("status").String().IsEqual("ok")

			// Redirect after remove
			testRedirectNotFound(t, alias)
		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "/" + alias,
	}

	redirectedToURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedToURL)
}

func testRedirectNotFound(t *testing.T, alias string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:  "/" + alias,
	}

	_, err := api.GetRedirect(u.String())
	require.Error(t, err, api.ErrInvalidStatusCode)
}
