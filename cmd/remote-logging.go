package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/grafana/loki-client-go/loki"
	"github.com/metal-stack/pixie/api"
	promconfig "github.com/prometheus/common/config"
	slogloki "github.com/samber/slog-loki/v3"
	slogmulti "github.com/samber/slog-multi"
)

func AddRemoteHandler(spec *Specification, handler slog.Handler) (slog.Handler, error) {
	metalConfig := spec.MetalConfig
	if metalConfig.Logging == nil || metalConfig.Logging.Endpoint == "" {
		return handler, nil
	}
	if metalConfig.Logging.Type != api.LogTypeLoki {
		slog.New(handler).Error("unsupported remote logging type, ignoring", "type", metalConfig.Logging.Type)
		return handler, nil
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

	// it is important that the loki handler does not block the metal-hammer under any circumstances
	// therefore we just throw away messages on error and solely rely on the default stdout logger
	// in case of backend unavailability
	failoverHandler := slogmulti.Failover()(
		slogmulti.Pipe(mdw).Handler(lokiHandler),
		newDropHandler(os.Stdout),
	)

	combinedHandler := slogmulti.Fanout(failoverHandler, handler)

	return combinedHandler, nil
}

func jsonFormattingMiddleware(ctx context.Context, record slog.Record, next func(context.Context, slog.Record) error) error {
	attrs := map[string]string{"msg": record.Message, "level": record.Level.String(), "time": record.Time.String()}

	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value.String()
		return true
	})

	r, err := json.Marshal(attrs)
	if err != nil {
		return fmt.Errorf("unable to marshal log attributes %w", err)
	}
	record = slog.NewRecord(record.Time, record.Level, string(r), record.PC)
	return next(ctx, record)
}

type dropHandler struct {
	w io.Writer
}

func newDropHandler(writer io.Writer) slog.Handler {
	return &dropHandler{
		w: writer,
	}
}

func (h *dropHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

func (h *dropHandler) Handle(ctx context.Context, record slog.Record) error {
	fmt.Fprintf(h.w, "dropped record %s", record.Message)
	return nil
}

func (h *dropHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *dropHandler) WithGroup(name string) slog.Handler {
	return h
}
