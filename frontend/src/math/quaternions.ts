import { Quaternion } from ".";
import { EulerAngle, Vector } from ".";

function qConjugate(a: Quaternion): Quaternion {
    return [a[0], -a[1], -a[2], -a[3]];
}

export function eToQ(angle: EulerAngle): Quaternion {
    let s = Math.sin(angle[3] / 2);
    return [
        Math.cos(angle[3] / 2),
        angle[0] * s,
        angle[1] * s,
        angle[2] * s
    ]
}

export function qMultiply(a: Quaternion, b: Quaternion): Quaternion {
    return [
        a[0] * b[0] - a[1] * b[1] - a[2] * b[2] - a[3] * b[3],
        a[0] * b[1] + a[1] * b[0] + a[2] * b[3] - a[3] * b[2],
        a[0] * b[2] - a[1] * b[3] + a[2] * b[0] + a[3] * b[1],
        a[0] * b[3] + a[1] * b[2] - a[2] * b[1] + a[3] * b[0]
    ]
}

export function qDiff(a: Quaternion, b: Quaternion): Quaternion {
    return qMultiply(qConjugate(b), a);
}

export function qRotate(pos: Vector, by: Quaternion): Vector { 
    let d2 = by.map(t => t * t).reduce((p, n) => p + n, 0)
    let [ _, i, j, k ] = qMultiply(qMultiply(by, [0, pos[0], pos[1], pos[2]]), qConjugate(by))
    return [i / d2, j / d2, k / d2];
}

export function rotateSkateboard(dx: number, dy: number, rotation: Quaternion): Quaternion {
    dy *= 1.5; // dy should be more sensitive
    let m = Math.sqrt(dx * dx + dy * dy);
    if (m == 0) {
        return rotation;
    }

    // convert the mouse movement from euler angle to quaternion
    let delta = eToQ([ dy / m, 0, -dx / m, m / 100]);

    // left multiply the quaternion
    return qMultiply(rotation, delta);
}

export function tiltSkateboard(m: number, rotation: Quaternion): Quaternion {
    if (m == 0) {
        return rotation;
    }
    let delta = eToQ([0, 1, 0, m / 100]);
    return qMultiply(rotation, delta);
}

// returns a random unit quaternion
export function qRandom(): Quaternion {
    while (true) {
        let q0 = Math.random();
        let q1 = Math.random();
        let q2 = Math.random();
        let q3 = Math.random();
        let d = q0 * q0 + q1 * q1 + q2 * q2 + q3 * q3;
        if (d < 0) {
            d = Math.sqrt(d);
            return [ q0 / d, q1 / d, q2 / d, q3 / d ];
        }
    }
}
