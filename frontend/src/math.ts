
export type Quaternion = [ number, number, number, number ];

export type EulerAngle = [ number, number, number, number ];

export type Vector = [ number, number, number ];

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

export function qRotate(pos: Vector, by: Quaternion): Vector { 
    let d2 = by.map(t => t * t).reduce((p, n) => p + n, 0)
    let [ _, i, j, k ] = qMultiply(qMultiply(by, [0, pos[0], pos[1], pos[2]]), qConjugate(by))
    return [i / d2, j / d2, k / d2];
}