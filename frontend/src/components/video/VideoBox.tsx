import * as React from 'react';

import { Box, Bounds } from './box';

interface VideoBoxProps {
    video: React.RefObject<HTMLVideoElement>
    layout: Bounds
    videoWidth: number
    videoHeight: number

    setHasPlayed: (hasPlayed: boolean) => void
    onPause: (data: string) => void
    onSubmit: (bounds: Bounds) => void
    onRefuse: () => void
}

export function VideoBox(props: VideoBoxProps) {
    let canvas = React.createRef<HTMLCanvasElement>()

    // bind play/pause to our state variable
    let [ isPlaying, setIsPlaying ] = React.useState(false);
    React.useEffect(() => {
        if (!props.video.current) {
            return;
        }
        props.video.current.onpause = ({target}) => {
            let v = target as HTMLVideoElement;
            let data = getVideoData(v);
            if (data) {
                props.onPause(data);
            } else {
                props.onRefuse();
            }
            setIsPlaying(false);
        }
        props.video.current.onplay = () => {
            props.setHasPlayed(true);
            setIsPlaying(true);
        }
    }, [props.video.current, canvas.current])

    // apply our state to the canvas
    React.useEffect(() => {
        if (!canvas.current) {
            return;
        }
        canvas.current.style.display = isPlaying ? "none" : "block";
    }, [canvas.current, isPlaying])

    // bind canvas events to canvas box
    React.useEffect(() => {
        if (!canvas.current) {
            return;
        }
        let ctx = canvas.current.getContext("2d");
        if (!ctx) {
            return;
        }
        let box = new Box(props.videoWidth, props.videoHeight, 10, ctx);
        const convertCoordinates = ({clientX, clientY}: {clientX: number, clientY: number}): [number, number] => {
            let { top, left, width, height } = props.layout;
            return [
                props.videoWidth * (clientX - left) / width,
                props.videoHeight * (clientY - top) / height,
            ]
        }
        canvas.current.onmousedown = (e) => {
            e.preventDefault();
            box.grab(convertCoordinates(e), "click");
        }
        canvas.current.onmousemove = (e) => {
            e.preventDefault();
            box.drag(convertCoordinates(e));
        }
        canvas.current.onmouseup = (e) => {
            e.preventDefault();
            box.release(props.onSubmit, props.onRefuse);
        }
        canvas.current.onmouseleave = (e) => {
            e.preventDefault();
            box.release(() => {}, () => {});
        }
        canvas.current.onmouseout = (e) => {
            e.preventDefault();
            box.release(() => {}, () => {});
        }
        canvas.current.ontouchstart = (e) => {
            e.preventDefault();
            box.grab(convertCoordinates(e.targetTouches[0]), "tap");
        }
        canvas.current.ontouchmove = (e) => {
            e.preventDefault();
            box.drag(convertCoordinates(e.targetTouches[0]));
        }
        canvas.current.ontouchend = (e) => {
            e.preventDefault();
            box.release(props.onSubmit, props.onRefuse);
        }
    }, [canvas.current, props.layout])

    return <canvas
        ref={canvas}
        width={props.videoWidth}
        height={props.videoHeight}
        style={{
            ...props.layout,
            position: "absolute",
            display: "none"
        }}
    > </canvas>
}

function getVideoData(v: HTMLVideoElement) {
    let canvas = document.createElement("canvas");
    canvas.width = v.videoWidth;
    canvas.height = v.videoHeight;
    let ctx = canvas.getContext("2d");
    if (!ctx) {
        return;
    }
    ctx.drawImage(v, 0, 0);
    return canvas.toDataURL();
}