package database

import (
	"os"
	"sync"
	"testing"

	"github.com/code19m/errx"
	"github.com/redis/go-redis/v9"
	"github.com/rise-and-shine/pkg/cfgloader"
	"github.com/rise-and-shine/pkg/rediswr"
)

//nolint:gochecknoglobals // lazy singleton for test Redis connection
var (
	redisOnce sync.Once
	redisInst redis.Cmdable
)

type redisConfig struct {
	Redis rediswr.Config `yaml:"redis" validate:"required"`
}

func GetTestRedis(t *testing.T) redis.Cmdable {
	t.Helper()

	redisOnce.Do(func() {
		client, err := initRedis()
		if err != nil {
			t.Fatalf("failed to initialize test Redis: %v", err)
		}
		redisInst = client
	})

	return redisInst
}

func initRedis() (redis.Cmdable, error) {
	originalWd, err := os.Getwd()
	if err != nil {
		return nil, errx.Newf("get working directory: %v", err)
	}

	root, err := projectRoot()
	if err != nil {
		return nil, errx.Newf("find project root: %v", err)
	}

	if err = os.Chdir(root); err != nil {
		return nil, errx.Newf("change to project root: %v", err)
	}
	defer func() { err = os.Chdir(originalWd) }()

	cfg := cfgloader.MustLoad[redisConfig](cfgloader.WithSilent())

	return rediswr.New(cfg.Redis), nil
}
