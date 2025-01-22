import { CONSTANTS } from "./constants";
import { Vector3DMath } from "./math";
import { Vector3D } from "./types";

// src/physics.ts
export class ElectromagnetismCalculator {
  // 1. Coulomb's Law
  static calculateElectricForce(
    charge1: number,
    position1: Vector3D,
    charge2: number,
    position2: Vector3D,
  ): Vector3D {
    const r = Vector3DMath.subtract(position2, position1);
    const rMagnitude = Vector3DMath.magnitude(r);
    const rUnit = Vector3DMath.normalize(r);

    // F = k * q1 * q2 / r^2
    const forceMagnitude = CONSTANTS.K * charge1 * charge2 /
      (rMagnitude * rMagnitude);
    return Vector3DMath.scale(rUnit, forceMagnitude);
  }

  // 2. Magnetic Field from Current (Biot-Savart Law)
  static calculateMagneticField(
    current: number,
    wirePath: Vector3D,
    observationPoint: Vector3D,
  ): Vector3D {
    const r = Vector3DMath.subtract(observationPoint, wirePath);
    const rMagnitude = Vector3DMath.magnitude(r);

    // dB = (μ₀/4π) * (I dl × r̂)/r²
    const factor = (CONSTANTS.MU_0 * current) /
      (4 * Math.PI * rMagnitude * rMagnitude);
    const crossProduct = Vector3DMath.crossProduct(
      wirePath,
      Vector3DMath.normalize(r),
    );
    return Vector3DMath.scale(crossProduct, factor);
  }

  // 3. Lorentz Force
  static calculateLorentzForce(
    charge: number,
    velocity: Vector3D,
    electricField: Vector3D,
    magneticField: Vector3D,
  ): Vector3D {
    // F = q(E + v × B)
    const magneticForce = Vector3DMath.crossProduct(velocity, magneticField);
    const scaledMagneticForce = Vector3DMath.scale(magneticForce, charge);
    const electricForce = Vector3DMath.scale(electricField, charge);

    return {
      x: electricForce.x + scaledMagneticForce.x,
      y: electricForce.y + scaledMagneticForce.y,
      z: electricForce.z + scaledMagneticForce.z,
    };
  }

  // 4. Induced EMF (Faraday's Law)
  static calculateInducedEMF(
    fluxChange: number,
    timeInterval: number,
    turns: number,
  ): number {
    // ε = -N * ΔΦ/Δt
    return turns * (-fluxChange / timeInterval);
  }

  // 5. Cyclotron Frequency
  static calculateCyclotronFrequency(
    charge: number,
    magneticField: number,
    mass: number,
  ): number {
    // f = |q|B/(2πm)
    return Math.abs(charge * magneticField) / (2 * Math.PI * mass);
  }

  // Electric potential energy between point charges
  static calculateElectricPotentialEnergy(
    charge1: number,
    charge2: number,
    distance: number,
  ): number {
    // U = k * q1 * q2 / r
    return CONSTANTS.K * charge1 * charge2 / distance;
  }

  // Magnetic flux through a surface
  static calculateMagneticFlux(
    magneticField: Vector3D,
    area: number,
    angle: number,
  ): number {
    // Φ = B * A * cos(θ)
    return Vector3DMath.magnitude(magneticField) * area * Math.cos(angle);
  }

  // Capacitor energy storage
  static calculateCapacitorEnergy(
    capacitance: number,
    voltage: number,
  ): number {
    // U = (1/2) * C * V^2
    return 0.5 * capacitance * voltage * voltage;
  }

  // Inductance calculation for a solenoid
  static calculateSolenoidInductance(
    turns: number,
    length: number,
    area: number,
  ): number {
    // L = (μ₀ * N^2 * A) / l
    return (CONSTANTS.MU_0 * turns * turns * area) / length;
  }

  // RC circuit time constant
  static calculateRCTimeConstant(
    resistance: number,
    capacitance: number,
  ): number {
    // τ = R * C
    return resistance * capacitance;
  }
}
