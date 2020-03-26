import { Graphable111 } from "../math";

interface sbWheelParams {
    radius: Graphable111
    length: number 
}

export function wheelGeometry(gl: WebGL2RenderingContext) {
    const params: sbWheelParams = {
        length: .05,
        radius: (t) => [1, 0],
    }
}