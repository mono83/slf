package health

import (
	"github.com/mono83/slf"
	"runtime"
	"time"
)

var now = time.Now()

// StartHealthMonitor starts goroutine, that will poll memory stats
// and send them as metrics every second
//
// Common usage: StartHealthMonitor(wd.New("", "app.").WithParams(util.HostParam()))
func StartHealthMonitor(rec slf.StatsReporter) {
	go func() {
		for {
			time.Sleep(time.Second)
			sendStats(rec)
		}
	}()
}

// sendStats takes and sends various health information
func sendStats(rec slf.StatsReporter) {
	mem := runtime.MemStats{}
	gor := runtime.NumGoroutine()
	runtime.ReadMemStats(&mem)

	rec.UpdateGauge("gcs", int64(mem.NumGC))
	rec.UpdateGauge("goroutines", int64(gor))
	rec.UpdateGauge("sys.malloc", int64(mem.Mallocs))
	rec.UpdateGauge("sys.free", int64(mem.Frees))
	rec.UpdateGauge("heap.alloc", int64(mem.HeapAlloc))
	rec.UpdateGauge("heap.inuse", int64(mem.HeapInuse))
	rec.UpdateGauge("heap.sys", int64(mem.HeapSys))
	rec.UpdateGauge("heap.objects", int64(mem.HeapObjects))
	rec.UpdateGauge("heap.nextgc", int64(mem.NextGC))
	rec.UpdateGauge("uptime", int64(time.Now().Sub(now).Seconds()))
}
