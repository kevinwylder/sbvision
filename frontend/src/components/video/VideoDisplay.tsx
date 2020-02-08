import * as React from 'react';

import "../../styles/video.css";

import { Box } from './box';
import { Video } from '../../api';
import { VideoBox } from './VideoBox';
import { VideoScrubber } from './VideoScrubber';

interface VideoDisplayProps {
    video: Video
}

export function VideoDisplay(props: VideoDisplayProps) {

    let canvas = React.createRef<HTMLCanvasElement>();
    let video = React.createRef<HTMLVideoElement>();

    // get state callbacks from video events 
    let [ videoSize, setVideoSize ] = React.useState([100, 100]);
    React.useEffect(() => {
        if (!video.current) {
            return;
        }
        video.current.onresize = ({target}) => {
            let v = target as HTMLVideoElement;
            setVideoSize([v.videoWidth, v.videoHeight]);
        }
        video.current.load();
    }, [ video.current, canvas.current ]);

    // use aspect ratio and container bounds to place the video
    let container = React.createRef<HTMLDivElement>();
    let [ bounds, setDisplayBounds ] = React.useState({ top: 0, left: 0, width: 0, height: 0 });
    React.useEffect(() => {
        if (!container.current) return;
        let [ vWidth, vHeight ] = videoSize;
        let aspectRatio = vWidth / vHeight;
        let { width, height } = container.current.getBoundingClientRect();
        let displayLeft = Math.max((width - aspectRatio * height), 0) / 2;
        let displayWidth = width - 2 * displayLeft;
        let displayHeight = displayWidth / aspectRatio;
        let displayTop = Math.max((height - displayHeight) / 2, 0);
        setDisplayBounds({
            top: displayTop,
            left: displayLeft,
            width: displayWidth,
            height: displayHeight,
        });
    }, [container.current, videoSize]);

    let [ hasPlayed, setHasPlayed ] = React.useState(false);
    let [ frame, setFrame ] = React.useState(0);

    // oh god this is where it starts to fall apart
    let [ advanceFrame, setAdvanceFrame ] = React.useState(0);
    let [ playFlag, setPlayFlag ] = React.useState(false)
    React.useEffect(() => {
        if (advanceFrame && video.current) {
            video.current.currentTime += advanceFrame / props.video.fps;
            setAdvanceFrame(0);
        }
        if (playFlag && video.current) {
            video.current.play();
            setPlayFlag(false);
        }
    }, [advanceFrame, playFlag, video.current])

    return <div 
        ref={container}
        className="video-container">
        <video
            ref={video} 
            onClick={_ => video.current?.pause()}
            style={{
                position: "absolute",
                ...bounds
            }}>
            <source
                src={`/video?type=${props.video.type}&id=${props.video.id}`}
                type={props.video.format}
            />
        </video>
        <VideoBox 
            layout={bounds}
            video={video}
            setHasPlayed={setHasPlayed}
            videoWidth={videoSize[0]}
            videoHeight={videoSize[1]}
            onPause={(data) => {

            }}
            onSubmit={(bounds) => {
                setAdvanceFrame(1);
            }}
            onRefuse={() => {
                setPlayFlag(true);
            }}
        />
        <VideoScrubber
            video={video}
            fps={props.video.fps}
            onFrame={setFrame}
            bounds={{
                top: bounds.top + bounds.height - 15,
                left: bounds.left,
                width: bounds.width
            }}
            ></VideoScrubber>
        <img
            style={{
                ...bounds,
                position: "absolute",
                objectFit: "contain",
                display: hasPlayed ? "none" : "block"
            }}
            src={`/images/${props.video.thumbnail}`}
            onClick={() => video.current && video.current.play() }
        />
    </div> 
}