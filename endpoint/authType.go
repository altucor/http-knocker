package endpoint

import (
	"errors"
	"net/http"
	"os"
	"strings"

	auth "github.com/abbot/go-http-auth"
	"github.com/altucor/http-knocker/logging"

	"golang.org/x/exp/slices"
)

type AuthType string

const (
	AUTH_TYPE_NONE       AuthType = "none"
	AUTH_TYPE_BASIC_AUTH AuthType = "basic-auth"
	AUTH_TYPE_AUTHELIA   AuthType = "authelia"
)

var (
	authTypeArr = []AuthType{
		AUTH_TYPE_NONE,
		AUTH_TYPE_BASIC_AUTH,
		AUTH_TYPE_AUTHELIA,
	}
)

func (ctx *AuthType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tempStr string = ""
	err := unmarshal(&tempStr)
	if err != nil {
		return err
	}
	tempAuthType := AuthType(tempStr)
	if !slices.Contains(authTypeArr, tempAuthType) {
		logging.CommonLog().Error("Cannot init from string")
		return errors.New("cannot init from string")
	}
	*ctx = tempAuthType
	return nil
}

type Auth struct {
	Type      AuthType `yaml:"auth-type"`
	UsersFile string   `yaml:"users-file"`
	Users     []string `yaml:"users"`
	authUsers map[string]string
}

func parseHtpasswdUserLine(line string) (string, string, error) {
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return "", "", errors.New("cannot parse htpasswd line")
	}
	return parts[0], parts[1], nil
}

func (ctx *Auth) setHtpasswdUsersFromArray(users []string) error {
	for _, line := range users {
		user, passHash, err := parseHtpasswdUserLine(line)
		if err != nil {
			logging.CommonLog().Error("cannot parse htpasswd line")
			return errors.New("cannot parse htpasswd line")
		}
		ctx.authUsers[user] = passHash
	}
	return nil
}

func (ctx *Auth) setHtpasswdUsersFromFile(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		logging.CommonLog().Errorf("cannot parse htpasswd file: %s\n", file)
		return errors.New("cannot parse htpasswd file")
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		user, passHash, err := parseHtpasswdUserLine(line)
		if err != nil {
			logging.CommonLog().Error("cannot parse htpasswd line")
			return errors.New("cannot parse htpasswd line")
		}
		ctx.authUsers[user] = passHash
	}
	return nil
}

func (ctx *Auth) basicAuthCheck(user string, realm string) string {
	passHash, ok := ctx.authUsers[user]
	if ok {
		return passHash
	}
	return ""
}

func (ctx *Auth) SetDefaults() {
	ctx.authUsers = make(map[string]string)
	if ctx.Type == "" {
		ctx.Type = AUTH_TYPE_NONE
	}
	switch ctx.Type {
	case AUTH_TYPE_BASIC_AUTH:
		if len(ctx.Users) != 0 {
			err := ctx.setHtpasswdUsersFromArray(ctx.Users)
			if err != nil {
				logging.CommonLog().Fatal("Cannot process htpassd users array")
			}
		}
		if ctx.UsersFile != "" {
			err := ctx.setHtpasswdUsersFromFile(
				ctx.UsersFile)
			if err != nil {
				logging.CommonLog().Fatal("Cannot process htpassd users file")
			}
		}
		if len(ctx.authUsers) == 0 {
			logging.CommonLog().Fatalf("Basic auth users list is empty")
		}
	}
}

// type authenticatorFunc func(auth.AuthenticatedHandlerFunc) http.HandlerFunc

func (ctx Auth) GetAuthenticator(mw http.HandlerFunc) http.HandlerFunc {
	// TODO: Read more about "realm" and maybe change it to something other
	authenticator := auth.NewBasicAuthenticator("http-knocker", ctx.basicAuthCheck)
	return authenticator.Wrap(func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		mw.ServeHTTP(w, &r.Request)
	})
}
