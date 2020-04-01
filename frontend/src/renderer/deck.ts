import { Graphable111, Graphable213, rationalEase, bezier, compose, scale, translate, identity, cross, one, absolute } from "../math";

interface sbDeckParams {
    width: number
    height: number
    thickness: number
    perimeterZNormalSize: number
    perimeterStepFactor: number
    tailEllipseExponent: number
    tailLiftX: number
    tailLiftFactor: number
    tailStretchFactor: number
    railLiftX: number
    railLiftFactor: number
    railTailTransition: Graphable111
}

const skateboardTopology: (p: sbDeckParams) => Graphable213 = ({
    tailLiftX,
    tailLiftFactor,
    tailStretchFactor,
    railLiftX,
    railLiftFactor,
    railTailTransition
}) => (x, y) => {
    let fx: Graphable111;
    let fy: Graphable111;
    if (Math.abs(x) > tailLiftX) {
        fx = scale( compose(rationalEase, scale( translate(absolute, -tailLiftX), tailStretchFactor)), tailLiftFactor)
        fy = one;
    } else if (Math.abs(x) < railLiftX) {
        fx = one;
        fy = scale(compose(rationalEase, absolute), railLiftFactor)
    } else {
        fx = compose(railTailTransition, scale(translate(absolute, -railLiftX), 1.0 / (tailLiftX - railLiftX) ))
        fy = scale(compose(rationalEase, absolute), railLiftFactor)
    }
    return cross(fx, fy)(x, y);
}

export const skateboardPerimeter: (p: sbDeckParams) => Graphable111 = ({
    width,
    tailEllipseExponent,
    tailLiftX,
}) => (x) => {
    x = Math.abs(x);
    if (x < tailLiftX) {
        return [1, 0];
    }
    let scale = 1.0 / (width / 2 - tailLiftX);
    x = (x - tailLiftX) * scale;
    let c = tailEllipseExponent;
    let p = 1 - Math.pow(x, c)
    let y = Math.pow(p, 1 / c);
    let d = (1 / c - 1) * Math.pow(p, 1 / c - 1) * (-c + 1) * Math.pow(x, c - 1);
    return [ y, d * scale ]
}

function generateDeck({width, height}: sbDeckParams, deck: Graphable213, perm: Graphable111) {
    let xCells = 50;
    let yCells = 15;
    let xSize = width / (xCells * 2);
    let yOffset = -height / 2;
    let ySize = height / yCells;
    let i = 0;
    let points: number[] = [];
    let indexes: number[] = [];
    for (let iy = 0; iy <= yCells; iy++) {
        for (let ix = 0; ix <= xCells; ix++) {
            let x = ((iy % 2) ? ix : (xCells - ix)) * xSize;
            let y = (iy * ySize + yOffset) * perm(x)[0];
            let [ z, norm ] = deck(x, y);
            points.push(x, y, z, ...norm);
            if (iy) {
                let curr = i;
                let below = curr - 2 * (ix) - 1;
                if (iy % 2)  {
                    indexes.push(below, curr);  
                } else {
                    indexes.push(curr, below);  
                }
            }
            i++;
        }
    }
    return [points, indexes];
}

function generatePerm({height, width, perimeterStepFactor, perimeterZNormalSize, thickness}: sbDeckParams, d: Graphable213, p: Graphable111) {
    let verts: number[] = [];
    let x = -perimeterStepFactor;
    let s = Math.sqrt(1 + perimeterZNormalSize * perimeterZNormalSize);
    let dz = perimeterZNormalSize / s;
    const addPermVertex = (x: number, ySign: number) => {
        let [ y, yp ] = p(x);
        y *= ySign * height / 2;
        let [ z ] = d(x, y)
        let dx, dy;
        if (Math.abs(yp) > 50) {
            dx = -1;
            dy = 0
        } else {
            let n = Math.sqrt(1 + yp * yp);
            dx = -yp / n;
            dy = 1 / n;
        }
        dx /= s;
        dy /= -s;
        dy *= ySign;
        verts.push(
            x,  y,  z - thickness / 2,
            dx, dy, dz,
            x,  y,  z + thickness / 2,
            dx, dy, dz,
        )
        return perimeterStepFactor * (Math.abs(dy) + 0.01);
    }
    while (x <= width / 2.0) {
        x += addPermVertex(x, 1)
    }
    x = width / 2.0;
    while (x >= 0) {
        x -= addPermVertex(x, -1);
    }
    return verts;
}

export function skateboardGeometry(gl: WebGL2RenderingContext, sb: sbDeckParams): [WebGLBuffer, WebGLBuffer, number, WebGLBuffer, number] {
    let d = skateboardTopology(sb);
    let p = skateboardPerimeter(sb);
    let [ deck, deckIndexes ] = generateDeck(sb, d, p);

    let buffer = gl.createBuffer();
    if (!buffer) {
        throw new Error("Could not create skateboard deck vertex buffer");
    }
    gl.bindBuffer(gl.ARRAY_BUFFER, buffer);
    gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(deck), gl.STATIC_DRAW);
    gl.bindBuffer(gl.ARRAY_BUFFER, null);

    let verts = gl.createBuffer();
    if (!verts) {
        throw new Error("Could not create skateboard deck index buffer");
    }
    gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, verts);
    gl.bufferData(gl.ELEMENT_ARRAY_BUFFER, new Uint16Array(deckIndexes), gl.STATIC_DRAW);
    gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, null);

    let perimeter = generatePerm(sb, d, p);
    let perimeterBuffer = gl.createBuffer();
    if (!perimeterBuffer) {
        throw new Error("Could not create skateboard perimeter buffer");
    }
    gl.bindBuffer(gl.ARRAY_BUFFER, perimeterBuffer);
    gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(perimeter), gl.STATIC_DRAW);
    gl.bindBuffer(gl.ARRAY_BUFFER, null);

    return [ buffer, verts, deckIndexes.length, perimeterBuffer, perimeter.length / 6 ];
}