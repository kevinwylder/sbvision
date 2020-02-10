import * as React from 'react';

import "../../styles/video.css";

import { Video, FrameManager, getVideoById } from '../../api';
import { VideoBox } from './VideoBox';
import { VideoScrubber } from './VideoScrubber';
import { useParams } from 'react-router-dom';

interface VideoDisplayProps {
    video?: Video
}

interface RouteParams {
    id?: string
}

export function VideoDisplay(props: VideoDisplayProps) {

    // allow for route based video loading
    let [ videoInfo, setVideoInfo ] = React.useState<Video>()
    let params: RouteParams = useParams();

    React.useEffect(() => {
        if (props.video) {
            setVideoInfo({...props.video});
            return;
        }
        getVideoById(parseInt(params.id as string))
        .then(v => setVideoInfo(v.videos[0]));
    }, [props.video])

    let video = React.createRef<HTMLVideoElement>();

    // get state callbacks from video events 
    let [ hasPlayed, setHasPlayed ] = React.useState(false);
    let [ videoSize, setVideoSize ] = React.useState([100, 100]);
    let [ frameManager, setFrameManager ] = React.useState<FrameManager>();
    React.useEffect(() => {
        if (!video.current || !videoInfo) {
            return;
        }
        let frameManager = new FrameManager(videoInfo, video.current)
        setFrameManager(frameManager);
        function onPaused() {
            frameManager?.uploadFrame();
        }
        function onPlay() {
            setHasPlayed(true);
        }
        function onResize({target}: UIEvent) {
            let v = target as HTMLVideoElement;
            setVideoSize([v.videoWidth, v.videoHeight]);
        }
        video.current.addEventListener("resize", onResize);
        video.current.addEventListener("play", onPlay);
        video.current.addEventListener("pause", onPaused);
        video.current.load();
        return () => {
            video.current?.removeEventListener("resize", onResize)
            video.current?.removeEventListener("play", onPlay);
            video.current?.removeEventListener("pause", onPaused);
        };
    }, [ video.current, videoInfo ]);

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

    return <div 
        ref={container}
        className="video-container">
{videoInfo ? <>
        <video
            ref={video} 
            onClick={_ => video.current?.pause()}
            style={{
                position: "absolute",
                ...bounds
            }}>
            <source
                src={`/stream?type=${videoInfo.type}&id=${videoInfo.id}`}
                type={videoInfo.format}
            />
        </video>
        <VideoBox 
            layout={bounds}
            video={video}
            videoWidth={videoSize[0]}
            videoHeight={videoSize[1]}
            onSubmit={(bounds) => {
                frameManager?.giveBounds(bounds);
            }}
        />
        <VideoScrubber
            video={video}
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
            src={`/images/${videoInfo.thumbnail}`}
            onClick={() => video.current && video.current.play() }
        />
</> : <div>Loading Video Info....</div>}
    </div> 
}