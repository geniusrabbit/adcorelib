package gometrics

import (
	"context"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	memAlloc = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "go_memory_alloc_bytes",
		Help: "Current number of bytes allocated and still in use",
	})
	memSys = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "go_memory_sys_bytes",
		Help: "Total bytes of memory obtained from the OS",
	})
	memHeapIdle = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "go_memory_heap_idle_bytes",
		Help: "Bytes in idle (unused) spans",
	})
	memHeapInuse = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "go_memory_heap_inuse_bytes",
		Help: "Bytes in in-use spans",
	})
	memGCCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "go_gc_total_count",
		Help: "Number of completed GC cycles",
	})
	lastPause = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "go_gc_last_pause_ns",
		Help: "The most recent GC pause duration in nanoseconds",
	})
	heapObjects = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "go_heap_objects",
		Help: "Number of allocated heap objects",
	})
	gcCpuFraction = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "go_gc_cpu_fraction",
		Help: "Fraction of CPU time used by GC",
	})
	numGoroutines = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "go_goroutines_count",
		Help: "Number of goroutines",
	})
	cgoCalls = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "go_cgo_calls_total",
		Help: "Total number of cgo calls",
	})
)

func TrackRuntimeMetrics(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var lastGCCount uint32

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			var stats runtime.MemStats
			runtime.ReadMemStats(&stats)

			memAlloc.Set(float64(stats.Alloc))
			memSys.Set(float64(stats.Sys))
			memHeapIdle.Set(float64(stats.HeapIdle))
			memHeapInuse.Set(float64(stats.HeapInuse))
			heapObjects.Set(float64(stats.HeapObjects))
			lastPause.Set(float64(stats.PauseNs[(stats.NumGC+255)%256]))
			gcCpuFraction.Set(stats.GCCPUFraction)

			deltaGC := stats.NumGC - lastGCCount
			if deltaGC > 0 {
				memGCCount.Add(float64(deltaGC))
				lastGCCount = stats.NumGC
			}

			numGoroutines.Set(float64(runtime.NumGoroutine()))
			cgoCalls.Set(float64(runtime.NumCgoCall()))
		}
	}
}
