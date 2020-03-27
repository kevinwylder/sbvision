
interface sbWheelParams {
    wheelRadius: number
    wheelLength: number 
    wheelMeshRes: number
    wheelScale: number
    wheelNormTaper: number
}

function wheelCap({wheelRadius, wheelMeshRes, wheelNormTaper} : sbWheelParams) {
    let points: number[] = []
    for (let side = -1; side <= 1; side += 2) {
        points.push(
            0, 0, 0,
            0, side, 0
        );
        let r = wheelNormTaper;
        let n = Math.sqrt(1 + r * r);
        r /= -n;
        for (let ti = 0; ti <= wheelMeshRes; ti++) {
            let theta = ti * Math.PI * 2 / wheelMeshRes;
            points.push(
                wheelRadius * Math.cos(theta), 0,        wheelRadius * Math.sin(theta),
                r * Math.cos(theta),           side / n, r * Math.sin(theta)
            )
        }
    }
    return points;
}

function wheelPoints({wheelRadius, wheelLength, wheelMeshRes, wheelNormTaper} : sbWheelParams) {
    let points: number[] = [];
    let indexes: number[] = [];
    for (let yi = -1; yi <= 1; yi += 2) {
        let dy = -wheelNormTaper;
        let n = Math.sqrt(1 + dy * dy)
        dy /= n * yi;
        let y = yi * wheelLength;
        for (let ti = 0; ti <= wheelMeshRes; ti++) {
            let theta = ti * Math.PI * 2 / wheelMeshRes;
            points.push(
                wheelRadius * Math.cos(theta), y,      wheelRadius * Math.sin(theta),
                -Math.cos(theta) / n,          dy / n, -Math.sin(theta) / n
            )
            if (yi == 1) {
                indexes.push(
                    ti, ti + wheelMeshRes + 1,
                )
            }
        }
    }
    return [ points, indexes ];
}

export function wheelGeometry(gl: WebGL2RenderingContext, params: sbWheelParams): [WebGLBuffer, WebGLBuffer, number, WebGLBuffer, number] {
    let [ points, verts ] = wheelPoints(params);

    let wheelPointBuffer = gl.createBuffer();
    if (!wheelPointBuffer) {
        throw new Error("Could not create wheel point buffer");
    }
    gl.bindBuffer(gl.ARRAY_BUFFER, wheelPointBuffer);
    gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(points), gl.STATIC_DRAW);
    gl.bindBuffer(gl.ARRAY_BUFFER, null);

    let vertBuffer = gl.createBuffer();
    if (!vertBuffer) {
        throw new Error("Could not create wheel point buffer");
    }
    gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, vertBuffer);
    gl.bufferData(gl.ELEMENT_ARRAY_BUFFER, new Uint16Array(verts), gl.STATIC_DRAW);
    gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, null);

    let cap = wheelCap(params);
    let wheelCapBuffer = gl.createBuffer();
    if (!wheelCapBuffer) {
        throw new Error("Could not create wheel cap buffer");
    }
    gl.bindBuffer(gl.ARRAY_BUFFER, wheelCapBuffer);
    gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(cap), gl.STATIC_DRAW);
    gl.bindBuffer(gl.ARRAY_BUFFER, null);

    return [wheelPointBuffer, vertBuffer, verts.length, wheelCapBuffer, cap.length / 12];
}