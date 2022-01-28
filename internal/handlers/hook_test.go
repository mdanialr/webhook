package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mdanialr/webhook/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fakeConfigFile = `
env: prod
port: 5050
secret: secret
log: /var/log/webhook/log
service:
  - repo:
      name: fiber-ln
      root: /home/nzk/dir/Fiber/light_novel/
      opt_cmd: "go build -o bin/fiber-ln main.go && systemctl restart fiber-ln"
  - repo:
      name: cd_test
      root: /home/nzk/dir/Laravel/cd_test/
      opt_cmd: pwd
`

var sampleRequestBody = []string{
	`
{
  "ref": "refs/heads/main",
  "commits": [
    {
      "message": ":test your ass",
      "committer": {
        "username": "nzk"
      }
    }
  ]
}
`, `
{
  "ref": "refs/heads/master",
  "commits": [
    {
      "message": ":testing",
      "committer": {
        "username": "nzk"
      }
    }
  ]
}
`,
}

type responseJSON struct {
	Committer string `json:"committer"`
	Message   string `json:"message"`
	Branch    string `json:"branch"`
	Reload    bool   `json:"reload"`
	Repo      string `json:"repo"`
}

func TestHook_SimpleTest(t *testing.T) {
	testCases := []struct {
		name             string
		route            string
		method           string
		expectedResponse responseJSON
		expectedCode     int
		expectedMIME     string
		isErr            bool
	}{
		{
			name: "POST /hook/x-repo : Should return 200 and JSON response that" +
				"match with the params 'x-repo'",
			route:            "/hook/x-repo",
			method:           fiber.MethodPost,
			expectedResponse: responseJSON{Repo: "x-repo"},
			expectedCode:     200,
			expectedMIME:     fiber.MIMEApplicationJSON,
		},
		{
			name:         "POST /hook/ : Route w/o params should not found and return 404",
			route:        "/hook/",
			method:       fiber.MethodPost,
			expectedCode: 404,
			expectedMIME: fiber.MIMEApplicationForm,
			isErr:        true,
		},
		{
			name:         "POST /err : Route Not found Should return 404",
			route:        "/err",
			method:       fiber.MethodPost,
			expectedCode: 404,
			expectedMIME: fiber.MIMEApplicationForm,
			isErr:        true,
		},
	}

	buf := bytes.NewBufferString(fakeConfigFile)

	mod, err := config.NewConfig(buf)
	require.NoError(t, err)

	app := fiber.New()
	app.Post("/hook/:repo", Hook(mod))

	for i, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.isErr {
			case false:
				buf := bytes.NewBufferString(sampleRequestBody[0])
				req := httptest.NewRequest(tt.method, tt.route, buf)
				req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
				res, _ := app.Test(req)
				assert.Equal(t, tt.expectedCode, res.StatusCode)
				assert.Equal(t, tt.expectedMIME, res.Header.Get("Content-Type"))
				var rJSON responseJSON
				if err := json.NewDecoder(res.Body).Decode(&rJSON); err != nil {
					t.Fatalf("failed to decode json in test #%v: %v", i+1, err)
				}
				assert.Equal(t, tt.expectedResponse.Repo, rJSON.Repo)
			case true:
				req := httptest.NewRequest(tt.method, tt.route, nil)
				res, _ := app.Test(req)
				assert.Equal(t, tt.expectedCode, res.StatusCode)
				assert.NotEqual(t, tt.expectedMIME, res.Header.Get("Content-Type"))
			}

		})
	}

}

func TestHook_TriggerErrorTest(t *testing.T) {
	buf := bytes.NewBufferString(fakeConfigFile)

	mod, err := config.NewConfig(buf)
	require.NoError(t, err)

	app := fiber.New()
	app.Post("/hook/:repo", Hook(mod))

	const route = "/hook/:x-repo"
	const expectedStatusCode = fiber.StatusBadRequest
	const expectedMIME = fiber.MIMEApplicationJSON

	t.Run("1# Posting with anything else than JSON format should error and return 400", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, route, nil)
		res, _ := app.Test(req)
		assert.Equal(t, expectedStatusCode, res.StatusCode)
		assert.Equal(t, expectedMIME, res.Header.Get("Content-Type"))
	})

	t.Run("2# Posting with empty JSON request body should error and return 400", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, route, nil)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		res, _ := app.Test(req)
		assert.Equal(t, expectedStatusCode, res.StatusCode)
		assert.Equal(t, expectedMIME, res.Header.Get("Content-Type"))
	})
}

