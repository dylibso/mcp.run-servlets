// src/main.ts
import { Vector3DMath } from "./math";
import {
  CallToolRequest,
  CallToolResult,
  ContentType,
  ListToolsResult,
} from "./pdk";
import { ElectromagnetismCalculator } from "./physics";

export function callImpl(input: CallToolRequest): CallToolResult {
  const args = input.params.arguments;

  try {
    switch (input.params.name) {
      case "coulomb_force": {
        const force = ElectromagnetismCalculator.calculateElectricForce(
          args.charge1,
          args.position1,
          args.charge2,
          args.position2,
        );
        return {
          content: [{
            type: ContentType.Text,
            text: JSON.stringify(
              {
                force,
                magnitude: Vector3DMath.magnitude(force),
                unit: "Newtons",
              },
              null,
              2,
            ),
          }],
        };
      }

      case "magnetic_field": {
        const field = ElectromagnetismCalculator.calculateMagneticField(
          args.current,
          args.wirePath,
          args.observationPoint,
        );
        return {
          content: [{
            type: ContentType.Text,
            text: JSON.stringify(
              {
                field,
                magnitude: Vector3DMath.magnitude(field),
                unit: "Tesla",
              },
              null,
              2,
            ),
          }],
        };
      }

      case "lorentz_force": {
        const force = ElectromagnetismCalculator.calculateLorentzForce(
          args.charge,
          args.velocity,
          args.electricField,
          args.magneticField,
        );
        return {
          content: [{
            type: ContentType.Text,
            text: JSON.stringify(
              {
                force,
                magnitude: Vector3DMath.magnitude(force),
                unit: "Newtons",
              },
              null,
              2,
            ),
          }],
        };
      }

      case "induced_emf": {
        const emf = ElectromagnetismCalculator.calculateInducedEMF(
          args.fluxChange,
          args.timeInterval,
          args.turns || 1,
        );
        return {
          content: [{
            type: ContentType.Text,
            text: JSON.stringify(
              {
                emf,
                unit: "Volts",
              },
              null,
              2,
            ),
          }],
        };
      }

      case "cyclotron_frequency": {
        const frequency = ElectromagnetismCalculator
          .calculateCyclotronFrequency(
            args.charge,
            args.magneticField,
            args.mass,
          );
        return {
          content: [{
            type: ContentType.Text,
            text: JSON.stringify(
              {
                frequency,
                unit: "Hertz",
              },
              null,
              2,
            ),
          }],
        };
      }

      case "electric_potential_energy": {
        const energy = ElectromagnetismCalculator
          .calculateElectricPotentialEnergy(
            args.charge1,
            args.charge2,
            args.distance,
          );
        return {
          content: [{
            type: ContentType.Text,
            text: JSON.stringify(
              {
                energy,
                unit: "Joules",
              },
              null,
              2,
            ),
          }],
        };
      }

      case "magnetic_flux": {
        const flux = ElectromagnetismCalculator.calculateMagneticFlux(
          args.magneticField,
          args.area,
          args.angle,
        );
        return {
          content: [{
            type: ContentType.Text,
            text: JSON.stringify(
              {
                flux,
                unit: "Weber",
              },
              null,
              2,
            ),
          }],
        };
      }

      case "capacitor_energy": {
        const energy = ElectromagnetismCalculator.calculateCapacitorEnergy(
          args.capacitance,
          args.voltage,
        );
        return {
          content: [{
            type: ContentType.Text,
            text: JSON.stringify(
              {
                energy,
                unit: "Joules",
              },
              null,
              2,
            ),
          }],
        };
      }

      case "solenoid_inductance": {
        const inductance = ElectromagnetismCalculator
          .calculateSolenoidInductance(
            args.turns,
            args.length,
            args.area,
          );
        return {
          content: [{
            type: ContentType.Text,
            text: JSON.stringify(
              {
                inductance,
                unit: "Henry",
              },
              null,
              2,
            ),
          }],
        };
      }

      case "rc_time_constant": {
        const timeConstant = ElectromagnetismCalculator.calculateRCTimeConstant(
          args.resistance,
          args.capacitance,
        );
        return {
          content: [{
            type: ContentType.Text,
            text: JSON.stringify(
              {
                timeConstant,
                unit: "Seconds",
              },
              null,
              2,
            ),
          }],
        };
      }

      default:
        return {
          content: [{
            type: ContentType.Text,
            text: `Unknown tool: ${input.params.name}`,
          }],
          isError: true,
        };
    }
  } catch (error: any) {
    return {
      content: [{
        type: ContentType.Text,
        text: `Calculation error: ${error.message}`,
      }],
      isError: true,
    };
  }
}

