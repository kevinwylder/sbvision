import { Vector, Quaternion, qRotate, eToQ, qMultiply } from './math';

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

const triangles: Triangle[] = [
    [0, 1, 2],
    [0, 1, 3],
    [1, 3, 4],
    [3, 4, 5],
    [6, 7, 8],
    [9, 10, 11],
]

export function renderSkateboard(ctx: CanvasRenderingContext2D, rotation: Quaternion, box: [number, number, number, number]) {


    let rotated = points.map(point => qRotate(point, rotation));
    let translated = rotated.map(([x, y, z]) => [x, y, z - 10]);
    let perspective = translated.map(([x, y, z]) => [-1.5 * x / z, -1.5 * y / z, z]);
    let painters = triangles.map(([a, b, c], i) => [Math.min(perspective[a][2] + perspective[b][2] + perspective[c][2]), i])
    painters = painters.sort(([d1], [d2]) => d1 - d2);

    ctx.fillStyle = "white";
    ctx.fillRect(box[0], box[1], box[2]-box[0], box[3]-box[1]);
    let scale = 4 * Math.abs(box[0] - box[2]) / 2;
    let dx = (box[0] + box[2]) / 2;
    let dy = (box[1] + box[3]) / 2;
    painters.forEach(([_, t]) => {
        let [ p0, p1, p2 ] = triangles[t];
        let [ x0, y0 ] = perspective[p0];
        let [ x1, y1 ] = perspective[p1];
        let [ x2, y2 ] = perspective[p2];
        ctx.fillStyle = "black";
        ctx.strokeStyle = "#d8d8d8";
        ctx.beginPath();
        ctx.moveTo(x0 * scale + dx, y0 * scale + dy);
        ctx.lineTo(x1 * scale + dx, y1 * scale + dy);
        ctx.lineTo(x2 * scale + dx, y2 * scale + dy);
        ctx.closePath();
        ctx.fill();
        ctx.stroke();
    })

}

export function rotateSkateboard(dx: number, dy: number, rotation: Quaternion): Quaternion {
    dy *= 1.5; // dy should be more sensitive
    let m = Math.sqrt(dx * dx + dy * dy);

    // convert the mouse movement from euler angle to quaternion
    let delta = eToQ([ dy / m, 0, -dx / m, m / 100]);

    // left multiply the quaternion
    let newQuaternion = qMultiply(rotation, delta);

    // check if the x coordinate is positive
    let [ x ] = qRotate([1, 0, 0], newQuaternion);
    if (x > 0) {
        return newQuaternion;
    }
    return rotation;
}

export function tiltSkateboard(m: number, rotation: Quaternion): Quaternion {
    let delta = eToQ([0, 1, 0, m / 100]);
    let newQuaternion = qMultiply(rotation, delta);
    // check if the x coordinate is positive
    let [ x ] = qRotate([1, 0, 0], newQuaternion);
    if (x > 0) {
        return newQuaternion;
    }
    return rotation;
}