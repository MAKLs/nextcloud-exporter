package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Namespace for exported metrics in Prometheus.
	Namespace         = "nextcloud"
	exporterSubsystem = "exporter"
	// MetricTag tags struct fields to be exported as metrics.
	MetricTag = "metric"
	// LabelValueTag tags struct fields to be used as label values in exported metrics.
	LabelValueTag = "label"
)

var (
	// ExporterRegistry is the Prometheus registry from which metrics are gathered.
	ExporterRegistry *prometheus.Registry = prometheus.NewRegistry()
	// MetricsStore stores templates to create metrics for Nextcloud server info.
	MetricsStore *MetricTemplateStore = &MetricTemplateStore{templates: make(map[string]MetricTemplate)}
)

func init() {
	// Populate metrics store
	for name, metricInfo := range ncMetrics {
		template := newMetricTemplate(name, metricInfo.help, metricInfo.valueType, metricInfo.variableLabels, metricInfo.constLabels)
		MetricsStore.mustAddTemplate(name, template)
	}
}

// Exporter metrics
var (
	ScrapeDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: Namespace,
		Subsystem: exporterSubsystem,
		Name:      "scrape_duration_seconds",
		Help:      "Duration of scrapes for Nextcloud metrics.",
	})
	ScrapeCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: exporterSubsystem,
		Name:      "scrape_count",
		Help:      "Count of scrapes partitioned by response code.",
	}, []string{"status_code"})
	NcUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: exporterSubsystem,
		Name:      "up",
		Help:      "Flag indicating whether last scrape was successful.",
	})
)

