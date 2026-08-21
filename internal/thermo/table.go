package thermo

// This file pins the thermodynamic tables for every species the solver can
// encounter. Each entry stores the NASA seven-coefficient polynomial for two
// temperature ranges plus the standard enthalpy of formation at 298.15 K in
// kJ/mol. The polynomial is written in the usual NASA form with
//   cp/R = a1 + a2*T + a3*T^2 + a4*T^3 + a5*T^4
// and the a6/a7 constants are kept for the Gibbs-energy evaluation used by
// the optional CO/H2 dissociation path. The values are the standard
// GRI-Mech / Burcat sets, in J/(mol K) units scaled by the gas constant.

// nasaRange holds the seven NASA coefficients for one temperature range.
type nasaRange struct {
	a1, a2, a3, a4, a5 float64 // cp/R polynomial coefficients
	a6, a7             float64 // enthalpy and entropy integration constants
}

// nasaEntry is the full thermodynamic description of one species.
type nasaEntry struct {
	name    string
	formula string
	dhf     float64 // standard enthalpy of formation at 298.15 K, kJ/mol
	low     nasaRange
	high    nasaRange
	mid     float64 // polynomial switch temperature, K
}

// speciesTable holds every species used by the combustion solver. The order
// matters for stable output: complete-combustion products first, then fuel
// species, then dissociation-only species.
var speciesTable = []nasaEntry{
	{
		name:    "CO2",
		formula: "CO2",
		dhf:     -393.522,
		mid:     1000.0,
		low: nasaRange{
			a1: 2.356773522, a2: 8.984596770e-03, a3: -7.123562690e-06,
			a4: 2.459190860e-09, a5: -1.436995480e-13,
			a6: -4.837196970e+04, a7: 9.901052220,
		},
		high: nasaRange{
			a1: 3.857460280, a2: 4.414370260e-03, a3: -2.214814040e-06,
			a4: 5.234901880e-10, a5: -4.720841640e-14,
			a6: -4.875916600e+04, a7: 2.271638060,
		},
	},
	{
		name:    "H2O",
		formula: "H2O",
		dhf:     -241.826,
		mid:     1000.0,
		low: nasaRange{
			a1: 4.198640560, a2: -2.036434100e-03, a3: 6.520402110e-06,
			a4: -5.487970620e-09, a5: 1.771978170e-12,
			a6: -3.029372670e+04, a7: -8.490322080e-01,
		},
		high: nasaRange{
			a1: 2.677037870, a2: 2.973183290e-03, a3: -7.737696490e-07,
			a4: 9.443366890e-11, a5: -4.269009590e-15,
			a6: -2.988589380e+04, a7: 6.882555710,
		},
	},
	{
		name:    "O2",
		formula: "O2",
		dhf:     0.0,
		mid:     1000.0,
		low: nasaRange{
			a1: 3.782456360, a2: -2.996734150e-03, a3: 9.847302010e-06,
			a4: -9.681295080e-09, a5: 3.243728360e-12,
			a6: -1.063943560e+03, a7: 3.657675730,
		},
		high: nasaRange{
			a1: 3.282537840, a2: 1.483087540e-03, a3: -7.579666690e-07,
			a4: 2.094705550e-10, a5: -2.167177940e-15,
			a6: -1.088457720e+03, a7: 5.453231290,
		},
	},
	{
		name:    "N2",
		formula: "N2",
		dhf:     0.0,
		mid:     1000.0,
		low: nasaRange{
			a1: 3.298677000, a2: 1.408240400e-03, a3: -3.963222000e-06,
			a4: 5.641515000e-09, a5: -2.444854000e-12,
			a6: -1.020899900e+03, a7: 3.950372000,
		},
		high: nasaRange{
			a1: 2.926640000, a2: 1.487976800e-03, a3: -5.684760000e-07,
			a4: 1.009703800e-10, a5: -6.753351000e-15,
			a6: -9.227977000e+02, a7: 5.980528000,
		},
	},
	{
		name:    "CH4",
		formula: "CH4",
		dhf:     -74.873,
		mid:     1000.0,
		low: nasaRange{
			a1: 5.149876130, a2: -1.367097880e-02, a3: 4.918005990e-05,
			a4: -4.847430260e-08, a5: 1.666939560e-11,
			a6: -1.024664760e+04, a7: -4.641303760,
		},
		high: nasaRange{
			a1: 7.485149500e-02, a2: 1.339094670e-02, a3: -5.732858090e-06,
			a4: 1.222925350e-09, a5: -1.018152300e-13,
			a6: -9.468344590e+03, a7: 6.869197250,
		},
	},
	{
		name:    "C2H4",
		formula: "C2H4",
		dhf:     52.467,
		mid:     1000.0,
		low: nasaRange{
			a1: 3.959201480, a2: -7.570522470e-03, a3: 5.709902740e-05,
			a4: -6.915887530e-08, a5: 2.698843730e-11,
			a6: 5.089775930e+03, a7: 4.097330960,
		},
		high: nasaRange{
			a1: 2.036111160, a2: 1.441054020e-02, a3: -6.019643100e-06,
			a4: 1.360261350e-09, a5: -1.252145240e-13,
			a6: 4.339493130e+03, a7: 1.056337470,
		},
	},
	{
		name:    "C2H6",
		formula: "C2H6",
		dhf:     -84.684,
		mid:     1000.0,
		low: nasaRange{
			a1: 4.291424920, a2: -5.501542700e-03, a3: 5.994238850e-05,
			a4: -7.084662850e-08, a5: 2.686857710e-11,
			a6: -1.152220550e+04, a7: 2.666823160,
		},
		high: nasaRange{
			a1: 1.071880150, a2: 2.168526690e-02, a3: -1.002560670e-05,
			a4: 2.315621340e-09, a5: -2.092921900e-13,
			a6: -1.142639320e+04, a7: 1.511561070,
		},
	},
	{
		name:    "CO",
		formula: "CO",
		dhf:     -110.527,
		mid:     1000.0,
		low: nasaRange{
			a1: 3.579533470, a2: -6.103536800e-04, a3: 1.016814330e-06,
			a4: 9.070058840e-10, a5: -9.044244880e-13,
			a6: -1.434408600e+04, a7: 3.508409280,
		},
		high: nasaRange{
			a1: 2.715185610, a2: 2.062527430e-03, a3: -9.988257710e-07,
			a4: 2.300530090e-10, a5: -2.036477160e-14,
			a6: -1.415187240e+04, a7: 7.818687720,
		},
	},
	{
		name:    "H2",
		formula: "H2",
		dhf:     0.0,
		mid:     1000.0,
		low: nasaRange{
			a1: 3.337279200, a2: -4.940247410e-05, a3: 4.994567780e-07,
			a4: -1.795663940e-10, a5: 2.002553760e-14,
			a6: -9.501589220e+02, a7: -3.205023310,
		},
		high: nasaRange{
			a1: 2.932865050, a2: 8.266080200e-04, a3: -1.464023350e-07,
			a4: 1.541004140e-11, a5: -6.888048320e-16,
			a6: -8.130655810e+02, a7: -1.024328650,
		},
	},
}

// tableIndex maps a species name to its position in speciesTable.
var tableIndex = func() map[string]int {
	m := make(map[string]int, len(speciesTable))
	for i, e := range speciesTable {
		m[e.name] = i
	}
	return m
}()
