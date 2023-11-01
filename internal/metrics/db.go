package metrics

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/superfly/flyctl/agent"
)

var metrics []metricsMessage = make([]metricsMessage, 0)

func queueMetric(metric metricsMessage) {
	metrics = append(metrics, metric)
}

// Spawns a forked `flyctl metrics send` process that sends metrics to the flyctl-metrics server
func FlushMetrics() error {
	json, err := json.Marshal(metrics)
	if err != nil {
		return err

	}

	flyctl, err := os.Executable()
	if err != nil {
		return err
	}

	cmd := exec.Command(flyctl, "metrics", "send")

	buffer := bytes.Buffer{}
	buffer.Write(json)

	cmd.Stdin = &buffer
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "FLY_NO_UPDATE_CHECK=1")

	agent.SetSysProcAttributes(cmd)

	if err := cmd.Process.Release(); err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
}