// Nextcloud metrics
var ncMetrics = map[string]struct {
	help           string
	valueType      prometheus.ValueType
	variableLabels []string
	constLabels    prometheus.Labels
}{
	"nc_version": {
		help:           "Version of Nextcloud installed on this instance.",
		valueType:      prometheus.UntypedValue,
		variableLabels: []string{"version"},
		constLabels:    nil,
	},
	"avatars_enabled": {
		help:           "Flag indicating whether avatars are enabled.",
		valueType:      prometheus.UntypedValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"previews_enabled": {
		help:           "Flag indicating whether previews are enabled.",
		valueType:      prometheus.UntypedValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"memcache_type": {
		help:           "Type of cache configured for this instance.",
		valueType:      prometheus.UntypedValue,
		variableLabels: []string{"cache_location", "cache_type"},
		constLabels:    nil,
	},
	"file_locking_enabled": {
		help:           "Flag indicating whether file locking is enabled.",
		valueType:      prometheus.UntypedValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"memcache_locking_type": {
		help:           "Type of memcache used for file locking.",
		valueType:      prometheus.UntypedValue,
		variableLabels: []string{"cache_type"},
		constLabels:    nil,
	},
	"debug_mode_enabled": {
		help:           "Flag indicating whether debug mode is enabled.",
		valueType:      prometheus.UntypedValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"free_space_bytes": {
		help:           "Free storage space in bytes on this instance.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},

	"installed_apps": {
		help:           "Number of apps installed on this instance.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"app_updates_available": {
		help:           "Number of app updates available on this instance.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"users": {
		help:           "Number of users on this instance.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"files": {
		help:           "Number of files on this instance.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"storages": {
		help:           "Number of storages available on this instance, partitioned by location.",
		valueType:      prometheus.GaugeValue,
		variableLabels: []string{"storage_location"},
		constLabels:    nil,
	},
	"shares": {
		help:           "Number of shares on this instance, partitioned by share type.",
		valueType:      prometheus.GaugeValue,
		variableLabels: []string{"share_type"},
		constLabels:    nil,
	},
	"fed_shares_sent": {
		help:           "Number of federated shares sent.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"fed_shares_received": {
		help:           "Number of federated shares received.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"web_server_type": {
		help:           "Type of web server hosting this instance.",
		valueType:      prometheus.UntypedValue,
		variableLabels: []string{"server_type"},
		constLabels:    nil,
	},
	"php_version": {
		help:           "Version of PHP installed on this instance.",
		valueType:      prometheus.UntypedValue,
		variableLabels: []string{"version"},
		constLabels:    nil,
	},
	"php_memory_limit_bytes": {
		help:           "Configured PHP memory limit.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_max_execution_time_seconds": {
		help:           "Configured PHP max execution time.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_upload_max_file_size_bytes": {
		help:           "Configured PHP upload max file size.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_enabled": {
		help:           "Flag indicating whether PHP OPcache is enabled.",
		valueType:      prometheus.UntypedValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_full": {
		help:           "Flag indicating whether PHP OPcache is full.",
		valueType:      prometheus.UntypedValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_restart_pending": {
		help:           "Flag indicating whether PHP OPcache is pending a restart.",
		valueType:      prometheus.UntypedValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_restart_in_progress": {
		help:           "Flag indicating whether PHP OPcache is restarting.",
		valueType:      prometheus.UntypedValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_memory_used_bytes": {
		help:           "Memory usage of PHP OPcache.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_memory_free_bytes": {
		help:           "Memory available to PHP OPcache.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_memory_wasted_bytes": {
		help:           "Memory wasted by PHP OPcache.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_interned_strings_buffer_size_bytes": {
		help:           "Size of PHP OPcache interned strings buffer.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_interned_strings_memory_used_bytes": {
		help:           "Memory used by PHP OPcache interned strings buffer.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_interned_strings_memory_free_bytes": {
		help:           "Memory available to PHP OPcache interned strings buffer.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_interned_strings_count": {
		help:           "Count of interned strings in PHP OPcache interned strings buffer.",
		valueType:      prometheus.CounterValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_cached_scripts_count": {
		help:           "Count of cached scripts in PHP OPcache.",
		valueType:      prometheus.CounterValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_cached_keys_count": {
		help:           "Count of cached scripts in PHP OPcache.",
		valueType:      prometheus.CounterValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_hits_count": {
		help:           "Count of PHP OPcache hits.",
		valueType:      prometheus.CounterValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_misses_count": {
		help:           "Count of PHP OPcache misses.",
		valueType:      prometheus.CounterValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_blacklist_misses_count": {
		help:           "Count of PHP OPcache blacklist misses.",
		valueType:      prometheus.CounterValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_start_time_ticks": {
		help:           "Start time of PHP OPcache.",
		valueType:      prometheus.CounterValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_last_restart_time_ticks": {
		help:           "Last restart time of PHP OPcache.",
		valueType:      prometheus.CounterValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_opcache_restart_count": {
		help:           "Count of PHP OPcache restarts, partitioned by restart type.",
		valueType:      prometheus.CounterValue,
		variableLabels: []string{"restart_type"},
		constLabels:    nil,
	},
	"php_jit_enabled": {
		help:           "Flag indicating whether PHP JIT is enabled.",
		valueType:      prometheus.UntypedValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_jit_on": {
		help:           "Flag indicating whether PHP JIT is on.",
		valueType:      prometheus.UntypedValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_jit_kind": {
		help:           "Kind of PHP JIT.",
		valueType:      prometheus.UntypedValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_jit_optimization_level": {
		help:           "Optimization level of PHP JIT.",
		valueType:      prometheus.UntypedValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_jit_optimization_flags": {
		help:           "Optimization flags of PHP JIT.",
		valueType:      prometheus.UntypedValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_jit_buffer_size_bytes": {
		help:           "Size of PHP JIT buffer.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_jit_buffer_free_bytes": {
		help:           "Free space in PHP JIT buffer.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_apcu_cache_slots": {
		help:           "Number of slots in PHP APCU cache.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_apcu_cache_ttl": {
		help:           "TTL for entries in PHP APCU cache.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_apcu_cache_hits_count": {
		help:           "Count of PHP APCU cache hits.",
		valueType:      prometheus.CounterValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_apcu_cache_misses_count": {
		help:           "Count of PHP APCU cache misses.",
		valueType:      prometheus.CounterValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_apcu_cache_inserts_count": {
		help:           "Count of PHP APCU cache inserts.",
		valueType:      prometheus.CounterValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_apcu_cache_expunges_count": {
		help:           "Count of PHP APCU cache expunges.",
		valueType:      prometheus.CounterValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_apcu_cache_entries": {
		help:           "Number of entries in PHP APCU cache.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_apcu_cache_start_time_ticks": {
		help:           "Start time of PHP APCU cache.",
		valueType:      prometheus.CounterValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_apcu_cache_memory_free_bytes": {
		help:           "Free memory available to PHP APCU cache.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_apcu_cache_memory_type": {
		help:           "PHP APCU cache memory type.",
		valueType:      prometheus.UntypedValue,
		variableLabels: []string{"memory_type"},
		constLabels:    nil,
	},
	"php_apcu_sma_seg": {
		help:           "Number of PHP APCU shared memory allocation segments.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_apcu_sma_seg_size_bytes": {
		help:           "Size of PHP APCU shared memory allocation segments.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"php_apcu_sma_memory_free_bytes": {
		help:           "Free memory available to PHP APCU shared memory allocation.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"database_type": {
		help:           "Type of database backing this instance.",
		valueType:      prometheus.UntypedValue,
		variableLabels: []string{"database_type"},
		constLabels:    nil,
	},
	"database_version": {
		help:           "Version of database backing this instance.",
		valueType:      prometheus.UntypedValue,
		variableLabels: []string{"version"},
		constLabels:    nil,
	},
	"database_size_bytes": {
		help:           "Size of database backing this instance.",
		valueType:      prometheus.GaugeValue,
		variableLabels: nil,
		constLabels:    nil,
	},
	"active_users": {
		help:           "Number of active users on this instance, partitioned by last t time.",
		valueType:      prometheus.GaugeValue,
		variableLabels: []string{"t"},
		constLabels:    nil,
	},
}

// MetricTemplate stores metadata for a metric exported by the exporter.
// It is useful for generating fixed-value metrics with prometheus.NewConstMetric.
type MetricTemplate struct {
	Desc      *prometheus.Desc     // Description given to metrics generated with this template
	ValueType prometheus.ValueType // Value type of metrics generate with this template
}

// MetricTemplateStore is the collection of all MetricTemplates exported by the exporter.
type MetricTemplateStore struct {
	templates map[string]MetricTemplate
}

// MetricTemplateStoreItem represents an item in the global metrics template collection.
//
// Each item contains the iternal key of the metric in the collection and the template needed to generate the metric.
type MetricTemplateStoreItem struct {
	Key      string
	Template *MetricTemplate
}

// Iter returns a channel that receives a MetricTemplateCollectionItem for each metric template in the collection.
func (store *MetricTemplateStore) Iter() <-chan MetricTemplateStoreItem {
	// TODO: any other way to emulate iterators?
	ch := make(chan MetricTemplateStoreItem)
	go func() {
		defer close(ch)
		for metricKey, metricTemplate := range store.templates {
			ch <- MetricTemplateStoreItem{Key: metricKey, Template: &metricTemplate}
		}
	}()
	return ch
}

func (store *MetricTemplateStore) mustAddTemplate(key string, template MetricTemplate) {
	if _, ok := store.templates[key]; !ok {
		// Template not present... add it
		store.templates[key] = template
	} else {
		// Template already exists with that name
		panic(fmt.Sprintf("template already exists with key '%s': %v", key, template))
	}
}

// WithName returns the metric template from the collection with the given name.
// If no such template exists, returns false.
func (store *MetricTemplateStore) WithName(name string) (MetricTemplate, bool) {
	template, ok := store.templates[name]
	return template, ok
}

func newMetricTemplate(name string, help string, valueType prometheus.ValueType, variableLabels []string, constLabels prometheus.Labels) MetricTemplate {
	fqName := prometheus.BuildFQName(Namespace, "", name)
	return MetricTemplate{
		Desc:      prometheus.NewDesc(fqName, help, variableLabels, constLabels),
		ValueType: valueType,
	}
}

// MustEmitMetric generates a fixed-value metric from the provided value and labelValues.
// Panics if the metric cannot be generated.
func (template *MetricTemplate) MustEmitMetric(value float64, labelValues ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(template.Desc, template.ValueType, value, labelValues...)
}
