import { Frame } from "../api";
import { renderSkateboard } from "./skateboard";

const ANNOTATION_FRAME_RATE = 1000 / 24;

/**
 * FrameList keeps a list of frames sorted by time 
 */

export class FrameList {

    private index = 0;
    private frames: Frame[];

    constructor(frames: Frame[]) {
        if (!frames) {
            frames = [];
        }
        this.frames = frames.sort((a, b) => a.time - b.time);
        console.log(this.frames);
    }

    // getFrame looks up the object stored for the time
    public getFrame(time: number): Frame|undefined {
        if (this.frames.length == 0) {
            return undefined;
        }
        if (this.frames.length == 1) {
            return Math.abs(this.frames[0].time - time) < ANNOTATION_FRAME_RATE 
                ? this.frames[0] 
                : undefined;
        }

        while (this.index < this.frames.length - 2 && this.frames[this.index].time <= time) {
            this.index++;
        }
        while (this.index > 0 && this.frames[this.index].time > time) {
            this.index--;
        }
        let before = this.frames[this.index];
        let beforeOff = Math.abs(before.time - time);
        let after = this.frames[this.index + 1];
        let afterOff = Math.abs(after.time - time);
        if (Math.min(beforeOff, afterOff) > ANNOTATION_FRAME_RATE) {
            return undefined;
        } else if (afterOff > beforeOff) {
            return before;
        } else {
            return after;
        }
    }

    // addFrame adds the frame to the list but keeps the sorted order
    public addFrame(frame: Frame) {
        if (this.frames.length == 0 || this.frames[0].time > frame.time) {
            this.frames.unshift(frame);
            return
        }
        if (this.frames.length == 1) {
            this.frames.push(frame);
            return
        }

        let idx = this.index;
        while (idx < this.frames.length - 1 && this.frames[idx].time <= frame.time) {
            idx++;
        }
        while (idx >= 0 && this.frames[idx].time > frame.time) {
            idx--;
        }
        if (this.frames[idx].time == frame.time) {
            let other = this.frames[idx]
            if (other.bounds != null && frame.bounds != null) {
                other.bounds.splice(other.bounds.length - 1, 0, ...frame.bounds);
            }
        } else {
            this.frames.splice(idx + 1, 0, frame);
        }
    }

    public drawFrame(ctx: CanvasRenderingContext2D, time: number) {
        let annotation = this.getFrame(time);
        if (annotation?.bounds) {
            annotation.bounds.forEach((bound) => {
                if (bound.rotations?.length) {
                    let { r, i, j, k } = bound.rotations[0];
                    renderSkateboard(ctx, [r, i, j, k], [bound.x, bound.y, bound.x + bound.width, bound.y + bound.height]);
                } else {
                    ctx.lineWidth = 4;
                    ctx.strokeStyle = "#FF0000";
                    ctx.strokeRect(bound.x, bound.y, bound.width, bound.height);
                }
            });
        }
    }

}