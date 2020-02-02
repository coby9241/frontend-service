package users_test

import (
	"crypto/rsa"
	"io/ioutil"
	"testing"

	. "github.com/coby9241/frontend-service/internal/models/users"
	"github.com/dgrijalva/jwt-go"
)

func LoadRSAPublicKeyFromDisk(location string) *rsa.PublicKey {
	keyData, e := ioutil.ReadFile(location)
	if e != nil {
		panic(e.Error())
	}
	key, e := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if e != nil {
		panic(e.Error())
	}
	return key
}

func TestToken(t *testing.T) {
	cases := []struct {
		name    string
		jwtKey  string
		token   string
		wantErr bool
	}{
		{
			name:   "test get claim success",
			jwtKey: LoadRSAPublicKeyFromDisk("../../../tests/testdata/sample_key.pub").N.String(),
			token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIifQ.8SyT66ZtDj-DgYSrvoMMoOO_jRl5tfGnAlHDFM2Cjkg",
		},
		{
			name:    "test get claim with with failure",
			token:   "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9eyJmb28iOiJiYXIifQ.EhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
			wantErr: true,
		},
		{
			name:    "test issue jwt token failure",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjF9.ymPzikm7vlQW4CAoEg_NzdbelV8aEDeRYa8x_0jk6uY",
			wantErr: true,
		},
	}

	for _, tt := range cases {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			if _, err := GetClaims(tc.token); (err != nil) != tc.wantErr {
				t.Errorf("failed to get claim: %v\n", err)
			}
		})
	}
}
