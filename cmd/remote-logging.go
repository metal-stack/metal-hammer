package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/grafana/loki-client-go/loki"
	"github.com/metal-stack/pixie/api"
	promconfig "github.com/prometheus/common/config"
	slogloki "github.com/samber/slog-loki/v3"
	slogmulti "github.com/samber/slog-multi"
)

func AddRemoteLoggerFrom(spec *Specification, handler slog.Handler) (*slog.Logger, error) {
	metalConfig := spec.MetalConfig
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

	config, err := loki.NewDefaultConfig(metalConfig.Logging.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to create loki default config %w", err)
	}
	// config.EncodeJson = true
	config.Client = httpClient
	client, err := loki.New(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create loki client %w", err)
	}

	lokiHandler := slogloki.Option{
		Level:  slog.LevelDebug,
		Client: client}.NewLokiHandler().WithAttrs(
		[]slog.Attr{
			{Key: "component", Value: slog.StringValue("metal-hammer")},
			{Key: "machineID", Value: slog.StringValue(spec.MachineUUID)},
		},
	)
	mdw := slogmulti.NewHandleInlineMiddleware(jsonFormattingMiddleware)
	logger := slog.New(slogmulti.Fanout(slogmulti.Pipe(mdw).Handler(lokiHandler), handler))
	logger.Info("remote logging to loki", "url", metalConfig.Logging.Endpoint, "machineID", spec.MachineUUID)
	return logger, nil
}

func jsonFormattingMiddleware(ctx context.Context, record slog.Record, next func(context.Context, slog.Record) error) error {
	attrs := map[string]string{"msg": record.Message, "level": record.Level.String(), "time": record.Time.Local().String()}

	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value.String()
		return true
	})

	r, _ := json.Marshal(attrs)
	record = slog.NewRecord(record.Time, record.Level, string(r), record.PC)
	return next(ctx, record)
}
