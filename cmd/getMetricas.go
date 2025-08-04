package main

import (
	"runtime"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// GetMetrixFromOS получает антайм метрики из runtime.MemStats
func GetMetrixFromOS() (gaugeMap map[string]uint64, err error) {

	v, err := mem.VirtualMemory() 
	if err != nil {
		return nil, err
	}
	
	cc, err := cpu.Counts(true)
	if err != nil {
		return nil, err
	}
	
	var mS runtime.MemStats

	runtime.ReadMemStats(&mS)

	gaugeMap = map[string]uint64{
		"Alloc":           (mS.Alloc),
		"BuckHashSys":     (mS.BuckHashSys),
		"Frees":           (mS.Frees),
		"GCCPUFraction":   uint64(mS.GCCPUFraction),
		"GCSys":           (mS.GCSys),
		"HeapAlloc":       (mS.HeapAlloc),
		"HeapIdle":        (mS.HeapIdle),
		"HeapInuse":       (mS.HeapInuse),
		"HeapObjects":     (mS.HeapObjects),
		"HeapReleased":    (mS.HeapReleased),
		"HeapSys":         (mS.HeapSys),
		"LastGC":          (mS.LastGC),
		"Lookups":         (mS.Lookups),
		"MCacheInuse":     (mS.MCacheInuse),
		"MCacheSys":       (mS.MCacheSys),
		"MSpanInuse":      (mS.MSpanInuse),
		"MSpanSys":        (mS.MSpanSys),
		"Mallocs":         (mS.Mallocs),
		"NextGC":          (mS.NextGC),
		"NumForcedGC":     uint64(mS.NumForcedGC),
		"NumGC":           uint64(mS.NumGC),
		"OtherSys":        (mS.OtherSys),
		"PauseTotalNs":    (mS.PauseTotalNs),
		"StackInuse":      (mS.StackInuse),
		"StackSys":        (mS.StackSys),
		"Sys":             (mS.Sys),
		"TotalAlloc":      (mS.TotalAlloc),
		"TotalMemory":     (v.Total),
		"FreeMemory":      (v.Free),
		"CPUutilization1": uint64(cc),
	}
	return
}
