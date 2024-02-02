package cmd

import (
	"fmt"
	"log/slog"

	"github.com/grafana/loki-client-go/loki"
	"github.com/grafana/loki-client-go/pkg/labelutil"
	"github.com/metal-stack/pixie/api"
	promconfig "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	slogloki "github.com/samber/slog-loki/v3"
	slogmulti "github.com/samber/slog-multi"
)

func AddRemoteLoggerFrom(pixieURL string, handler slog.Handler, machineID string) (*slog.Logger, error) {
	metalConfig, err := fetchConfig(pixieURL)
	if err != nil {
		return nil, err
	}
	if metalConfig.Logging == nil {
		return slog.New(handler), nil
	}
	if metalConfig.Logging.Type != api.LogTypeLoki {
		slog.New(handler).Error("unsupported remote logging type, ignoring", "type", metalConfig.Logging.Type)
		return slog.New(handler), nil
	}
	httpClient := promconfig.DefaultHTTPClientConfig
	if metalConfig.Logging.BasicAuth != nil {
		httpClient.BasicAuth = &promconfig.BasicAuth{
			Username: metalConfig.Logging.BasicAuth.User,
			Password: promconfig.Secret(metalConfig.Logging.BasicAuth.Password),
		}
	}
	if metalConfig.Logging.CertificateAuth != nil {
		httpClient.TLSConfig = promconfig.TLSConfig{
			Cert:               metalConfig.Logging.CertificateAuth.Cert,
			Key:                promconfig.Secret(metalConfig.Logging.CertificateAuth.Key),
			InsecureSkipVerify: metalConfig.Logging.CertificateAuth.InsecureSkipVerify,
		}
	}

	config, _ := loki.NewDefaultConfig(metalConfig.Logging.Endpoint)
	config.ExternalLabels = labelutil.LabelSet{
		LabelSet: model.LabelSet{
			"component": "metal-hammer",
			"machineID": model.LabelValue(machineID),
		},
	}
	config.Client = httpClient
	client, err := loki.New(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create loki client %w", err)
	}

	lokiHandler := slogloki.Option{Level: slog.LevelDebug, Client: client}.NewLokiHandler()

	logger := slog.New(slogmulti.Fanout(lokiHandler, handler))
	return logger, nil
}
