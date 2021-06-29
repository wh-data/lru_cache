package lru_cache

import (
	"flag"
	"runtime"
)

var (
	maxM = flag.Int("maxM", 5*1024*1024, "max stack mem (Bytes)")
	mem  *runtime.MemStats
)

func exceedMaxMem() bool {
	mem = new(runtime.MemStats)
	runtime.ReadMemStats(mem)
	return int(mem.Alloc) > *maxM
}
