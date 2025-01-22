import { Vector3D } from "./types";

// src/math.ts
export class Vector3DMath {
  static magnitude(v: Vector3D): number {
    return Math.sqrt(v.x * v.x + v.y * v.y + v.z * v.z);
  }

  static normalize(v: Vector3D): Vector3D {
    const mag = this.magnitude(v);
    return {
      x: v.x / mag,
      y: v.y / mag,
      z: v.z / mag,
    };
  }

  static subtract(a: Vector3D, b: Vector3D): Vector3D {
    return {
      x: a.x - b.x,
      y: a.y - b.y,
      z: a.z - b.z,
    };
  }

  static crossProduct(a: Vector3D, b: Vector3D): Vector3D {
    return {
      x: a.y * b.z - a.z * b.y,
      y: a.z * b.x - a.x * b.z,
      z: a.x * b.y - a.y * b.x,
    };
  }

  static scale(v: Vector3D, scalar: number): Vector3D {
    return {
      x: v.x * scalar,
      y: v.y * scalar,
      z: v.z * scalar,
    };
  }
}
