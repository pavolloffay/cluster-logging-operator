package functional

import (
	log "github.com/ViaQ/logerr/v2/log/static"
	"github.com/openshift/cluster-logging-operator/internal/runtime"
	"strings"
)

const (
	OTELReceiverConf = `
exporters:
  debug:
    verbosity: detailed
receivers:
  otlp:
    protocols:
      http:
        endpoint: 0.0.0.0:8090
service:
  telemetry:
    logs: 
      level: debug
  pipelines:
    logs:
      receivers: [otlp]
      exporters: [debug]
`
	OTELImage = "ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector:0.93.0"
)

func (f *CollectorFunctionalFramework) AddOTELCollector(b *runtime.PodBuilder, outputName string) error {
	log.V(3).Info("Adding OTEL collector", "name", outputName)
	name := strings.ToLower(outputName)

	config := runtime.NewConfigMap(b.Pod.Namespace, name, map[string]string{
		"config.yaml": OTELReceiverConf,
	})
	log.V(2).Info("Creating configmap", "namespace", config.Namespace, "name", config.Name, "config.yaml", OTELReceiverConf)
	if err := f.Test.Client.Create(config); err != nil {
		return err
	}

	log.V(2).Info("Adding container", "name", name, "image", OTELImage)
	b.AddContainer(name, OTELImage).
		AddVolumeMount(config.Name, "/etc/functional", "", false).
		WithCmdArgs([]string{"--config", "/etc/functional/config.yaml"}).
		End().
		AddConfigMapVolume(config.Name, config.Name)
	return nil
}

func (f *CollectorFunctionalFramework) AddEcho(b *runtime.PodBuilder, outputName string) error {
	log.V(3).Info("Adding echo server", "name", outputName)

	name := strings.ToLower(outputName)
	b.AddContainer(name, "mendhak/http-https-echo:31").
		AddEnvVar("HTTP_PORT", "8090").
		End()
	return nil
}
