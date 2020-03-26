import { getFrames, Frame, addRotation } from "../../api";
import { renderSkateboard } from "../../renderer";
import { Quaternion, rotateSkateboard, tiltSkateboard } from '../../math';

export interface ImageWrapperFrameStatus {
    frame: Frame|undefined
    selectedBound: number|undefined
}

/**
 * ImageWrapper allows review of individual frames for a specific video
 */
export class ImageWrapper {

    private frames: Frame[] = [];
    private idx: number = 0;
    private frame: Frame|undefined;
    private selectedBound: number = 0;

    private image: HTMLImageElement;
    private ctx: CanvasRenderingContext2D;
    private rotation: Quaternion;

    constructor(
        videoID: number,
        private canvas: HTMLCanvasElement,
        private onFrameUpdate: (status: ImageWrapperFrameStatus) => void
    ) {
        let ctx = canvas.getContext("2d");
        if (!ctx) {
            throw new Error("Invalid canvas");
        }
        this.ctx = ctx;
        this.image = document.createElement("img");
        this.image.style.visibility = 'hidden';
        this.image.onload = () => this.draw();
        this.rotation = [ 1, 0, 0, 0 ];
        getFrames(videoID)
        .then(frames => {
            this.frames = frames;
            this.updateFrame();
        })
        .catch(err => console.log(err));
    }

    public move(e: {movementX: number, movementY: number}) {
        this.rotation = rotateSkateboard(e.movementX, e.movementY, this.rotation);
        this.draw();
    }

    public scroll(dx: number) {
        this.rotation = tiltSkateboard(dx, this.rotation);
        this.draw();
    }

    public click(e: {clientX: number, clientY: number}) { 
        if (!this.frame?.bounds) {
            return;
        }
        if (!this.selectedBound) {
            // compute the position on the image from the css position
            let { clientX, clientY } = e;
            let { top, left, width } = this.canvas.getBoundingClientRect();
            let ratio = this.image.width / width;
            let px = (clientX - left) * ratio;
            let py = (clientY - top) * ratio;
            this.selectNearestBound(px, py);
            return;
        }

        let rotation = { r: this.rotation[0], i: this.rotation[1], j: this.rotation[2], k: this.rotation[3] };

        // push this rotation to the current bound
        for (let bound of this.frame.bounds) {
            if (bound.id == this.selectedBound) {
                bound.rotations.push(rotation);
            }
        }

        addRotation(this.selectedBound, rotation)
        .catch(err => console.log(err));

        this.selectNextBound();
    }

    public remove() {
        // TODO: how do we trust you're not trolling here?
    }

    public next() {
        this.idx ++;
        this.updateFrame();
    }

    public prev() {
        this.idx --;
        this.updateFrame();
    }

    public setBound(boundID: number) {
        this.selectedBound = boundID;
        this.draw();
        this.onFrameUpdate({
            frame: this.frame,
            selectedBound: this.selectedBound,
        })
    }

    public destroy() {
        this.image.remove();
    }

    private updateFrame() {
        if (this.idx < 0 || this.idx >= this.frames.length) {
            this.frame = undefined;
            this.idx = 0;
            this.draw();
        } else {
            this.frame = this.frames[this.idx];
            this.image.src = `/api/image?frame=${this.frame.id}`;
        }
        this.selectedBound = 0;
        // pick the first rotationless
        for (let bound of this.frame?.bounds || []) {
            if (bound.rotations.length == 0) {
                this.selectedBound = bound.id;
                break;
            }
        }
        this.onFrameUpdate({
            frame: this.frame,
            selectedBound: this.selectedBound,
        })
    }

    private draw() {
        if (this.image.height == 0 || this.image.width == 0) {
            this.canvas.width = 150;
            this.canvas.height = 100;
            this.ctx.fillText("No Video Frames Here!", 10, 10);
        }
        this.canvas.width = this.image.width;
        this.canvas.height = this.image.height;
        let { width } = this.canvas.getBoundingClientRect();
        let height = width * this.image.height / this.image.width; 
        this.canvas.style.height = height + "px";
        this.ctx.drawImage(this.image, 0, 0);

        let skateboardBox: [number, number, number, number] | undefined;
        this.frame?.bounds.forEach(({id, x, y, width, height}) => {
            if (this.selectedBound == id) {
                let long = Math.min(this.image.width, this.image.height) * .25;
                let offsetX = (2 * x + width < this.image.width) ? width : -long;
                let offsetY = (2 * y + height < height) ? 0 : height - long;
                skateboardBox = [ x + offsetX, y + offsetY, x + offsetX + long, y + offsetY + long];
                this.ctx.strokeStyle = "#00FF00";
            } else {
                this.ctx.strokeStyle = "#FF0000";
            }
            this.ctx.lineWidth = 4;
            this.ctx.strokeRect(x, y, width, height);
        });

        if (skateboardBox) {
            renderSkateboard(this.ctx, this.rotation, skateboardBox);
        }
    }


    private selectNearestBound(px: number, py: number) {
        if (!this.frame) {
            return;
        }
        let minDistance = Infinity;
        let nearestBound = 0;
        for (let {id, x, y, width, height, rotations} of this.frame.bounds) {
            if (Math.abs(x - px) < minDistance) {
                nearestBound = id;
                minDistance = Math.abs(x - px);
            }
            if (Math.abs(x + width - px) < minDistance) {
                nearestBound = id;
                minDistance = Math.abs(x + width - px);
            }
            if (Math.abs(y - py) < minDistance) {
                nearestBound = id;
                minDistance = Math.abs(y - py);
            }
            if (Math.abs(y + height - py) < minDistance) {
                nearestBound = id;
                minDistance = Math.abs(y + height - py);
            }
        }
        this.setBound(nearestBound);
    }

    private selectNextBound() {
        if (!this.frame) {
            return;
        }
        for (let bound of this.frame.bounds) {
            if (bound.rotations.length == 0) {
                this.setBound(bound.id);
                return
            }
        }
        this.next();
    }
}