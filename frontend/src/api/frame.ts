import { Quaternion } from "../math";
import { API_URL } from "./url";

export interface Frame {
    id: number
    time: number
    bounds: Bound[]
}

export interface Bound {
    id: number
    x: number
    y: number
    width: number
    height: number
    rotations: Rotation[]
}

export interface Rotation {
    r: number
    i: number
    j: number
    k: number
}

export function toRot(q: Quaternion): Rotation  {
    return {
        r: q[0],
        i: q[1],
        j: q[2],
        k: q[3],
    }
}

export function toQuat({r, i, j, k}: Rotation): Quaternion {
    return [r, i, j, k]
}

export function getFrames(videoId: number): Promise<Frame[]> {
    return fetch(`${API_URL}/api/frames?video=${videoId}`)
    .then(res => res.json())
    .then(({frames}) => frames);
}