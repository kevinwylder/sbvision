import { session } from "./session";

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

export class RotationlessFramePager {

    private offset = 0;
    private frames: Frame[] = [];

    public nextFrame(): Promise<Frame|undefined> {
        if (this.frames.length) {
            return Promise.resolve(this.frames.shift());
        }
        return fetch(`/api/frames?rotationless=1&offset=${this.offset}`)
        .then(res => res.json())
        .then(({frames, offset}) => {
            this.offset = offset;
            if (frames) {
                this.frames = frames;
                return frames[0];
            } else {
                this.frames = [];
                return undefined;
            }
        })
    }

    public addRotation(boundId: number, quaternion: {r: number, i: number, j: number, k: number}) {
        return fetch(`/app/contribute/rotation?bound=${boundId}`,{
            method: "POST",
            headers: session,
            body: JSON.stringify(quaternion),
        })
        .then(res => res.json())
        .then(_ => { this.offset -= 1});
    }

    public skipRotation() {
        // TODO: 
    }
}