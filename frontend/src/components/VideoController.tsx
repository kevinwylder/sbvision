import * as React from 'react';

interface VideoComponentProps {
    url: string
    width: number
    height: number
}

type HighlightBox = [ number, number, number, number ];

export function VideoController(props: VideoComponentProps) {

    let [ box, setBox ] = React.useState<HighlightBox>([0, 0, props.width, props.height])

    let canvas = React.createRef<HTMLCanvasElement>();
    let youtube = React.createRef<HTMLIFrameElement>()

    const renderCanvas = () => {
        if (!youtube.current) {
            console.log("no youtube")
            return;
        }
        let videos = youtube.current.contentWindow?.document.getElementsByTagName("video")
        if (!videos || videos.length == 0) {
            console.log("videos", videos);
            return
        }
        let [ video ] = videos;
        video.pause();
        if (!canvas.current) {
            return;
        }
        let ctx = canvas.current.getContext("2d");
        if (!ctx) {
            return
        }
        ctx.drawImage(video, 0, 0);
        ctx.fillStyle = "rgba(0, 0, 0, .1)";
        let [ x0, y0, x1, y1 ] = box;
        ctx.fillRect(0, 0, Math.min(x0, x1), props.height);
        ctx.fillRect(0, 0, props.width, Math.min(y0, y1));
        ctx.fillRect(Math.max(x0, x1), 0, props.width, props.height);
        ctx.fillRect(0, Math.max(y0, y1), props.width, props.height);
    }

    const setBoxStartingCorner = (event: React.MouseEvent) => {
        if (!canvas.current) {
            return;
        }
        let { clientX, clientY } = event;
        let { top, left } = canvas.current.getBoundingClientRect();
        setBox([clientX - left, clientY - top, clientX - left, clientY - top]);
    }

    const setBoxEndingCorner = (event: React.MouseEvent) => {
        if (!canvas.current) {
            return;
        }
        let { clientX, clientY } = event;
        let { top, left } = canvas.current.getBoundingClientRect();
        let [ x, y ] = box;
        setBox([x, y, clientX - left, clientY - top]);
    }

    return (
        <div>
            <iframe 
                ref={youtube} 
                width={props.width} 
                height={props.height} 
                src={props.url} > 
            </iframe>

            <canvas 
                ref={canvas} 
                width={props.width} 
                height={props.height} 
                style={{ 
                    width: props.width,
                    height: props.height,
                    border: "1px solid black"
                }} 
                onMouseDown={setBoxStartingCorner} 
                onMouseMove={setBoxEndingCorner}
                onMouseUp={setBoxEndingCorner}>
            </canvas>
            <br/>
            <button onClick={() => renderCanvas()}>Clip</button>
        </div>
    )
}