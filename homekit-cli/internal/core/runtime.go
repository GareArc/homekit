package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/homekit/homekit-cli/internal/util/bufutil"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type ctxKey struct{}

// Options controls bootstrap behaviour for the CLI runtime.
type Options struct {
	ConfigPath string
	LogLevel   string
	LogFormat  string
	NoColor    bool
	DryRun     bool
}

// Runtime represents initialized application state shared across commands.
type Runtime struct {
	Context context.Context
	Config  Config
	Logger  zerolog.Logger
	Version VersionInfo
	DryRun  bool
	BufPool *bufutil.Pool
}

// VersionInfo carries build metadata injected at link-time.
type VersionInfo struct {
	Version string
	Commit  string
	Date    string
	Source  string
}

// Config holds the application configuration.
type Config struct {
	AssetOverrides string   `mapstructure:"asset_overrides"`
	PluginPaths    []string `mapstructure:"plugin_paths"`
	TempDir        string   `mapstructure:"temp_dir"`
	LogLevel       string   `mapstructure:"log_level"`
	// Add other fields as needed
}

// DefaultConfigPath returns the default user config file location.
func DefaultConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "homekit", "config.yaml"), nil
}

// WithRuntime attaches the runtime instance to a context for downstream use.
func WithRuntime(ctx context.Context, rt *Runtime) context.Context {
	return context.WithValue(ctx, ctxKey{}, rt)
}

// FromContext extracts the runtime from context when previously registered.
func FromContext(ctx context.Context) (*Runtime, bool) {
	rt, ok := ctx.Value(ctxKey{}).(*Runtime)
	return rt, ok
}

// Bootstrap initializes configuration and logging, returning the runtime model.
func Bootstrap(parent context.Context, opts Options, version VersionInfo) (*Runtime, error) {
	ctx := parent
	if ctx == nil {
		ctx = context.Background()
	}

	var cfg Config
	var err error
	if opts.ConfigPath != "" {
		cfg, err = loadConfig(opts.ConfigPath)
		if err != nil {
			return nil, err
		}
	}

	logger, err := configureLogger(opts)
	if err != nil {
		return nil, err
	}

	bufPool := bufutil.NewPool(1024, 1024*1024)

	rt := &Runtime{
		Context: ctx,
		Config:  cfg,
		Logger:  logger,
		Version: normalizeVersion(version),
		DryRun:  opts.DryRun,
		BufPool: bufPool,
	}

	rt.Context = WithRuntime(ctx, rt)
	return rt, nil
}

func normalizeVersion(v VersionInfo) VersionInfo {
	if v.Version == "" {
		v.Version = "dev"
	}
	if v.Date == "" {
		v.Date = time.Now().UTC().Format(time.RFC3339)
	}
	if v.Source == "" {
		_, file, _, _ := runtime.Caller(1)
		v.Source = filepath.Dir(file)
	}
	return v
}

func configureLogger(opts Options) (zerolog.Logger, error) {
	level := zerolog.InfoLevel
	if opts.LogLevel != "" {
		l, err := zerolog.ParseLevel(strings.ToLower(opts.LogLevel))
		if err != nil {
			return zerolog.Logger{}, fmt.Errorf("parse log level %q: %w", opts.LogLevel, err)
		}
		level = l
	}

	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		NoColor:    opts.NoColor || strings.ToLower(opts.LogFormat) == "json",
	}

	var logger zerolog.Logger
	if strings.ToLower(opts.LogFormat) == "json" {
		logger = zerolog.New(os.Stderr).Level(level).With().Timestamp().Logger()
	} else {
		logger = zerolog.New(output).Level(level).With().Timestamp().Logger()
	}

	log.Logger = logger
	return logger, nil
}

func loadConfig(path string) (Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	v.SetEnvPrefix("HOMEKIT")
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		var cfgErr viper.ConfigFileNotFoundError
		if !errors.As(err, &cfgErr) {
			return Config{}, fmt.Errorf("read config at %s: %w", path, err)
		}
	}

	// Unmarshal into structured config
	var cfg Config

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "mapstructure",
		Result:  &cfg,
	})
	if err != nil {
		return Config{}, fmt.Errorf("create decoder: %w", err)
	}

	if err := decoder.Decode(v.AllSettings()); err != nil {
		return Config{}, fmt.Errorf("decode config: %w", err)
	}

	return cfg, nil
}
