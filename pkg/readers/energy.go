package readers

// Energy : the structure that hold the energy measurements
// Pkg => Package,
// PP0 => Core,
// PP1 => Uncore (L3 cache, integrated GPU if present),
// DRAM => Ram,
// PSys => Platform (if available)
type Energy struct {
	Pkg, PP0, PP1, DRAM, PSys float64
}

func (e Energy) Add(e2 Energy) Energy {
	return Energy{
		Pkg:  e.Pkg + e2.Pkg,
		PP0:  e.PP0 + e2.PP0,
		PP1:  e.PP1 + e2.PP1,
		DRAM: e.DRAM + e2.DRAM,
		PSys: e.PSys + e2.PSys,
	}
}

func (e Energy) Sub(e2 Energy) Energy {
	return Energy{
		Pkg:  e.Pkg - e2.Pkg,
		PP0:  e.PP0 - e2.PP0,
		PP1:  e.PP1 - e2.PP1,
		DRAM: e.DRAM - e2.DRAM,
		PSys: e.PSys - e2.PSys,
	}
}

type Measurement map[int64]map[int]Energy

func (m Measurement) Delta(m2 Measurement) Measurement {
	m3 := Measurement{}

	for pkgId, cores := range m {
		if _, exists := m3[pkgId]; !exists {
			coreMap := make(map[int]Energy)
			m3[pkgId] = coreMap
		}
		for coreId, core := range cores {
			if _, exists := m3[pkgId][coreId]; !exists {
				m3[pkgId][coreId] = core.Sub(m2[pkgId][coreId])
			}
		}
	}

	return m3
}

func (m Measurement) DeltaSum(m2 Measurement) Measurement {
	m3 := Measurement{}

	for pkgId, cores := range m {
		if _, exists := m3[pkgId]; !exists {
			coreMap := make(map[int]Energy)
			m3[pkgId] = coreMap
			m3[pkgId][0] = Energy{}
		}
		for coreId, core := range cores {
			m3[pkgId][0] = core.Sub(m2[pkgId][coreId])
		}
	}

	return m3
}

type Units struct {
	Power, Time, CpuEnergy, DramEnergy float64
}
