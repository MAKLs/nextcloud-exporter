package exporter

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"

	"github.com/MAKLs/nextcloud-exporter/client"
	"github.com/MAKLs/nextcloud-exporter/metrics"
	"github.com/MAKLs/nextcloud-exporter/models"
	"github.com/prometheus/client_golang/prometheus"
)

// NCExporter collects server info from Nextcloud.
type NCExporter struct {
	client         client.Client
	lock           sync.Mutex
	excludePHP     bool
	excludeStrings bool
	filterMetrics  []string
}

// NewNCExporter creates a new Exporter instance.
func NewNCExporter(client client.Client, excludePHP bool, excludeStrings bool, filterMetrics []string) *NCExporter {
	return &NCExporter{
		client:         client,
		excludePHP:     excludePHP,
		excludeStrings: excludeStrings,
		filterMetrics:  filterMetrics,
	}
}

// Wrapper around `NCClient` to time duration of requests to Nextcloud
func (col *NCExporter) fetchNCServerInfo() (*models.NCServerInfo, error) {
	timer := prometheus.NewTimer(metrics.ScrapeDuration)
	defer timer.ObserveDuration()
	return col.client.FetchNCServerInfo()
}

// Collect fetches the server info from Nextcloud and exposes it as metrics.
func (col *NCExporter) Collect(ch chan<- prometheus.Metric) {
	col.lock.Lock()
	defer col.lock.Unlock()
	log.Println("collecting metrics")

	serverInfo, err := col.fetchNCServerInfo()
	if err != nil {
		metrics.NcUp.Set(0)
		log.Println(err)
	} else {
		metrics.NcUp.Set(1)
		col.mustCollectTaggedMetrics(serverInfo, ch)
	}

	metrics.ScrapeDuration.Collect(ch)
	metrics.NcUp.Collect(ch)
	metrics.ScrapeCount.Collect(ch)
}

// Describe describes the metrics collected by this exporter.
func (col *NCExporter) Describe(ch chan<- *prometheus.Desc) {
	metrics.MetricsCollection.Describe(ch)

	metrics.ScrapeCount.Describe(ch)
	metrics.ScrapeDuration.Describe(ch)
	metrics.NcUp.Describe(ch)
}

func (col *NCExporter) shouldSkipMetric(name string, metricKind reflect.Kind) bool {
	return (strings.HasPrefix(name, "php") && col.excludePHP) || func() bool {
		for _, filter := range col.filterMetrics {
			if strings.Compare(filter, prometheus.BuildFQName(metrics.Namespace, "", name)) == 0 {
				return true
			}
		}
		return false
	}() || (metricKind == reflect.String && col.excludeStrings)
}

func (col *NCExporter) mustCollectTaggedMetrics(v interface{}, ch chan<- prometheus.Metric) error {
	if err := col.collectTaggedMetrics(v, ch); err != nil {
		panic(err)
	} else {
		return nil
	}
}

func (col *NCExporter) collectTaggedMetrics(v interface{}, ch chan<- prometheus.Metric) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// MUST be called with a struct
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("can only reflect fields of structs, received %s", val.Kind())
	}

	for fi := 0; fi < val.NumField(); fi++ {
		field := val.Field(fi)
		fieldKind := field.Type().Kind()

		if fieldKind == reflect.Struct {
			// Recurse through nested structs
			col.collectTaggedMetrics(field.Interface(), ch)
		} else if metricName, ok := val.Type().Field(fi).Tag.Lookup(metrics.MetricTag); ok && !col.shouldSkipMetric(metricName, fieldKind) {
			// If field is tagged with a metric, collect it.
			// A metric label is optional.
			labelValues := make([]string, 0)
			if label, ok := val.Type().Field(fi).Tag.Lookup(metrics.LabelValueTag); ok {
				labelValues = append(labelValues, label)
			}

			if metricTemplate, ok := metrics.MetricsCollection.WithName(metricName); ok {
				switch fieldKind {
				case reflect.Float64:
					ch <- metricTemplate.MustEmitMetric(field.Float(), labelValues...)
				case reflect.Bool:
					var val float64
					if field.Bool() {
						val = 1
					} else {
						val = 0
					}
					ch <- metricTemplate.MustEmitMetric(val, labelValues...)
				case reflect.String:
					labelValues := append(labelValues, field.String())
					ch <- metricTemplate.MustEmitMetric(1, labelValues...)
				default:
					// TODO
				}
			} else {
				log.Printf("%s tagged for export but no corresponding metric template found", metricName)
			}
		}
	}

	return nil
}
