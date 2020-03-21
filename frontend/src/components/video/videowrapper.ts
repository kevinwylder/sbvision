import { getFrames, getVideoById, addBounds, uploadFrame } from "../../api";
import { Box } from "./boxmanager";
import { FrameList } from "./framesearch";
import { PlayButton } from "./playbutton";
import { renderSkateboard } from "../../renderer";

export class VideoWrapper {

    private video: HTMLVideoElement;

    private ctx: CanvasRenderingContext2D|undefined;
    private frames: FrameList|undefined;
    private box: Box|undefined;
    private playbutton: PlayButton|undefined;

    constructor (
        private videoID: number,
        onload: (width: number, height: number) => void,
    ) {

        this.video = document.createElement("video");
        this.video.style.visibility = 'hidden';
        this.video.style.width = "1";
        this.video.style.height = "1";
        this.video.controls = true;
        this.video.setAttribute("playsinline", "");
        this.video.loop = true;
        document.body.appendChild(this.video);

        this.video.onplay = () => {
            if (this.playbutton) {
                this.playbutton.cancel();
                this.playbutton = undefined;
            }
            this.render();
        }

        this.video.onpause = () => {
            this.render();
        }

        this.video.onresize = () => {
            this.box = new Box(6, this.video.videoWidth, this.video.videoHeight);
            onload(this.video.videoWidth, this.video.videoHeight);
            this.video.play();
        }

        this.video.ontimeupdate = () => {
            if (this.video.paused) {
                uploadFrame(this.videoID, this.video)
                .then(frame => {
                    this.frames?.addFrame(frame);
                })
                .catch(err => console.log(err));
                this.render();
            }
        }

        getFrames(videoID)
        .then(frames => {
            this.frames = new FrameList(frames);
        })

        getVideoById(videoID)
        .then(videoInfo => {
            var source = document.createElement('source');
            source.src = `/app/video/stream?id=${videoInfo.id}`;
            source.type = videoInfo.format;
            this.video.appendChild(source);
            this.video.load();
            this.render();
        })
    }

    public mobileClick() {
        this.video.play();
    }

    public grab(pos: [number, number]) {
        this.box?.grab(pos);
        this.render();
    }
    
    public drag(pos: [number, number]) {
        if (this.box?.drag(pos)) {
            this.render();
        }
    }

    public up() {
        if (!this.box) {
            return;
        }

        let [ wasTap, wasInsideBox, wasEmptyBox ] = this.box.release();

        if (wasTap && wasEmptyBox) {
            if (this.video.paused) {
                this.video.play();
            } else {
                this.video.pause();
            }
            return;
        }

        if (!this.video.paused) {
            this.box.reset();
            return;
        }

        if (!wasTap) {
            if (this.ctx) {
                this.box.drawHelpDelayed(this.ctx);
            }
            return;
        } 

        if (wasInsideBox) {
            let frame = this.frames?.getFrame(this.video.currentTime * 1000);
            if (frame && this.box) {
                addBounds(frame.id, this.box.bounds())
                .then(() => { 
                    this.forward();
                })
                .catch(err => {
                    console.log(err);
                    this.video.play();
                });
                return;
            }
        }
        console.log("Click outside box");
        this.box.reset();
        this.video.play();
    }

    public forward() {
        this.video.currentTime += 0.02; 
        this.render();
    }

    public back() {
        this.video.currentTime -= 0.02; 
        this.render();
    }
    
    public setContext(ctx: CanvasRenderingContext2D) {
        this.ctx = ctx;
        this.playbutton = new PlayButton(this.videoID, ctx, this.video.videoWidth, this.video.videoHeight);
    }

    public destroy() {
        this.video.pause();
        this.ctx = undefined;
        this.video.remove();
    }

    private render() {
        if (!this.ctx) return;

        this.ctx.drawImage(this.video, 0, 0);
        this.drawAnnotations();
        if (this.video.paused) {
            this.box?.draw(this.ctx);
        } else {
            window.requestAnimationFrame(() => this.render());
        }
    }

    private drawAnnotations() {
        if (!this.ctx || !this.frames) {
            return;
        }
        let annotation = this.frames.getFrame(this.video.currentTime * 1000);
        for (let i = 0; annotation?.bounds && i < annotation.bounds.length; i++) {
            let bound = annotation.bounds[i];
            if (bound.rotations?.length) {
                let { r, i, j, k } = bound.rotations[0];
                renderSkateboard(this.ctx, [r, i, j, k], [bound.x, bound.y, bound.x + bound.width, bound.y + bound.height]);
            } else {
                this.ctx.lineWidth = 4;
                this.ctx.strokeStyle = "#FF0000";
                this.ctx.strokeRect(bound.x, bound.y, bound.width, bound.height);
            }
        }
    }

}