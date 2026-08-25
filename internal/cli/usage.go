package cli

const usage = `adiab-flame: constant-pressure adiabatic flame temperature.

Given a fuel, an equivalence ratio and an inlet temperature in a JSON file,
adiab-flame balances the C/H/O/N atoms and iterates the product temperature
until the product enthalpy equals the reactant enthalpy. It prints the flame
temperature, the product mole fractions and the atom-balance residuals.

Supported fuels: CH4, C2H4, C2H6.
Air model: 21% O2 / 79% N2 by mole.
Products: CO2, H2O, O2, N2, plus unburned fuel for rich mixtures and CO/H2
when the optional dissociation equilibrium is enabled.

usage:
  adiab-flame tad <config.json> [--dissoc]
  adiab-flame checks <fuel> [--inlet K]
  adiab-flame fuels
  adiab-flame help

A config JSON looks like:

  {
    "fuel": "CH4",
    "equivalence_ratio": 1.0,
    "inlet_temperature_k": 298.15,
    "dissociation": false
  }

Illegal inputs (phi <= 0, inlet temperature <= 0, unknown fuel, unknown JSON
fields, missing files, a temperature iteration that fails to converge) are
reported on stderr and exit non-zero.
`
