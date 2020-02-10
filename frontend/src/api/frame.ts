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
    succeeded?: boolean
    boundsToUpload?: Bounds
}

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
        let frame = Math.floor(this.elem.currentTime * 30);
        this.elem.currentTime = frame * 30;
        return frame
    }

    public uploadFrame() {
        let canvas = document.createElement("canvas");
        canvas.width = this.elem.videoWidth;
        canvas.height = this.elem.videoHeight;
        let ctx = canvas.getContext("2d");
        if (!ctx) {
            return;
        }
        ctx.drawImage(this.elem, 0, 0);
        let data = canvas.toDataURL();
        let frame = this.currentFrame();
        this.frameIDs[frame] = { isUploading: true }
        fetch(`/frame?video=${this.video.id}&frame=${frame}`, {
            method: "POST",
            body: data,
            headers: session,
        })
        .then(res => res.json())
        .then(({id}) => {
            this.frameIDs[frame].isUploading = false;
            this.frameIDs[frame].succeeded = true;
            this.frameIDs[frame].id = id;
            let waitingBounds = this.frameIDs[frame].boundsToUpload;
            if (waitingBounds) {
                this.uploadBoundsFor(frame, waitingBounds);
            }
        })
        .catch(err => {
            this.frameIDs[frame].isUploading = false;
            this.frameIDs[frame].succeeded = false;
            console.log("Error uploading frame: ", err)
        });
    }

    public giveBounds(bounds: Bounds) {
        this.uploadBoundsFor(this.currentFrame(), bounds);
    }

    private uploadBoundsFor(frame: number, bounds: Bounds) {
        fetch(`/bounds?frame=${frame}`, {
            method: "POST",
            body: JSON.stringify(bounds),
            headers: session,
        })
    }


}