func TestHook_UsingBodyRequest(t *testing.T) {
	testCases := []struct {
		name             string
		route            string
		method           string
		configSample     config.Model
		bodyRequest      string
		expectedResponse responseJSON
		expectedCode     int
		expectedMIME     string
	}{
		{
			name: "1# This should make Reload true cause there is no need to" +
				"validate anything since both config username and keyword are empty",
			route:        "/hook/s-repo",
			method:       fiber.MethodPost,
			bodyRequest:  sampleRequestBody[1],
			configSample: config.Model{Secret: "1"},
			expectedCode: 200,
			expectedMIME: fiber.MIMEApplicationJSON,
			expectedResponse: responseJSON{
				Repo: "s-repo", Message: ":testing", Committer: "nzk", Branch: "master", Reload: true,
			},
		},
		{
			name: "2# This should make Reload false because config username does not " +
				"match w Committer",
			route:        "/hook/z-repo",
			method:       fiber.MethodPost,
			bodyRequest:  sampleRequestBody[0],
			configSample: config.Model{Secret: "1", Usr: "hasagi"},
			expectedCode: 200,
			expectedMIME: fiber.MIMEApplicationJSON,
			expectedResponse: responseJSON{
				Repo: "z-repo", Message: ":test your ass", Committer: "nzk", Branch: "main", Reload: false,
			},
		},
		{
			name: "3# This should make Reload false because config keyword does not" +
				"match w Message's prefix",
			route:        "/hook/z-repo",
			method:       fiber.MethodPost,
			bodyRequest:  sampleRequestBody[0],
			configSample: config.Model{Secret: "1", Keyword: "&"},
			expectedCode: 200,
			expectedMIME: fiber.MIMEApplicationJSON,
			expectedResponse: responseJSON{
				Repo: "z-repo", Message: ":test your ass", Committer: "nzk", Branch: "main", Reload: false,
			},
		},
		{
			name: "4# This should make Reload true because config username" +
				"does match w Committer",
			route:        "/hook/z-repo",
			method:       fiber.MethodPost,
			bodyRequest:  sampleRequestBody[0],
			configSample: config.Model{Secret: "1", Usr: "nzk"},
			expectedCode: 200,
			expectedMIME: fiber.MIMEApplicationJSON,
			expectedResponse: responseJSON{
				Repo: "z-repo", Message: ":test your ass", Committer: "nzk", Branch: "main", Reload: true,
			},
		},
		{
			name: "5# This should make Reload true because config keyword" +
				"does match w Message's prefix",
			route:        "/hook/z-repo",
			method:       fiber.MethodPost,
			bodyRequest:  sampleRequestBody[0],
			configSample: config.Model{Secret: "1", Keyword: ":"},
			expectedCode: 200,
			expectedMIME: fiber.MIMEApplicationJSON,
			expectedResponse: responseJSON{
				Repo: "z-repo", Message: ":test your ass", Committer: "nzk", Branch: "main", Reload: true,
			},
		},
		{
			name: "6# This should make Reload false because config username does not" +
				"match w Committer, even keyword is match w Message's prefix",
			route:        "/hook/z-repo",
			method:       fiber.MethodPost,
			bodyRequest:  sampleRequestBody[0],
			configSample: config.Model{Secret: "1", Keyword: ":", Usr: "hasagi"},
			expectedCode: 200,
			expectedMIME: fiber.MIMEApplicationJSON,
			expectedResponse: responseJSON{
				Repo: "z-repo", Message: ":test your ass", Committer: "nzk", Branch: "main", Reload: false,
			},
		},
		{
			name: "7# This should make Reload false because config keyword does not" +
				"match w Message's prefix, even username is match w Committer",
			route:        "/hook/z-repo",
			method:       fiber.MethodPost,
			bodyRequest:  sampleRequestBody[0],
			configSample: config.Model{Secret: "1", Keyword: "&", Usr: "nzk"},
			expectedCode: 200,
			expectedMIME: fiber.MIMEApplicationJSON,
			expectedResponse: responseJSON{
				Repo: "z-repo", Message: ":test your ass", Committer: "nzk", Branch: "main", Reload: false,
			},
		},
		{
			name: "8# This should make Reload false because config keyword does not" +
				"match w Message's prefix and username also does not match w Committer",
			route:        "/hook/z-repo",
			method:       fiber.MethodPost,
			bodyRequest:  sampleRequestBody[0],
			configSample: config.Model{Secret: "1", Keyword: "&", Usr: "nzk"},
			expectedCode: 200,
			expectedMIME: fiber.MIMEApplicationJSON,
			expectedResponse: responseJSON{
				Repo: "z-repo", Message: ":test your ass", Committer: "nzk", Branch: "main", Reload: false,
			},
		},
		{
			name: "8# This should make Reload true because config keyword does " +
				"match w Message's prefix and username also does match w Committer",
			route:        "/hook/z-repo",
			method:       fiber.MethodPost,
			bodyRequest:  sampleRequestBody[0],
			configSample: config.Model{Secret: "1", Keyword: ":", Usr: "nzk"},
			expectedCode: 200,
			expectedMIME: fiber.MIMEApplicationJSON,
			expectedResponse: responseJSON{
				Repo: "z-repo", Message: ":test your ass", Committer: "nzk", Branch: "main", Reload: true,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Post("/hook/:repo", Hook(&tt.configSample))

			// Validate config to fill in required username and keyword
			require.NoError(t, tt.configSample.Sanitization())

			req := httptest.NewRequest(tt.method, tt.route, bytes.NewBufferString(tt.bodyRequest))
			req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
			res, _ := app.Test(req)

			assert.Equal(t, tt.expectedCode, res.StatusCode)
			assert.Equal(t, tt.expectedMIME, res.Header.Get("Content-Type"))

			var rJSON responseJSON
			require.NoError(t,
				json.NewDecoder(res.Body).Decode(&rJSON),
				"failed to decode json in test",
			)

			assert.Equal(t, tt.expectedResponse, rJSON)
		})
	}
}
