package config

import (
	"fmt"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	"autoscaler/db"
	"autoscaler/helpers"
	"autoscaler/models"
)

const (
	DefaultLoggingLevel                   string        = "info"
	DefaultServerPort                     int           = 8080
	DefaultPolicyPollerInterval           time.Duration = 40 * time.Second
	DefaultAggregatorExecuteInterval      time.Duration = 40 * time.Second
	DefaultSaveInterval                   time.Duration = 5 * time.Second
	DefaultMetricPollerCount              int           = 20
	DefaultAppMonitorChannelSize          int           = 200
	DefaultAppMetricChannelSize           int           = 200
	DefaultEvaluationExecuteInterval      time.Duration = 40 * time.Second
	DefaultEvaluatorCount                 int           = 20
	DefaultTriggerArrayChannelSize        int           = 200
	DefaultBackOffInitialInterval         time.Duration = 5 * time.Minute
	DefaultBackOffMaxInterval             time.Duration = 2 * time.Hour
	DefaultBreakerConsecutiveFailureCount int64         = 3
)

type ServerConfig struct {
	Port      int             `yaml:"port"`
	TLS       models.TLSCerts `yaml:"tls"`
	NodeAddrs []string        `yaml:"node_addrs"`
	NodeIndex int             `yaml:"node_index"`
}
type DBConfig struct {
	PolicyDB    db.DatabaseConfig `yaml:"policy_db"`
	AppMetricDB db.DatabaseConfig `yaml:"app_metrics_db"`
}

type AggregatorConfig struct {
	MetricPollerCount         int           `yaml:"metric_poller_count"`
	AppMonitorChannelSize     int           `yaml:"app_monitor_channel_size"`
	AppMetricChannelSize      int           `yaml:"app_metric_channel_size"`
	AggregatorExecuteInterval time.Duration `yaml:"aggregator_execute_interval"`
	PolicyPollerInterval      time.Duration `yaml:"policy_poller_interval"`
	SaveInterval              time.Duration `yaml:"save_interval"`
}

type EvaluatorConfig struct {
	EvaluatorCount            int           `yaml:"evaluator_count"`
	TriggerArrayChannelSize   int           `yaml:"trigger_array_channel_size"`
	EvaluationManagerInterval time.Duration `yaml:"evaluation_manager_execute_interval"`
}

type ScalingEngineConfig struct {
	ScalingEngineUrl string          `yaml:"scaling_engine_url"`
	TLSClientCerts   models.TLSCerts `yaml:"tls"`
}

type MetricCollectorConfig struct {
	MetricCollectorUrl string          `yaml:"metric_collector_url"`
	TLSClientCerts     models.TLSCerts `yaml:"tls"`
}

type CircuitBreakerConfig struct {
	BackOffInitialInterval  time.Duration `yaml:"back_off_initial_interval"`
	BackOffMaxInterval      time.Duration `yaml:"back_off_max_interval"`
	ConsecutiveFailureCount int64         `yaml:"consecutive_failure_count"`
}

type Config struct {
	Logging                   helpers.LoggingConfig `yaml:"logging"`
	Server                    ServerConfig          `yaml:"server"`
	DB                        DBConfig              `yaml:"db"`
	Aggregator                AggregatorConfig      `yaml:"aggregator"`
	Evaluator                 EvaluatorConfig       `yaml:"evaluator"`
	ScalingEngine             ScalingEngineConfig   `yaml:"scalingEngine"`
	MetricCollector           MetricCollectorConfig `yaml:"metricCollector"`
	DefaultStatWindowSecs     int                   `yaml:"defaultStatWindowSecs"`
	DefaultBreachDurationSecs int                   `yaml:"defaultBreachDurationSecs"`
	CircuitBreaker            CircuitBreakerConfig  `yaml:"circuitBreaker"`
}

