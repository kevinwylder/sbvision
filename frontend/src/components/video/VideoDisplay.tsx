import * as React from 'react';

import "../../styles/video.css";

import { Box } from '../../box';
import { Video } from '../../api';

interface VideoDisplayProps {
    video: Video
    topBarSize: number
}

export function VideoDisplay(props: VideoDisplayProps) {

    // then attach the canvas size to respond to the video size and resolution
    let video = React.createRef<HTMLVideoElement>();
    let [ aspectRatio, setAspectRatio ] = React.useState(1);
    React.useEffect(() => {
        if (!video.current) {
            return;
        }
        video.current.onresize = ({target}) => {
            let v = target as HTMLVideoElement;
            setAspectRatio(v.videoWidth / v.videoHeight);
        }
        video.current.load();
    }, [ video.current ]);

    let canvas = React.createRef<HTMLCanvasElement>();

    // create a ref to set the pixel dimensions we are allocated for video
    let container = React.createRef<HTMLDivElement>();
    let [ bounds, setDisplayBounds ] = React.useState({
        top: 0,
        left: 0,
        width: 0,
        height: 0
    })
    React.useLayoutEffect(() => {
        console.log("Render");
        if (!container.current) return;
        let { width, height } = container.current.getBoundingClientRect();
        console.log(height, width / height, aspectRatio);
        let displayLeft = Math.max((width - aspectRatio * (height - 2 * props.topBarSize)), 0) / 2;
        let displayWidth = width - 2 * displayLeft;
        let displayHeight = displayWidth / aspectRatio;
        let displayTop = Math.max((height - displayHeight) / 2, props.topBarSize);
        setDisplayBounds({
            top: displayTop,
            left: displayLeft,
            width: displayWidth,
            height: displayHeight,
        });
    }, [container.current, aspectRatio]);

    return <div 
        ref={container}
        className="video-container">
        <div 
            className="video-control-bar"
            style={{
                top: bounds.top - props.topBarSize,
                left: bounds.left,
                width: bounds.width,
                height: props.topBarSize
            }}>
                Back button!
        </div>
        <video
            ref={video} 
            autoPlay
            style={{
                position: "absolute",
                ...bounds
            }}>
            <source
                src={`/video?type=${props.video.type}&id=${props.video.id}`}
                type={props.video.format}
            />
        </video>
        <canvas 
            ref={canvas}
            style={{
                position: "absolute",
                ...bounds
            }}
        >
        </canvas>
        <div
            className="video-control-bar"
            style={{
                top: bounds.top + bounds.height,
                left: bounds.left,
                width: bounds.width,
                height: props.topBarSize
            }}> Is there a skateboard in this video? Pause and let us know </div>

    </div> 
}

class VideoControls {

    public width = 0;
    public height = 0;
    public aspect = 0;
    public isFirstPlay = true;

    constructor(
        private video: HTMLVideoElement,
        private ctx: CanvasRenderingContext2D,
    ) {
    }

    public draw() { 
        this.ctx.clearRect(0, 0, this.width, this.height);
        if (this.video.paused) {
            if (this.isFirstPlay) {
                this.drawPlayArrow();
            }
        } else {
            this.isFirstPlay = false;
        }
    }

    private drawPlayArrow() {
        let size = this.width * .1;
        let top = (this.height - size) / 2;
        let bottom = (this.height + size) / 2;
        let left = this.width / 2 - size;
        let right = this.width / 2 + size;
        this.ctx.beginPath();
        this.ctx.moveTo(left, top);
        this.ctx.lineTo(right, this.height / 2);
        this.ctx.lineTo(left, bottom);
        this.ctx.closePath();
        this.ctx.fillStyle = "white";
        this.ctx.strokeStyle = "black";
        this.ctx.lineWidth = 6;
        this.ctx.lineCap = "round";
        this.ctx.fill();
        this.ctx.stroke();
    } 

}