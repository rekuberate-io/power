package readers

// Energy : the structure that hold the energy measurements
// Pkg => Package,
// PP0 => Core,
// PP1 => Uncore (L3 cache, integrated GPU if present),
// DRAM => Ram,
// PSys => Platform (if available)
type Energy struct {
	Pkg, PP0, PP1, DRAM, PSys uint64
}

type Units struct {
	Power, Time, CpuEnergy, DramEnergy float64
}
