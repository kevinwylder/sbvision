
export interface Frame {
    id: number
    time: number
    bounds: Bound[]
}

export interface Bound {
    id: number
    frameId: number
    x: number
    y: number
    width: number
    height: number
    rotations: Rotation[]
}

export interface Rotation {
    boundId: number
    id: number
    r: number
    i: number
    j: number
    k: number
}

export function getFrames(videoId: number): Promise<Frame[]> {
    return fetch(`/api/frames?video=${videoId}`)
    .then(res => res.json())
    .then(({frames}) => frames);
}

export function getRotations(): Promise<Frame[]> {
    return fetch(`/api/frames?rotationless=1`)
    .then(res => res.json())
    .then(({frames}) => frames);
}