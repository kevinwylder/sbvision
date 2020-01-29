import * as React from 'react';
import { Label, encodeImage } from '../model';
import { Quaternion } from '../math';
import { renderSkateboard, rotateSkateboard, tiltSkateboard } from '../skateboard';
import { drawBoundingBox, describeBox, BOX_BORDER, chooseDraggingIndexes, dragBoxBounds, isInsideBox } from '../bounds';

type Box = [ number, number, number, number ];

interface VideoComponentProps {
    url: string
    width: number
    callback: (label: Label) => void
}

/**
 * VideoController is an HTML5 video player that switches to a canvas when paused
 * there are 3 states: 
 *  1. video is playing or about to be played for the first time. Use native video controls
 *  2. Video is paused and the user is expected to create a clip of the area
 *  3. Clip has been selected and the user is expected to match the skateboard rotation
 */

export function VideoController(props: VideoComponentProps) {

    let [ isPlaying, setIsPlaying ] = React.useState(true);
    let [ lastMousePosition, setLastMousePosition ] = React.useState<[number, number]>([0, 0]);
    let [ draggingIndexes, setDraggingIndexes ] = React.useState<number[]>([]);
    let [ box, setBox ] = React.useState<Box>([0, 0, 0, 0]);
    let [ quaternion, setQuaternion ] = React.useState<Quaternion>([Math.SQRT1_2, Math.SQRT1_2, 0, 0]);

    let canvas = React.createRef<HTMLCanvasElement>();
    let [ video, setVideoElement ] = React.useState<HTMLVideoElement|null>(null)
    let videoRef = React.useCallback((v: HTMLVideoElement) => {
        if (v) {
            setVideoElement(v);
        }
    }, []);

    const getMousePosition: (e: React.MouseEvent) => [number, number]| null = (e) =>  {
        if (!canvas.current || !video) {
            return null;
        }
        let { clientX, clientY } = e;
        let { top, left } = canvas.current.getBoundingClientRect();
        let x = (clientX - left) * video.videoWidth / props.width;
        let y = (clientY - top) * video.videoWidth / props.width;
        return [ x, y ];
    }

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

        // draw the video
        ctx.clearRect(0, 0, width, height);
        ctx.drawImage(video, 0, 0);

        // draw the box on the video
        drawBoundingBox(ctx, width, height, box);

        if (draggingIndexes.length > 0) {
            return;
        }
        
        let { top, left, right, bottom } = describeBox(box);
        // draw the skateboard next to the box
        let boxWidth = Math.max(right - left, bottom - top);
        let yDelta = (bottom - top - boxWidth) / 2;
        if ( width - right > left ) {
            // skateboard is to the right of the clip
            renderSkateboard(ctx, quaternion, [right + BOX_BORDER, top + yDelta, right + boxWidth + BOX_BORDER, bottom - yDelta]);
        } else {
            // skateboard is to the left of the clip
            renderSkateboard(ctx, quaternion, [left - boxWidth - BOX_BORDER, top + yDelta, left - BOX_BORDER, bottom - yDelta]);
        }

    }
    React.useEffect(renderCanvas, [ isPlaying, video, box, quaternion, video?.currentTime ]);

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
        video.currentTime += 0.05;
    }

    const back = () => {
        if (!video) {
            return;
        }
        video.currentTime -= 0.05;
    }

    const confirmClip = (isSkateboard: boolean) => {
        if (!canvas.current) {
            return;
        }
        let ctx = canvas.current.getContext("2d");
        if (!ctx) {
            return;
        }
        let { top, left, right, bottom } = describeBox(box);
        let tmp = document.createElement("canvas");
        tmp.width = right - left;
        tmp.height = bottom - top;
        let image = ctx.getImageData(...box);
        ctx = tmp.getContext("2d");
        if (!ctx) {
            return;
        }
        ctx.putImageData(image, 0, 0);
        props.callback({
            input: {
                width: Math.round(right - left),
                height: Math.round(bottom - top),
                data: tmp.toDataURL(),
            },
            output: {
                isSkateboard: (isSkateboard) ? 1 : 0,
                rotation: quaternion
            }
        })
        forward();
    }

    const setUnitQuaternion = () => {
        setQuaternion([1, 0, 0, 0]);
    }

    window.onkeydown = (e: KeyboardEvent) => {
        switch (e.key) {
        case "r":
        case "R":
            setUnitQuaternion();
            break;
        case "Enter":
            confirmClip(!e.shiftKey);
            break;
        case "ArrowRight":
            forward();
            break;
        case "ArrowLeft":
            back();
            break;
        case " ":
            playPause();
            break;
        }
    }

    return (
        <div className="video-container" >
            <video 
                onTimeUpdate={renderCanvas}
                ref={videoRef} 
                autoPlay
                controls
                style={{
                    display: (isPlaying) ? "block" : "none",
                    width: props.width,
                }}> 
                <source src={props.url} type="video/mp4" />
            </video>

            { video != null && <canvas 
                ref={canvas} 
                width={video.videoWidth} 
                height={video.videoHeight} 
                style={{ 
                    display: (isPlaying) ? "none" : "block",
                    width: props.width,
                    height: props.width * video.videoHeight / (video.videoWidth + 1),
                    border: "1px solid black"
                }} 
                onMouseDown={(e) => {
                    let pos = getMousePosition(e);
                    if (!pos) {
                        return;
                    }
                    let [ x, y ] = pos;
                    setLastMousePosition(pos);
                    setDraggingIndexes(chooseDraggingIndexes(box, x, y));
                    if (!isInsideBox(box, x, y)) {
                        setBox([x, y, x, y]);
                    }
                }}
                onMouseMove={(e) => {
                    let pos = getMousePosition(e);
                    if (!pos) {
                        return;
                    }
                    let dx = pos[0] - lastMousePosition[0];
                    let dy = pos[1] - lastMousePosition[1];
                    e.preventDefault();
                    setLastMousePosition(pos);
                    if (draggingIndexes.length) {
                        setBox(dragBoxBounds(box, dx, dy, draggingIndexes));
                    } else {
                        setQuaternion(rotateSkateboard(dx, dy, quaternion));
                    }
                }}
                onWheel={(e) => {
                    setQuaternion(tiltSkateboard(e.deltaY, quaternion))
                }}
                onMouseLeave={() => setDraggingIndexes([]) }
                onMouseOut={() => setDraggingIndexes([]) }
                onMouseUp={() => setDraggingIndexes([]) }>
            </canvas>
            }
            <br/>
            <div className="video-controls" style={{width: props.width}}> 

                <div className="video-controls-button" 
                    onClick={playPause}> 
                    {isPlaying ? "Pause" : "Play"} (Spacebar)
                </div>

                <div className="video-controls-button"
                    onClick={back}> 
                    Back a little (Left Arrow)
                </div>

                <div className="video-controls-button"
                    onClick={forward}> 
                    Forward a little (Right Arrow)
                </div>

                <div className="video-controls-button"
                    onClick={() => confirmClip(true)}> 
                    Clip Skateboard (Enter)
                </div>

                <div className="video-controls-button"
                    onClick={() => confirmClip(false)}> 
                    Clip Not-a-Skateboard (Shift Enter)
                </div>

                <div className="video-controls-button"
                    onClick={setUnitQuaternion}>
                    Reset Skateboard Orientation (R)
                </div>

            </div>
        </div>
    )
}