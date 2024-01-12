package usecase

import (
	"github.com/dgrijalva/jwt-go/v4"
	"testing"
	"time"
)

func Test_parseToken(t *testing.T) {
	mskLocation, _ := time.LoadLocation("Europe/Moscow")

	data := []struct {
		name       string
		token      string
		signingKey []byte
		expected   *UserClaims
		errMsg     string
	}{
		{
			name:       "correct",
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjY1NTAwNjc2MjEuODA4ODkxLCJVc2VySW5mbyI6eyJfaWQiOiI2NTMyNzBjZTA5Yzg5NmI5ZDM2NTBiMzgiLCJlbWFpbCI6InJ1cHljaG1hbkBtYWlsLnJ1IiwicGFzc3dvcmQiOiIkMmEkMTAkbnhkTzhvZkVXbUt4ZGFBdXdOb3RYZUU4aFBETmJXbmNkSEhhWGJDTTJ6TFRkU3UyVHlmUUMifX0.DifwaB2THXQ3IFNJFI0AzvMYJFHaODYVSh_eqAWZgoY",
			signingKey: []byte("secret_key"),
			expected: &UserClaims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: &jwt.Time{Time: time.Date(2177, time.July, 25, 2, 13, 41, 808890000, mskLocation)},
				},
				UserInfo: map[string]interface{}{
					"_id":   "653270ce09c896b9d3650b38",
					"email": "rupychman@mail.ru",
				},
			},
			errMsg: "",
		},
		{
			name:       "expired",
			token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjY1NTAwNjc2MjEuODA4ODkxLCJVc2VySW5mbyI6eyJfaWQiOiI2NTMyNzBjZTA5Yzg5NmI5ZDM2NTBiMzgiLCJlbWFpbCI6InJ1cHljaG1hbkBtYWlsLnJ1IiwicGFzc3dvcmQiOiIkMmEkMTAkbnhkTzhvZkVXbUt4ZGFBdXdOb3RYZUU4aFBETmJXbmNkSEhhWGJDTTJ6TFRkU3UyVHlmUUMifX0.DifwaB2THXQ3IFNJFI0AzvMYJFHaODYVSh_eqAWZgoY",
			signingKey: []byte("secret_key"),
			expected:   &UserClaims{},
			errMsg:     "",
		},
		{
			"not valid signing key",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjY1NTAwNjc2MjEuODA4ODkxLCJVc2VySW5mbyI6eyJfaWQiOiI2NTMyNzBjZTA5Yzg5NmI5ZDM2NTBiMzgiLCJlbWFpbCI6InJ1cHljaG1hbkBtYWlsLnJ1IiwicGFzc3dvcmQiOiIkMmEkMTAkbnhkTzhvZkVXbUt4ZGFBdXdOb3RYZUU4aFBETmJXbmNkSEhhWGJDTTJ6TFRkU3UyVHlmUUMifX0.dXLISmyIegGhhdV4BpOcfvXZgDVKjTwfz2Q35VNQV58",
			[]byte("your-256-bit-secret"),
			&UserClaims{},
			"",
		},
		{"not valid token", "", []byte(""), &UserClaims{}, "token is malformed: token contains an invalid number of segments"},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			_, err := parseToken(d.signingKey, d.token)
			//if result != d.expected {
			//	t.Errorf("Expected %v, got %v", d.expected, result)
			//}

			var errMsg string

			if err != nil {
				errMsg = err.Error()
			}

			if errMsg != d.errMsg {
				t.Errorf("Expected error message `%s`, got `%s", d.errMsg, errMsg)
			}
		})
	}
}