func LoadConfig(bytes []byte) (*Config, error) {
	conf := &Config{
		Logging: helpers.LoggingConfig{
			Level: DefaultLoggingLevel,
		},
		Server: ServerConfig{
			Port: DefaultServerPort,
		},
		Aggregator: AggregatorConfig{
			AggregatorExecuteInterval: DefaultAggregatorExecuteInterval,
			PolicyPollerInterval:      DefaultPolicyPollerInterval,
			SaveInterval:              DefaultSaveInterval,
			MetricPollerCount:         DefaultMetricPollerCount,
			AppMonitorChannelSize:     DefaultAppMonitorChannelSize,
			AppMetricChannelSize:      DefaultAppMetricChannelSize,
		},
		Evaluator: EvaluatorConfig{
			EvaluationManagerInterval: DefaultEvaluationExecuteInterval,
			EvaluatorCount:            DefaultEvaluatorCount,
			TriggerArrayChannelSize:   DefaultTriggerArrayChannelSize,
		},
	}
	err := yaml.Unmarshal(bytes, &conf)
	if err != nil {
		return nil, err
	}

	conf.Logging.Level = strings.ToLower(conf.Logging.Level)
	if conf.CircuitBreaker.ConsecutiveFailureCount == 0 {
		conf.CircuitBreaker.ConsecutiveFailureCount = DefaultBreakerConsecutiveFailureCount
	}
	if conf.CircuitBreaker.BackOffInitialInterval == 0 {
		conf.CircuitBreaker.BackOffInitialInterval = DefaultBackOffInitialInterval
	}
	if conf.CircuitBreaker.BackOffMaxInterval == 0 {
		conf.CircuitBreaker.BackOffMaxInterval = DefaultBackOffMaxInterval
	}
	return conf, nil
}

func (c *Config) Validate() error {
	if c.DB.PolicyDB.Url == "" {
		return fmt.Errorf("Configuration error: Policy DB url is empty")
	}
	if c.DB.AppMetricDB.Url == "" {
		return fmt.Errorf("Configuration error: AppMetric DB url is empty")
	}
	if c.ScalingEngine.ScalingEngineUrl == "" {
		return fmt.Errorf("Configuration error: Scaling engine url is empty")
	}
	if c.MetricCollector.MetricCollectorUrl == "" {
		return fmt.Errorf("Configuration error: Metric collector url is empty")
	}
	if c.Aggregator.AggregatorExecuteInterval <= time.Duration(0) {
		return fmt.Errorf("Configuration error: aggregator execute interval is less-equal than 0")
	}
	if c.Aggregator.PolicyPollerInterval <= time.Duration(0) {
		return fmt.Errorf("Configuration error: policy poller interval is less-equal than 0")
	}
	if c.Aggregator.SaveInterval <= time.Duration(0) {
		return fmt.Errorf("Configuration error: save interval is less-equal than 0")
	}
	if c.Aggregator.MetricPollerCount <= 0 {
		return fmt.Errorf("Configuration error: metric poller count is less-equal than 0")
	}
	if c.Aggregator.AppMonitorChannelSize <= 0 {
		return fmt.Errorf("Configuration error: appMonitor channel size is less-equal than 0")
	}
	if c.Aggregator.AppMetricChannelSize <= 0 {
		return fmt.Errorf("Configuration error: appMetric channel size is less-equal than 0")
	}
	if c.Evaluator.EvaluationManagerInterval <= time.Duration(0) {
		return fmt.Errorf("Configuration error: evalution manager execeute interval is less-equal than 0")
	}
	if c.Evaluator.EvaluatorCount <= 0 {
		return fmt.Errorf("Configuration error: evaluator count is less-equal than 0")
	}
	if c.Evaluator.TriggerArrayChannelSize <= 0 {
		return fmt.Errorf("Configuration error: trigger-array channel size is less-equal than 0")
	}
	if c.DefaultBreachDurationSecs < 60 || c.DefaultBreachDurationSecs > 3600 {
		return fmt.Errorf("Configuration error: defaultBreachDurationSecs should be between 60 and 3600")
	}
	if c.DefaultStatWindowSecs < 60 || c.DefaultStatWindowSecs > 3600 {
		return fmt.Errorf("Configuration error: defaultStatWindowSecs should be between 60 and 3600")
	}

	if (c.Server.NodeIndex >= len(c.Server.NodeAddrs)) || (c.Server.NodeIndex < 0) {
		return fmt.Errorf("Configuration error: node_index out of range")
	}
	return nil

}
