package auth

import (
	"gopkg.in/square/go-jose.v2/jwt"
)

// Passport is a collection used to identify whether an identity is valid and store identity information,
// such as account number, account name, tenant status, timeout period, group, role, operable resources, etc.
type Passport struct {
	AccountID     string   `json:"a"`
	AccountName   string   `json:"b"`
	AccountType   int      `json:"c"` // 0-tenantAccount,1-subAccount,2-systemAccount
	TenantID      string   `json:"d"`
	TenantCode    string   `json:"e"`
	DefaultExpire int      `json:"f"`
	Groups        []string `json:"g"`
	Roles         []string `json:"h"`
	Resource      Resource `json:"i"`
}

// Resource sets a list of all resources that the role can operate
type Resource struct {
	Workspaces   []string `json:"j"`
	Environments []string `json:"k"`
	Programs     []string `json:"l"`
	Projects     []string `json:"m"`
	Resources    []string `json:"n"`
}

// ClaimSet is a child struct of jwt.Claims, for store Passport in JWT
type ClaimSet struct {
	jwt.Claims
	Passport *Passport `json:"passport"`
}

//UMVerifier is a global variable to store JwtVerifier in memory.
var UMVerifier *JwtVerifier

const (
	defaultSslDir = "./ssl"
)

//InitJwt creates a new UMVerifier, using defaultSslDir is the ssl directory is not specified
func InitJwt(sslDir ...string) {
	dir := defaultSslDir
	if len(sslDir) > 0 {
		dir = sslDir[0]
	}
	verifier, err := NewJwtVerifier(dir)
	if err != nil {
		panic("NewJwtSigner fail,as:" + err.Error())
	} else {
		UMVerifier = verifier
	}
}
