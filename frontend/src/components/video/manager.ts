import { Bounds } from './box';
import { CollectionStatistics } from '../CollectionStats';
import { Frame, getFrames, uploadFrame, uploadBounds } from '../../api';

/**
 * The frame manager is responsible for matching requests to upload a frame with 
 * bounds applied to the frames.
 * 
 * The server is responsible for keeping a unique mapping of frames
 */
export class VideoManager {

    private frames: Frame[] = [];

    constructor(
        private videoId: number,
        private elem: HTMLVideoElement,
        private ctx: CanvasRenderingContext2D,
        private updateStats: (stats: CollectionStatistics) => void,
    ) { 
        // bind time updates to our scheduler with a local function
        this.elem.addEventListener("timeupdate", this.onTimeUpdate.bind(this));
        // get the already collected frames of the video
        getFrames(videoId)
        .then(frames => {
            this.frames = frames.sort((a, b) => a.time - b.time);
            this.countStats();
        })
        .catch(err => console.log(err));
    }

    private countStats() {
        let frameCount = 0;
        let boundCount = 0; 
        let rotationCount = 0;
        if (frames) {
            this.frames.forEach(({bounds}) => {
                frameCount++;
                if (bounds) {
                    bounds.forEach(({rotations}) => {
                        boundCount++;
                        if (rotations) {
                            rotationCount += rotations.length;
                        }
                    })
                }
            })
        }
        this.updateStats({ bounds: boundCount, frames: frameCount, rotations: rotationCount });
    }


    private onTimeUpdate() {
        if (!this.elem.paused) {
            this.scheduleBounds();
            return;
        } else {
            this.uploadFrame();
        }
    }

    private lastScheduledTime: number = 0;
    private boundsSchedule: number[] = [];
    private scheduleBounds() {
        let previousDelta = this.elem.currentTime - this.lastScheduledTime;
        this.lastScheduledTime = this.elem.currentTime;
        if (previousDelta < 0 || previousDelta > 1000) {
            return;
        }
        // schedule the next previousDelta of boxes
        let [ nextFrame, _ ] = this.findFrame()
        this.boundsSchedule.forEach(timeout => clearTimeout(timeout));
        this.ctx.clearRect(0, 0, this.elem.videoWidth, this.elem.videoHeight);
        this.boundsSchedule = [];
        while (++nextFrame < this.frames.length && this.frames[nextFrame].time < (this.elem.currentTime + previousDelta) * 1000 ) {
            let frame = this.frames[nextFrame];
            if (!frame.bounds) {
                continue;
            }
            this.boundsSchedule.push(setTimeout(() => {
                frame.bounds.forEach(({x, y, width, height}) => {
                    this.ctx.strokeStyle = "#FF0000";
                    this.ctx.lineWidth = 3;
                    this.ctx.strokeRect(x, y, width, height);
                })
            }, frame.time - this.elem.currentTime * 1000));

        }
    }

    private uploadFrame() {
        uploadFrame(this.videoId, this.elem)
        .then(frame => {
            this.frames.push(frame);
            this.frames.sort((a, b) => a.time - b.time);
            this.countStats();
        });
    }

    // binary search for the last frame that has a timestamp less than or equal to the current time
    private findFrame(): [ number, boolean ] {
        let frameTime = Math.ceil(this.elem.currentTime * 1000);
        let low = 0;
        let high = this.frames.length;
        while (high > low) {
            let middle = Math.floor((low + high) / 2);
            let frame = this.frames[middle];
            if (frameTime == frame.time) {
                return [ middle, true ];
            } else if (frameTime < frame.time) {
                high = middle;
            } else {
                low = middle + 1;
            }
        }
        return [ low, false ];
    }

    public giveBounds(bounds: Bounds) {
        let x = Math.min(Math.max(Math.round(bounds.left), 0), this.elem.videoWidth);
        let y = Math.min(Math.max(Math.round(bounds.top), 0), this.elem.videoHeight);
        let width = Math.min(Math.max(Math.round(bounds.width), 0), this.elem.videoWidth - x);
        let height = Math.min(Math.max(Math.round(bounds.height), 0), this.elem.videoHeight - y);
        let [ frameNumber, wasFound ] = this.findFrame();
        if (wasFound) {
            const frame = this.frames[frameNumber];
            uploadBounds(frame.id, { x, y, width, height })
            .then(({id}) => {
                let bound = {
                    id, 
                    frameId: frame.id, 
                    height, 
                    width, 
                    x, 
                    y, 
                    rotations: []
                };
                if (!frame.bounds) {
                    frame.bounds = [ bound ]
                } else {
                    frame.bounds.push( bound );
                }
                this.countStats();
            })
            .catch(err => console.log(err));
            this.elem.currentTime += 0.042; // slightly more than 24fps;
        }
    }

    public continue() {
        this.elem.play();
    }

    public unregister() {
        this.elem.removeEventListener("timeupdate", this.onTimeUpdate);
    }

}