export function describeImpl(): ListToolsResult {
  return {
    tools: [
      {
        name: "coulomb_force",
        description:
          "Calculate the electrostatic force between two point charges using Coulomb's law",
        inputSchema: {
          type: "object",
          properties: {
            charge1: {
              type: "number",
              description: "First charge in Coulombs",
            },
            position1: {
              type: "object",
              description: "Position of first charge",
              properties: {
                x: { type: "number" },
                y: { type: "number" },
                z: { type: "number" },
              },
              required: ["x", "y", "z"],
            },
            charge2: {
              type: "number",
              description: "Second charge in Coulombs",
            },
            position2: {
              type: "object",
              description: "Position of second charge",
              properties: {
                x: { type: "number" },
                y: { type: "number" },
                z: { type: "number" },
              },
              required: ["x", "y", "z"],
            },
          },
          required: ["charge1", "position1", "charge2", "position2"],
        },
      },
      {
        name: "magnetic_field",
        description:
          "Calculate the magnetic field due to a current element using the Biot-Savart law",
        inputSchema: {
          type: "object",
          properties: {
            current: {
              type: "number",
              description: "Current in Amperes",
            },
            wirePath: {
              type: "object",
              description: "Current element vector",
              properties: {
                x: { type: "number" },
                y: { type: "number" },
                z: { type: "number" },
              },
              required: ["x", "y", "z"],
            },
            observationPoint: {
              type: "object",
              description: "Point where field is calculated",
              properties: {
                x: { type: "number" },
                y: { type: "number" },
                z: { type: "number" },
              },
              required: ["x", "y", "z"],
            },
          },
          required: ["current", "wirePath", "observationPoint"],
        },
      },
      {
        name: "lorentz_force",
        description:
          "Calculate the Lorentz force on a charged particle in electromagnetic fields",
        inputSchema: {
          type: "object",
          properties: {
            charge: {
              type: "number",
              description: "Charge in Coulombs",
            },
            velocity: {
              type: "object",
              description: "Particle velocity",
              properties: {
                x: { type: "number" },
                y: { type: "number" },
                z: { type: "number" },
              },
              required: ["x", "y", "z"],
            },
            electricField: {
              type: "object",
              description: "Electric field vector",
              properties: {
                x: { type: "number" },
                y: { type: "number" },
                z: { type: "number" },
              },
              required: ["x", "y", "z"],
            },
            magneticField: {
              type: "object",
              description: "Magnetic field vector",
              properties: {
                x: { type: "number" },
                y: { type: "number" },
                z: { type: "number" },
              },
              required: ["x", "y", "z"],
            },
          },
          required: ["charge", "velocity", "electricField", "magneticField"],
        },
      },
      {
        name: "induced_emf",
        description:
          "Calculate the induced EMF using Faraday's law of electromagnetic induction",
        inputSchema: {
          type: "object",
          properties: {
            fluxChange: {
              type: "number",
              description: "Change in magnetic flux in Weber (Wb)",
            },
            timeInterval: {
              type: "number",
              description: "Time interval in seconds",
            },
            turns: {
              type: "number",
              description: "Number of turns in the coil",
            },
          },
          required: ["fluxChange", "timeInterval"],
        },
      },
      {
        name: "cyclotron_frequency",
        description:
          "Calculate the cyclotron frequency for a charged particle in a magnetic field",
        inputSchema: {
          type: "object",
          properties: {
            charge: {
              type: "number",
              description: "Particle charge in Coulombs",
            },
            magneticField: {
              type: "number",
              description: "Magnetic field strength in Tesla",
            },
            mass: {
              type: "number",
              description: "Particle mass in kilograms",
            },
          },
          required: ["charge", "magneticField", "mass"],
        },
      },
      {
        name: "electric_potential_energy",
        description:
          "Calculate the electric potential energy between two point charges",
        inputSchema: {
          type: "object",
          properties: {
            charge1: {
              type: "number",
              description: "First charge in Coulombs",
            },
            charge2: {
              type: "number",
              description: "Second charge in Coulombs",
            },
            distance: {
              type: "number",
              description: "Distance between charges in meters",
            },
          },
          required: ["charge1", "charge2", "distance"],
        },
      },
      {
        name: "magnetic_flux",
        description: "Calculate the magnetic flux through a surface",
        inputSchema: {
          type: "object",
          properties: {
            magneticField: {
              type: "object",
              description: "Magnetic field vector",
              properties: {
                x: { type: "number" },
                y: { type: "number" },
                z: { type: "number" },
              },
              required: ["x", "y", "z"],
            },
            area: {
              type: "number",
              description: "Surface area in square meters",
            },
            angle: {
              type: "number",
              description: "Angle between field and surface normal in radians",
            },
          },
          required: ["magneticField", "area", "angle"],
        },
      },
      {
        name: "capacitor_energy",
        description: "Calculate the energy stored in a capacitor",
        inputSchema: {
          type: "object",
          properties: {
            capacitance: {
              type: "number",
              description: "Capacitance in Farads",
            },
            voltage: {
              type: "number",
              description: "Voltage across capacitor in Volts",
            },
          },
          required: ["capacitance", "voltage"],
        },
      },
      {
        name: "solenoid_inductance",
        description: "Calculate the inductance of a solenoid",
        inputSchema: {
          type: "object",
          properties: {
            turns: {
              type: "number",
              description: "Number of turns in the solenoid",
            },
            length: {
              type: "number",
              description: "Length of solenoid in meters",
            },
            area: {
              type: "number",
              description: "Cross-sectional area in square meters",
            },
          },
          required: ["turns", "length", "area"],
        },
      },
      {
        name: "rc_time_constant",
        description: "Calculate the time constant of an RC circuit",
        inputSchema: {
          type: "object",
          properties: {
            resistance: {
              type: "number",
              description: "Resistance in Ohms",
            },
            capacitance: {
              type: "number",
              description: "Capacitance in Farads",
            },
          },
          required: ["resistance", "capacitance"],
        },
      },
    ],
  };
}
