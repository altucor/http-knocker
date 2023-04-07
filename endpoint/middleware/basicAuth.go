package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	auth "github.com/abbot/go-http-auth"
	"github.com/altucor/http-knocker/logging"
	"gopkg.in/yaml.v3"
)

type BasicAuthCfg struct {
	UsersFile string   `yaml:"users-file"`
	Users     []string `yaml:"users"`
}

type BasicAuth struct {
	authUsers map[string]string
}

func parseHtpasswdUserLine(line string) (string, string, error) {
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return "", "", errors.New("cannot parse htpasswd line")
	}
	return parts[0], parts[1], nil
}

func (ctx *BasicAuth) setHtpasswdUsersFromArray(users []string) error {
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

func (ctx *BasicAuth) setHtpasswdUsersFromFile(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		logging.CommonLog().Errorf("cannot parse htpasswd file: %s\n", file)
		return errors.New("cannot parse htpasswd file")
	}

	lines := strings.Split(string(data), "\n")
	return ctx.setHtpasswdUsersFromArray(lines)
}

func BasicAuthNew(cfg BasicAuthCfg) *BasicAuth {
	auth := &BasicAuth{}

	if len(cfg.Users) != 0 {
		err := auth.setHtpasswdUsersFromArray(cfg.Users)
		if err != nil {
			logging.CommonLog().Fatal("Cannot process htpassd users array")
		}
	}
	if cfg.UsersFile != "" {
		err := auth.setHtpasswdUsersFromFile(
			cfg.UsersFile)
		if err != nil {
			logging.CommonLog().Fatal("Cannot process htpassd users file")
		}
	}
	if len(auth.authUsers) == 0 {
		logging.CommonLog().Warn("Basic-Auth users list is empty. No one able to authenticate")
	}
	return auth
}

func BasicAuthNewFromYaml(value *yaml.Node) (*BasicAuth, error) {
	var cfg struct {
		Conn BasicAuthCfg `yaml:"config"`
	}
	if err := value.Decode(&cfg); err != nil {
		return nil, err
	}
	return BasicAuthNew(cfg.Conn), nil
}

func (ctx *BasicAuth) basicAuthCheck(user string, realm string) string {
	passHash, ok := ctx.authUsers[user]
	if ok {
		return passHash
	}
	return ""
}

func (ctx *BasicAuth) Register(handler http.HandlerFunc) http.HandlerFunc {
	// TODO: Read more about "realm" and maybe change it to something other
	authenticator := auth.NewBasicAuthenticator("http-knocker", ctx.basicAuthCheck)
	return authenticator.Wrap(func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		handler.ServeHTTP(w, &r.Request)
	})
}

func (ctx *BasicAuth) Intercept(w http.ResponseWriter, r *http.Request) error {
	return nil
}
