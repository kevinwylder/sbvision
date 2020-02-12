import * as React from 'react';
import { useParams } from 'react-router-dom';

import { Video, getVideoById } from '../../api';

import { VideoManager } from './manager';

import { VideoBox } from './VideoBox';
import { VideoScrubber } from './VideoScrubber';
import { VideoControls } from './VideoControls';

import "./video.css";
import { CollectionStatistics } from '../CollectionStats';

const CONTROLS_HEIGHT = 49;
const SCRUBBER_HEIGHT = 30;

export function VideoDisplay() {

    // load the video info based on the route
    let [ videoInfo, setVideoInfo ] = React.useState<Video>()
    let params: { id?: string } = useParams();
    React.useEffect(() => {
        getVideoById(parseInt(params.id as string))
        .then(v => setVideoInfo(v.videos[0]));
    }, [params.id])

    // use aspect ratio and container bounds to place the video
    let container = React.createRef<HTMLDivElement>();
    let [ bounds, setDisplayBounds ] = React.useState({ top: 0, left: 0, width: 0, height: 0 });
    let [ videoSize, setVideoSize ] = React.useState([100, 100]);
    React.useEffect(() => {
        if (!container.current) return;
        let [ vWidth, vHeight ] = videoSize;
        let aspectRatio = vWidth / vHeight;
        let { width, height } = container.current.getBoundingClientRect();
        let displayLeft = Math.max((width - aspectRatio * (height - SCRUBBER_HEIGHT - CONTROLS_HEIGHT)), 0) / 2;
        let displayWidth = width - 2 * displayLeft;
        let displayHeight = displayWidth / aspectRatio;
        let displayTop = Math.max((height - displayHeight - SCRUBBER_HEIGHT - CONTROLS_HEIGHT) / 2 + CONTROLS_HEIGHT, CONTROLS_HEIGHT) ;
        setDisplayBounds({
            top: displayTop,
            left: displayLeft,
            width: displayWidth,
            height: displayHeight,
        });
    }, [container.current, videoSize]);

    // get state callbacks from video events 
    let video = React.createRef<HTMLVideoElement>();
    let prestigeCanvas = React.createRef<HTMLCanvasElement>();
    let [ hasPlayed, setHasPlayed ] = React.useState(false);
    let [ collectionStats, setCollectionStats ] = React.useState<CollectionStatistics>({bounds: 0, frames: 0, rotations: 0});
    let [ videoManager, setVideoManager ] = React.useState<VideoManager>()
    React.useEffect(() => {
        if (!video.current || !prestigeCanvas.current || !videoInfo) {
            return;
        }
        let ctx = prestigeCanvas.current.getContext("2d");
        if (!ctx) {
            return;
        }
        let videoManager = new VideoManager(videoInfo.id, video.current, ctx, setCollectionStats);
        setVideoManager(videoManager);
        function onPlay() {
            setHasPlayed(true);
        }
        function onResize({target}: UIEvent) {
            let v = target as HTMLVideoElement;
            setVideoSize([v.videoWidth, v.videoHeight]);
        }
        video.current.addEventListener("resize", onResize);
        video.current.addEventListener("play", onPlay);
        video.current.load();
        return () => {
            video.current?.removeEventListener("resize", onResize)
            video.current?.removeEventListener("play", onPlay);
            videoManager.unregister();
        };
    }, [ video.current, prestigeCanvas.current, videoInfo ]);

    return <div 
        ref={container}
        className="video-container">
{videoInfo && collectionStats ? <>
        <VideoControls
            stats={collectionStats}
            bounds={{
                left: bounds.left,
                top: bounds.top - CONTROLS_HEIGHT,
                width: bounds.width,
                height: CONTROLS_HEIGHT
            }}
        />
        <video
            ref={video} 
            playsInline={true}
            autoPlay={true}
            controls={false}
            style={{
                position: "absolute",
                ...bounds
            }}>
            <source
                src={`/stream?type=${videoInfo.type}&id=${videoInfo.id}`}
                type={videoInfo.format}
            />
        </video>
        <canvas 
            ref={prestigeCanvas}
            width={videoSize[0]}
            height={videoSize[1]}
            style={{
                position: "absolute",
                ...bounds
            }}
            onClick={_ => video.current?.pause()}
        />
        <VideoBox 
            layout={bounds}
            video={video}
            videoWidth={videoSize[0]}
            videoHeight={videoSize[1]}
            onSubmit={(bounds) => {
                videoManager?.giveBounds(bounds)
            }}
            onReject={() => {
                videoManager?.continue();
            }}
        />
        <VideoScrubber
            height={SCRUBBER_HEIGHT}
            video={video}
            bounds={{
                top: bounds.top + bounds.height,
                left: bounds.left,
                width: bounds.width
            }}
            ></VideoScrubber>
        <img
            src="/play.png"
            onClick={() => video.current?.play()}
            style={{
                ...bounds,
                position: "absolute",
                objectFit: "contain",
                display: hasPlayed ? "none" : "block"
            }}
        />
</> : <div>Loading Video Info....</div>}
    </div> 
}