package metrics

import (
	"errors"
	"math/rand"
	"runtime"
	"strconv"
)

const (
	Alloc         = "Alloc"
	BuckHashSys   = "BuckHashSys"
	Frees         = "Frees"
	GCCPUFraction = "GCCPUFraction"
	GCSys         = "GCSys"
	HeapAlloc     = "HeapAlloc"
	HeapIdle      = "HeapIdle"
	HeapInuse     = "HeapInuse"
	HeapObjects   = "HeapObjects"
	HeapReleased  = "HeapReleased"
	HeapSys       = "HeapSys"
	LastGC        = "LastGC"
	Lookups       = "Lookups"
	MCacheInuse   = "MCacheInuse"
	MCacheSys     = "MCacheSys"
	MSpanInuse    = "MSpanInuse"
	MSpanSys      = "MSpanSys"
	Mallocs       = "Mallocs"
	NextGC        = "NextGC"
	NumForcedGC   = "NumForcedGC"
	NumGC         = "NumGC"
	OtherSys      = "OtherSys"
	PauseTotalNs  = "PauseTotalNs"
	StackInuse    = "StackInuse"
	StackSys      = "StackSys"
	Sys           = "Sys"
	TotalAlloc    = "TotalAlloc"
	//
	RandomValue = "RandomValue"
)

const PollCount = "PollCount"

type Metricer interface {
	UpdateMetricsGauge() *map[string]string
	UpdateMetricsCounter() (uint64, error)
}

type Metrics struct {
	MetricsGauge   map[string]string
	MetricsCounter map[string]uint64
}

func NewMetriсs() Metricer {
	var meticsGauge = map[string]string{
		Alloc:         "",
		BuckHashSys:   "",
		Frees:         "",
		GCCPUFraction: "", //float64 :(
		GCSys:         "",
		HeapAlloc:     "",
		HeapIdle:      "",
		HeapInuse:     "",
		HeapObjects:   "",
		HeapReleased:  "",
		HeapSys:       "",
		LastGC:        "",
		Lookups:       "",
		MCacheInuse:   "",
		MCacheSys:     "",
		MSpanInuse:    "",
		MSpanSys:      "",
		Mallocs:       "",
		NextGC:        "",
		NumForcedGC:   "",
		NumGC:         "",
		OtherSys:      "",
		PauseTotalNs:  "",
		StackInuse:    "",
		StackSys:      "",
		Sys:           "",
		TotalAlloc:    "",
		//
		RandomValue: "",
	}
	var metricsCounter = map[string]uint64{
		PollCount: 0,
	}
	return &Metrics{MetricsGauge: meticsGauge, MetricsCounter: metricsCounter}
}

func (m *Metrics) UpdateMetricsGauge() *map[string]string {
	var runtimeMetrics runtime.MemStats
	runtime.ReadMemStats(&runtimeMetrics)

	m.MetricsGauge[Alloc] = strconv.FormatUint(runtimeMetrics.Alloc, 10)
	m.MetricsGauge[BuckHashSys] = strconv.FormatUint(runtimeMetrics.BuckHashSys, 10)
	m.MetricsGauge[Frees] = strconv.FormatUint(runtimeMetrics.Frees, 10)
	m.MetricsGauge[GCCPUFraction] = strconv.FormatFloat(runtimeMetrics.GCCPUFraction, 'f', -1, 64)
	m.MetricsGauge[GCSys] = strconv.FormatUint(runtimeMetrics.GCSys, 10)
	m.MetricsGauge[HeapAlloc] = strconv.FormatUint(runtimeMetrics.HeapAlloc, 10)
	m.MetricsGauge[HeapIdle] = strconv.FormatUint(runtimeMetrics.HeapIdle, 10)
	m.MetricsGauge[HeapInuse] = strconv.FormatUint(runtimeMetrics.HeapInuse, 10)
	m.MetricsGauge[HeapObjects] = strconv.FormatUint(runtimeMetrics.HeapObjects, 10)
	m.MetricsGauge[HeapReleased] = strconv.FormatUint(runtimeMetrics.HeapReleased, 10)
	m.MetricsGauge[HeapSys] = strconv.FormatUint(runtimeMetrics.HeapSys, 10)
	m.MetricsGauge[LastGC] = strconv.FormatUint(runtimeMetrics.LastGC, 10)
	m.MetricsGauge[Lookups] = strconv.FormatUint(runtimeMetrics.Lookups, 10)
	m.MetricsGauge[MCacheInuse] = strconv.FormatUint(runtimeMetrics.MCacheInuse, 10)
	m.MetricsGauge[MCacheSys] = strconv.FormatUint(runtimeMetrics.MCacheSys, 10)
	m.MetricsGauge[MSpanInuse] = strconv.FormatUint(runtimeMetrics.MSpanInuse, 10)
	m.MetricsGauge[MSpanSys] = strconv.FormatUint(runtimeMetrics.MSpanSys, 10)
	m.MetricsGauge[Mallocs] = strconv.FormatUint(runtimeMetrics.Mallocs, 10)
	m.MetricsGauge[NextGC] = strconv.FormatUint(runtimeMetrics.NextGC, 10)
	m.MetricsGauge[NumForcedGC] = strconv.Itoa(int(runtimeMetrics.NumForcedGC)) //?
	m.MetricsGauge[NumGC] = strconv.Itoa(int(runtimeMetrics.NumGC))             //?
	m.MetricsGauge[OtherSys] = strconv.FormatUint(runtimeMetrics.OtherSys, 10)
	m.MetricsGauge[PauseTotalNs] = strconv.FormatUint(runtimeMetrics.PauseTotalNs, 10)

	m.MetricsGauge[StackInuse] = strconv.FormatUint(runtimeMetrics.StackInuse, 10)
	m.MetricsGauge[StackSys] = strconv.FormatUint(runtimeMetrics.StackSys, 10)
	m.MetricsGauge[Sys] = strconv.FormatUint(runtimeMetrics.Sys, 10)
	m.MetricsGauge[TotalAlloc] = strconv.FormatUint(runtimeMetrics.TotalAlloc, 10)
	m.MetricsGauge[RandomValue] = strconv.FormatFloat(rand.Float64(), 'f', -1, 64)
	return &m.MetricsGauge
}
func (m *Metrics) UpdateMetricsCounter() (uint64, error) {
	value, ok := m.MetricsCounter[PollCount]
	if !ok {
		return 0, errors.New("")
	}
	newValue := value + 1
	m.MetricsCounter[PollCount] = newValue
	return newValue, nil
}

//func (m *Metrics) GetMetrics() *map[string]string {
//	return &m.MetricsGauge
//}
