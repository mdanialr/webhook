package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModel_Sanitization(t *testing.T) {
	testCases := []struct {
		name    string
		sample  Model
		expect  Model
		wantErr bool
	}{
		{
			name:    "Should error when `user` field not provided",
			sample:  Model{},
			wantErr: true,
		},
		{
			name:    "Should error when `pass` field not provided",
			sample:  Model{User: "nzk"},
			wantErr: true,
		},
		{
			name:    "Should error when `repo` field not provided",
			sample:  Model{User: "nzk", Pass: "secret"},
			wantErr: true,
		},
		{
			name:   "Should pass when all required fields are provided",
			sample: Model{User: "nzk", Pass: "secret", Repo: "private"},
			expect: Model{User: "nzk", Pass: "secret", Repo: "private", Tag: "latest"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.sample.Sanitization()

			switch tc.wantErr {
			case true:
				require.Error(t, err)
			case false:
				require.NoError(t, err)
				assert.Equal(t, tc.expect.User, tc.sample.User)
				assert.Equal(t, tc.expect.Pass, tc.sample.Pass)
				assert.Equal(t, tc.expect.Repo, tc.sample.Repo)
			}
		})
	}
}

func TestModel_SanitizationOptionalFields(t *testing.T) {
	testCases := []struct {
		name    string
		sample  Model
		expect  Model
		wantErr bool
	}{
		{
			name:   "Default value for `tag` field should be 'latest' when not provided",
			sample: Model{User: "nzk", Pass: "secret", Repo: "private"},
			expect: Model{User: "nzk", Pass: "secret", Repo: "private", Tag: "latest", Image: "nzk/private:latest", Name: "private_latest", Id: "nzk_private_latest"},
		},
		{
			name:   "Should use `tag` from config file when provided instead the default one",
			sample: Model{User: "nzk", Pass: "secret", Repo: "private", Tag: "test-case"},
			expect: Model{User: "nzk", Pass: "secret", Repo: "private", Tag: "test-case", Image: "nzk/private:test-case", Name: "private_test-case", Id: "nzk_private_test-case"},
		},
		{
			name:   "`Id` field should be formatted as expected",
			sample: Model{User: "hi", Pass: "secret", Repo: "from", Tag: "earth"},
			expect: Model{User: "hi", Pass: "secret", Repo: "from", Tag: "earth", Image: "hi/from:earth", Name: "from_earth", Id: "hi_from_earth"},
		},
		{
			name:   "`Image` field should be formatted as expected",
			sample: Model{User: "nzk", Pass: "secret", Repo: "private"},
			expect: Model{User: "nzk", Pass: "secret", Repo: "private", Tag: "latest", Image: "nzk/private:latest", Name: "private_latest", Id: "nzk_private_latest"},
		},
		{
			name:   "`Name` field should be formatted as expected",
			sample: Model{User: "nzk", Pass: "secret", Repo: "private"},
			expect: Model{User: "nzk", Pass: "secret", Repo: "private", Tag: "latest", Image: "nzk/private:latest", Name: "private_latest", Id: "nzk_private_latest"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.sample.Sanitization()

			switch tc.wantErr {
			case true:
				require.Error(t, err)
			case false:
				require.NoError(t, err)
				assert.Equal(t, tc.expect, tc.sample)
			}
		})
	}
}

func TestModel_ParsePullCommand(t *testing.T) {
	testCases := []struct {
		name   string
		sample Model
		expect string
	}{
		{
			name:   "Should result as expected with minimal required fields are provided",
			sample: Model{User: "nzk", Pass: "secret", Repo: "private"},
			expect: "docker login -u nzk -p secret && " +
				"docker pull -q nzk/private:latest && " +
				"docker stop private_latest && " +
				"docker container prune -f && " +
				"docker run --name private_latest -d  nzk/private:latest",
		},
		{
			name:   "Should result as expected with minimal required fields are provided and args also provided",
			sample: Model{User: "nzk", Pass: "secret", Repo: "private", Args: "-p 1000:1000"},
			expect: "docker login -u nzk -p secret && " +
				"docker pull -q nzk/private:latest && " +
				"docker stop private_latest && " +
				"docker container prune -f && " +
				"docker run --name private_latest -d -p 1000:1000 nzk/private:latest",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, tc.sample.Sanitization())
			out := tc.sample.ParsePullCommand()

			assert.Equal(t, tc.expect, out)
		})
	}
}
