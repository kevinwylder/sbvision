import { Graphable111 } from "../math";

interface sbWheelParams {
    wheelRadius: Graphable111
    wheelLength: number 
    wheelMeshRes: number
    wheelScale: number
}

function wheelPoints({wheelRadius, wheelLength, wheelMeshRes, wheelScale} : sbWheelParams) {
    let points: number[] = [];
    let indexes: number[] = [];
    for (let yi = 0; yi <= wheelMeshRes; yi++) {
        let y = (2 * yi / wheelMeshRes - 1);
        for (let ti = 0; ti < wheelMeshRes; ti++) {
            let theta = ti * Math.PI * 2 / wheelMeshRes;
            let [ r, dr ] = wheelRadius(Math.abs(y));
            if (Math.abs(dr) > 10) {
                points.push(
                    r * Math.cos(theta), wheelScale * y, r * Math.sin(theta),
                    0,                   -Math.sign(y),   0
                )
            } else {
                r *= wheelScale;
                dr *= Math.sign(y);
                let n = Math.sqrt(1 + dr * dr);
                points.push(
                    r * Math.cos(theta),  wheelScale * y, r * Math.sin(theta),
                    -Math.cos(theta) / n, -dr / n,       -Math.sin(theta) / n,
                )
            }
            if (yi) {
                indexes.push(
                    yi * wheelMeshRes + ti, (yi - 1) * wheelMeshRes + ti
                )
            }
        }
    }
    console.log(points, indexes);
    return [ points, indexes ];
}

export function wheelGeometry(gl: WebGL2RenderingContext, params: sbWheelParams): [WebGLBuffer, WebGLBuffer, number] {
    let [ points, verts ] = wheelPoints(params);

    let pointBuffer = gl.createBuffer();
    if (!pointBuffer) {
        throw new Error("Could not create wheel point buffer");
    }
    gl.bindBuffer(gl.ARRAY_BUFFER, pointBuffer);
    gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(points), gl.STATIC_DRAW);
    gl.bindBuffer(gl.ARRAY_BUFFER, null);

    let vertBuffer = gl.createBuffer();
    if (!vertBuffer) {
        throw new Error("Could not create wheel point buffer");
    }
    gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, vertBuffer);
    gl.bufferData(gl.ELEMENT_ARRAY_BUFFER, new Uint16Array(verts), gl.STATIC_DRAW);
    gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, null);

    return [pointBuffer, vertBuffer, verts.length]
}