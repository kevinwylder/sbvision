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

type Triangle = [ number, number, number ] // these are "pointers" in the points array

const points: Vector[] = [
    [2, .5, 0],
    [2, -.5, 0],
    [2.5, 0, .25],
    [-2, .5, 0],
    [-2, -.5, 0],
    [-2.5, 0, .25],
    [1.75, 0, 0],
    [1.75, -.5, -.5],
    [1.75, .5, -.5],
    [-1.75, 0, 0],
    [-1.75, -.5, -.5],
    [-1.75, .5, -.5]
]

const board: Triangle[] = [
    [0, 1, 2],
    [0, 1, 3],
    [1, 3, 4],
    [3, 4, 5],
];

const trucks: Triangle[] = [
    [6, 7, 8],
    [9, 10, 11],
]

export function renderSkateboard(ctx: CanvasRenderingContext2D, rotation: Quaternion, box: [number, number, number, number]) {
    let xDepth  = qRotate([1, 0, 0], rotation)[0];
    let zDepth = qRotate([0, 0, 1], rotation)[2];

    let rotated = points.map(point => qRotate(point, rotation));
    let translated = rotated.map(([x, y, z]) => [x, y, z - 10]);
    let perspective = translated.map(([x, y, z]) => [-1.5 * x / z, -1.5 * y / z, z]);

    ctx.lineCap = "round";
    ctx.lineWidth = 2;
    ctx.fillStyle = "white";
    ctx.strokeStyle = "black";
    ctx.fillRect(box[0], box[1], box[2]-box[0], box[3]-box[1]);
    let scale = 2 * Math.abs(box[0] - box[2]) / 2;
    let dx = (box[0] + box[2]) / 2;
    let dy = (box[1] + box[3]) / 2;

    if (zDepth >= 0) {
        trucks.forEach(([p0, p1, p2]) => {
            let [ x0, y0 ] = perspective[p0];
            let [ x1, y1 ] = perspective[p1];
            let [ x2, y2 ] = perspective[p2];
            ctx.fillStyle = "#c9c9c9";
            ctx.beginPath();
            ctx.moveTo(x0 * scale + dx, y0 * scale + dy);
            ctx.lineTo(x1 * scale + dx, y1 * scale + dy);
            ctx.lineTo(x2 * scale + dx, y2 * scale + dy);
            ctx.closePath();
            ctx.stroke();
            ctx.fill();
        })
    }

    ((xDepth > 0) ? board : board.reverse()).forEach(([p0, p1, p2]) => {
        let [ x0, y0 ] = perspective[p0];
        let [ x1, y1 ] = perspective[p1];
        let [ x2, y2 ] = perspective[p2];
        ctx.fillStyle = (zDepth > 0) ? "black" : "red";
        ctx.strokeStyle = "#d8d8d8";
        ctx.beginPath();
        ctx.moveTo(x0 * scale + dx, y0 * scale + dy);
        ctx.lineTo(x1 * scale + dx, y1 * scale + dy);
        ctx.lineTo(x2 * scale + dx, y2 * scale + dy);
        ctx.closePath();
        ctx.stroke();
        ctx.fill();
    })

    if (zDepth < 0) {
        trucks.forEach(([p0, p1, p2]) => {
            let [ x0, y0 ] = perspective[p0];
            let [ x1, y1 ] = perspective[p1];
            let [ x2, y2 ] = perspective[p2];
            ctx.fillStyle = "#C9C9C9";
            ctx.beginPath();
            ctx.moveTo(x0 * scale + dx, y0 * scale + dy);
            ctx.lineTo(x1 * scale + dx, y1 * scale + dy);
            ctx.lineTo(x2 * scale + dx, y2 * scale + dy);
            ctx.closePath();
            ctx.stroke();
            ctx.fill();
        })
    }

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