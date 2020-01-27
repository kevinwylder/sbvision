import * as React from 'react';
import { Box } from '../model';

interface VideoComponentProps {
    url: string
    width: number
    callback: (data: ImageData) => void
}

export function VideoController(props: VideoComponentProps) {

    let [ video, setVideoElement ] = React.useState<HTMLVideoElement|null>(null)
    let [ box, setBox ] = React.useState<Box>([0, 0, 0, 0]);
    let [ isPlaying, setIsPlaying ] = React.useState(true);
    let [ isDragging, setIsDragging ] = React.useState(false);

    let canvas = React.createRef<HTMLCanvasElement>();
    let videoRef = React.useCallback((v: HTMLVideoElement) => {
        if (v) {
            setVideoElement(v);
        }
    }, []);

    const renderCanvas = () => {
        if (isPlaying) {
            return;
        }
        if (!video) {
            return;
        }
        if (!canvas.current) {
            return;
        }
        let { width, height } = canvas.current;
        let ctx = canvas.current.getContext("2d");
        if (!ctx) {
            return
        }

        // draw the box on the video
        ctx.clearRect(0, 0, width, height);
        ctx.drawImage(video, 0, 0);
        ctx.fillStyle = "rgba(0, 0, 0, .5)";
        let [ x0, y0, x1, y1 ] = box;
        ctx.fillRect(0, 0, Math.min(x0, x1), Math.max(y0, y1));
        ctx.fillRect(Math.min(x0, x1), 0, width, Math.min(y0, y1));
        ctx.fillRect(Math.max(x0, x1), Math.min(y0, y1), width, height);
        ctx.fillRect(0, Math.max(y0, y1), Math.max(x0, x1), height);
    }
    React.useEffect(renderCanvas, [ isPlaying, video, box ]);

    const setBoxPos = (event: React.MouseEvent, starting: boolean) => {
        if (!canvas.current) {
            return;
        }
        let { clientX, clientY } = event;
        let { top, left, width, height } = canvas.current.getBoundingClientRect();
        let x = canvas.current.width * (clientX - left) / width;
        let y = canvas.current.height * (clientY - top) / height;
        if (starting) {
            setBox([x, y, x, y]);
        } else {
            let [sx, sy] = box;
            setBox([sx, sy, x, y]);
        }
    }

    const playPause = () => {
        if (isPlaying) {
            setIsPlaying(false);
            video?.pause();
        } else {
            setIsPlaying(true);
            video?.play();
        }
    }

    const forward = () => {
        if (!video) {
            return;
        }
        video.currentTime += 0.1;
    }

    const back = () => {
        if (!video) {
            return;
        }
        video.currentTime -= 0.1;
    }

    const clipSection = () => {
        if (!canvas.current) {
            return
        }
        let ctx = canvas.current.getContext("2d");
        if (!ctx) {
            return;
        }
        props.callback(ctx.getImageData(box[0], box[1], box[2] - box[0], box[3] - box[1]));
        forward();
    }

    return (
        <div className="video-container">
            <video 
                onTimeUpdate={renderCanvas}
                ref={videoRef} 
                controls
                style={{
                    display: (isPlaying) ? "block" : "none",
                    width: props.width,
                }}> 
                <source src={props.url} type="video/mp4" />
            </video>

            { video != null && <><canvas 
                    ref={canvas} 
                    width={video.videoWidth} 
                    height={video.videoHeight} 
                    style={{ 
                        display: (isPlaying) ? "none" : "block",
                        width: props.width,
                        height: props.width * video.videoHeight / (video.videoWidth + 1),
                        border: "1px solid black"
                    }} 
                    onMouseDown={(event: React.MouseEvent) => {
                        setBoxPos(event, true);
                        setIsDragging(true);
                    }} 
                    onMouseMove={(event: React.MouseEvent) => {
                        if (isDragging) {
                            setBoxPos(event, false);
                        }
                    }}
                    onMouseUp={(event: React.MouseEvent) => {
                        setBoxPos(event, false);
                        setIsDragging(false);
                    }}>
                </canvas>
                <br/>
                <div className="video-controls" style={{width: props.width}}> 

                    <div className="video-controls-button"
                        onClick={back}> 
                        Back a little
                    </div>

                    <div className="video-controls-button" 
                        onClick={playPause}> 
                        {isPlaying ? "Pause" : "Play"} 
                    </div>

                    <div className="video-controls-button"
                        onClick={forward}> 
                        Forward a little
                    </div>

                    <div className="video-controls-button"
                        onClick={clipSection}> 
                        Clip 
                    </div>

                </div>
            </>}
        </div>
    )
}