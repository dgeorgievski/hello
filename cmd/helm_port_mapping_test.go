package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// helmDocument represents a single Kubernetes resource from "helm template" output.
type helmDocument struct {
	APIVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Metadata   map[string]interface{} `yaml:"metadata"`
	Spec       map[string]interface{} `yaml:"spec"`
}

// containerPortInfo holds the port configuration of a container.
type containerPortInfo struct {
	Name          string
	ContainerPort int
	Protocol      string
}

// probeInfo holds the configuration of an HTTP probe.
type probeInfo struct {
	Path string
	Port interface{} // string (named port) or int
}

// parseHelmDocuments splits multi-document YAML (separated by "---") into individual documents.
func parseHelmDocuments(yamlData []byte) ([]*helmDocument, error) {
	var docs []*helmDocument
	sections := strings.Split(string(yamlData), "\n---")
	for _, section := range sections {
		section = strings.TrimSpace(section)
		if section == "" {
			continue
		}
		var doc helmDocument
		if err := yaml.Unmarshal([]byte(section), &doc); err != nil {
			return nil, err
		}
		if doc.Kind != "" {
			docs = append(docs, &doc)
		}
	}
	return docs, nil
}

// extractContainerPortInfo parses the containers list from a Deployment spec
// and returns port information for the first container with matching port name.
func extractContainerPortInfo(spec map[string]interface{}, portName string) (*containerPortInfo, error) {
	template, ok := spec["template"].(map[string]interface{})
	if !ok {
		return nil, nil
	}
	podSpec, ok := template["spec"].(map[string]interface{})
	if !ok {
		return nil, nil
	}
	containers, ok := podSpec["containers"].([]interface{})
	if !ok {
		return nil, nil
	}
	for _, raw := range containers {
		container, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		ports, ok := container["ports"].([]interface{})
		if !ok {
			continue
		}
		for _, rawPort := range ports {
			p, ok := rawPort.(map[string]interface{})
			if !ok {
				continue
			}
			if p["name"] == portName {
				portNum := p["containerPort"]
				var port int
				switch v := portNum.(type) {
				case int:
					port = v
				case float64:
					port = int(v)
				default:
					return nil, nil
				}
				return &containerPortInfo{
					Name:          p["name"].(string),
					ContainerPort: port,
					Protocol:      p["protocol"].(string),
				}, nil
			}
		}
	}
	return nil, nil
}

// extractProbeInfo extracts probe configuration from a container spec.
func extractProbeInfo(container map[string]interface{}, probeName string) *probeInfo {
	probe, ok := container[probeName].(map[string]interface{})
	if !ok {
		return nil
	}
	httpGet, ok := probe["httpGet"].(map[string]interface{})
	if !ok {
		return nil
	}
	var port interface{}
	if p, exists := httpGet["port"]; exists {
		switch v := p.(type) {
		case string:
			port = v
		case float64:
			port = int(v)
		}
	}
	return &probeInfo{
		Path: httpGet["path"].(string),
		Port: port,
	}
}

// intToString converts an integer to its string representation.
func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}

func TestHelmPortMapping(t *testing.T) {
	tests := []struct {
		name          string
		servicePort   int
		containerPort int
	}{
		{
			name:          "default service port",
			servicePort:   80,
			containerPort: 8080,
		},
		{
			name:          "custom service port",
			servicePort:   9090,
			containerPort: 8080,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Render Helm chart with the selected Service port.
			chartPath, err := filepath.Abs("../deploy/helm/hello")
			if err != nil {
				t.Fatalf("failed to resolve chart path: %v", err)
			}
			args := []string{"template", "hello", chartPath}

			if tt.servicePort != 80 {
				// For non-default ports, use a temporary values file.
				dir := t.TempDir()
				valuesFile := filepath.Join(dir, "overrides.yaml")
				content := "service:\n  port: " + intToString(tt.servicePort) + "\n"
				if err := os.WriteFile(valuesFile, []byte(content), 0644); err != nil {
					t.Fatalf("failed to write override file: %v", err)
				}
				args = append(args, "--values", valuesFile)
			}

			cmd := exec.Command("helm", args...)
			output, err := cmd.Output()
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					t.Fatalf("helm template failed: %s", string(exitErr.Stderr))
				}
				t.Fatalf("helm template command failed: %v", err)
			}

			// Parse the resulting Service and Deployment objects.
			docs, err := parseHelmDocuments(output)
			if err != nil {
				t.Fatalf("failed to parse helm template output: %v", err)
			}

			var svcPort, svcTargetPort interface{}
			var containerInfo *containerPortInfo
			var liveProbe, readyProbe *probeInfo

			for _, doc := range docs {
				switch doc.Kind {
				case "Service":
					ports, ok := doc.Spec["ports"].([]interface{})
					if !ok || len(ports) == 0 {
						t.Fatal("Service has no ports defined")
					}
					p := ports[0].(map[string]interface{})
					svcPort = p["port"]
					svcTargetPort = p["targetPort"]

				case "Deployment":
					template, ok := doc.Spec["template"].(map[string]interface{})
					if !ok {
						t.Fatal("Deployment missing template spec")
					}
					podSpec, ok := template["spec"].(map[string]interface{})
					if !ok {
						t.Fatal("Deployment template missing spec")
					}
					containers, ok := podSpec["containers"].([]interface{})
					if !ok || len(containers) == 0 {
						t.Fatal("Deployment has no containers defined")
					}
					container := containers[0].(map[string]interface{})
					containerInfo, _ = extractContainerPortInfo(doc.Spec, "http")
					liveProbe = extractProbeInfo(container, "livenessProbe")
					readyProbe = extractProbeInfo(container, "readinessProbe")
				}
			}

			// Verify Service port matches the input value.
			if got := svcPort; got != tt.servicePort {
				t.Errorf("Service port = %v, want %d", got, tt.servicePort)
			}

			// Verify Service targetPort uses the named port "http" (resolves to containerPort at runtime).
			if got := svcTargetPort; got != "http" {
				t.Errorf("Service targetPort = %v, want \"http\"", got)
			}

			// Verify containerPort and HTTP probe resolution.
			if containerInfo == nil {
				t.Fatal("container port info not found")
			}
			if got := containerInfo.ContainerPort; got != tt.containerPort {
				t.Errorf("containerPort = %d, want %d", got, tt.containerPort)
			}
			if got := containerInfo.Name; got != "http" {
				t.Errorf("container port name = %q, want \"http\"", got)
			}

			// Verify liveness probe targets the correct path and named port.
			if liveProbe == nil {
				t.Fatal("livenessProbe not found in rendered Deployment")
			}
			if got := liveProbe.Path; got != "/healthz" {
				t.Errorf("livenessProbe path = %q, want \"/healthz\"", got)
			}
			if got := liveProbe.Port; got != "http" {
				t.Errorf("livenessProbe port = %v (resolves to containerPort %d), want \"http\"", got, tt.containerPort)
			}

			// Verify readiness probe targets the correct path and named port.
			if readyProbe == nil {
				t.Fatal("readinessProbe not found in rendered Deployment")
			}
			if got := readyProbe.Path; got != "/readyz" {
				t.Errorf("readinessProbe path = %q, want \"/readyz\"", got)
			}
			if got := readyProbe.Port; got != "http" {
				t.Errorf("readinessProbe port = %v (resolves to containerPort %d), want \"http\"", got, tt.containerPort)
			}
		})
	}
}
