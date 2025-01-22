# Electromagnetic Physics Calculator

This module provides a collection of tools for calculating various
electromagnetic phenomena. Each function handles specific electromagnetic
calculations and returns results in standard SI units.

## Available Functions

### 1. Coulomb Force (`coulomb_force`)

Calculates the electrostatic force between two point charges using Coulomb's
law.

- **Input Parameters:**
  - `charge1`: First charge in Coulombs
  - `position1`: Position vector of first charge (x, y, z)
  - `charge2`: Second charge in Coulombs
  - `position2`: Position vector of second charge (x, y, z)
- **Output:** Force vector in Newtons

### 2. Magnetic Field (`magnetic_field`)

Calculates the magnetic field due to a current element using the Biot-Savart
law.

- **Input Parameters:**
  - `current`: Current in Amperes
  - `wirePath`: Current element vector (x, y, z)
  - `observationPoint`: Point where field is calculated (x, y, z)
- **Output:** Magnetic field vector in Tesla

### 3. Lorentz Force (`lorentz_force`)

Calculates the Lorentz force on a charged particle in electromagnetic fields.

- **Input Parameters:**
  - `charge`: Charge in Coulombs
  - `velocity`: Particle velocity vector (x, y, z)
  - `electricField`: Electric field vector (x, y, z)
  - `magneticField`: Magnetic field vector (x, y, z)
- **Output:** Force vector in Newtons

### 4. Induced EMF (`induced_emf`)

Calculates the induced EMF using Faraday's law of electromagnetic induction.

- **Input Parameters:**
  - `fluxChange`: Change in magnetic flux in Weber (Wb)
  - `timeInterval`: Time interval in seconds
  - `turns`: Number of turns in the coil (optional, defaults to 1)
- **Output:** EMF in Volts

### 5. Cyclotron Frequency (`cyclotron_frequency`)

Calculates the cyclotron frequency for a charged particle in a magnetic field.

- **Input Parameters:**
  - `charge`: Particle charge in Coulombs
  - `magneticField`: Magnetic field strength in Tesla
  - `mass`: Particle mass in kilograms
- **Output:** Frequency in Hertz

### 6. Electric Potential Energy (`electric_potential_energy`)

Calculates the electric potential energy between two point charges.

- **Input Parameters:**
  - `charge1`: First charge in Coulombs
  - `charge2`: Second charge in Coulombs
  - `distance`: Distance between charges in meters
- **Output:** Energy in Joules

### 7. Magnetic Flux (`magnetic_flux`)

Calculates the magnetic flux through a surface.

- **Input Parameters:**
  - `magneticField`: Magnetic field vector (x, y, z) in Tesla
  - `area`: Surface area in square meters
  - `angle`: Angle between field and surface normal in radians
- **Output:** Magnetic flux in Weber

### 8. Capacitor Energy (`capacitor_energy`)

Calculates the energy stored in a capacitor.

- **Input Parameters:**
  - `capacitance`: Capacitance in Farads
  - `voltage`: Voltage across capacitor in Volts
- **Output:** Energy in Joules

### 9. Solenoid Inductance (`solenoid_inductance`)

Calculates the inductance of a solenoid.

- **Input Parameters:**
  - `turns`: Number of turns in the solenoid
  - `length`: Length of solenoid in meters
  - `area`: Cross-sectional area in square meters
- **Output:** Inductance in Henry

### 10. RC Time Constant (`rc_time_constant`)

Calculates the time constant of an RC circuit.

- **Input Parameters:**
  - `resistance`: Resistance in Ohms
  - `capacitance`: Capacitance in Farads
- **Output:** Time constant in Seconds

## Constants Used

- Vacuum permittivity (ε₀): 8.854 × 10⁻¹² F/m
- Vacuum permeability (μ₀): 4π × 10⁻⁷ H/m
- Coulomb constant (k): 1/(4πε₀)

## Error Handling

All functions include error handling and will return appropriate error messages
if:

- Invalid parameters are provided
- Calculations result in errors
- Unknown tool names are requested

## Response Format

All successful calculations return a JSON response containing:

- The calculated value (scalar or vector)
- The magnitude (for vector quantities)
- The appropriate SI unit
