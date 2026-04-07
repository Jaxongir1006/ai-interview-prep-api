//go:build system

package trigger

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/code19m/errx"
	"github.com/gavv/httpexpect/v2"
	"github.com/rise-and-shine/pkg/cfgloader"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/spf13/cast"
)

//nolint:gochecknoglobals // lazy singleton for base URL
var (
	baseURLOnce sync.Once
	baseURL     string
)

type config struct {
	HTTPServer server.Config `yaml:"http_server" validate:"required"`
}

// UserAction returns an httpexpect client configured for testing user_action use cases.
func UserAction(t *testing.T) *httpexpect.Expect {
	t.Helper()

	baseURLOnce.Do(func() {
		url, err := loadBaseURL()
		if err != nil {
			t.Fatalf("UserAction: failed to load base URL: %v", err)
		}
		baseURL = url
	})

	return httpexpect.Default(t, baseURL)
}

func loadBaseURL() (url string, err error) { //nolint:nonamedreturns // required to capture deferred error
	originalWd, err := os.Getwd()
	if err != nil {
		return "", errx.Newf("get working directory: %v", err)
	}

	root, err := projectRoot()
	if err != nil {
		return "", errx.Newf("find project root: %v", err)
	}

	if err = os.Chdir(root); err != nil {
		return "", errx.Newf("change to project root: %v", err)
	}
	defer func() { err = os.Chdir(originalWd) }()

	cfg := cfgloader.MustLoad[config](cfgloader.WithSilent())

	return fmt.Sprintf("http://%s", net.JoinHostPort(cfg.HTTPServer.Host, cast.ToString(cfg.HTTPServer.Port))), nil
}

func projectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", errx.Newf("get working directory: %v", err)
	}

	for {
		if _, err = os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errx.Newf("could not find project root (go.mod)")
		}
		dir = parent
	}
}
