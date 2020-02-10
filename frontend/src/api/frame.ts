import { Video } from './videos';
import { session } from './session';

export interface Bounds {
    left: number
    top: number
    width: number
    height: number
}

interface frameStatus {
    id?: number
    isUploading: boolean
    failed?: boolean
    boundsToUpload?: Bounds
}

const FRAME_RATE = 24;

/**
 * The frame manager is responsible for matching requests to upload a frame with 
 * bounds applied to the frames.
 * 
 * The server is responsible for keeping a unique mapping of frames
 */
export class FrameManager {

    private frameIDs: { [time: number]: frameStatus } = {};

    constructor(
        private video: Video,
        private elem: HTMLVideoElement
    ) { }

    private currentFrame() {
        let frame = Math.ceil(this.elem.currentTime * FRAME_RATE);
        this.elem.currentTime = frame / FRAME_RATE;
        return frame
    }

    public uploadFrame() {
        let frame = this.currentFrame();
        if (this.frameIDs[frame]) {
            return;
        }
        let canvas = document.createElement("canvas");
        canvas.width = this.elem.videoWidth;
        canvas.height = this.elem.videoHeight;
        let ctx = canvas.getContext("2d");
        if (!ctx) {
            return;
        }
        ctx.drawImage(this.elem, 0, 0);
        // check if the frame is all black
        let okay = false
        for (let i = 0; !okay && i < 10; i++) {
            let pixel = ctx.getImageData(Math.floor(Math.random() * this.elem.videoWidth), Math.floor(Math.random() * this.elem.videoHeight), 1, 1);
            if (pixel.data[0] + pixel.data[1] + pixel.data[2] > 20) {
                okay = true;
            }
        }
        if (!okay) {
            return;
        }
        let data = canvas.toDataURL();
        this.frameIDs[frame] = { isUploading: true }
        fetch(`/frame?video=${this.video.id}&frame=${frame}`, {
            method: "POST",
            body: data,
            headers: session,
        })
        .then(res => res.json())
        .then(({id}) => {
            this.frameIDs[frame].isUploading = false;
            this.frameIDs[frame].id = id;
            let waitingBounds = this.frameIDs[frame].boundsToUpload;
            if (waitingBounds) {
                this.uploadBoundsFor(frame, waitingBounds);
            }
        })
        .catch(err => {
            this.frameIDs[frame].isUploading = false;
            this.frameIDs[frame].failed = true;
            console.log("Error uploading frame: ", err)
        });
    }

    public giveBounds(bounds: Bounds) {
        this.uploadBoundsFor(this.currentFrame(), bounds);
    }

    private uploadBoundsFor(frame: number, bounds: Bounds) {
        let status = this.frameIDs[frame];
        if (!status || status.failed) {
            console.log("Missing frame", frame, status);
            return;
        }
        if (status.isUploading) {
            this.frameIDs[frame].boundsToUpload = bounds;
            return;
        }
        let x = Math.min(Math.max(Math.round(bounds.left), 0), this.elem.videoWidth);
        let y = Math.min(Math.max(Math.round(bounds.top), 0), this.elem.videoHeight);
        let width = Math.min(Math.max(Math.round(bounds.width), 0), this.elem.videoWidth - x);
        let height = Math.min(Math.max(Math.round(bounds.height), 0), this.elem.videoHeight - y);
        fetch(`/bounds?frame=${status.id}`, {
            method: "POST",
            body: JSON.stringify({x, y, width, height}),
            headers: session,
        })
    }


}