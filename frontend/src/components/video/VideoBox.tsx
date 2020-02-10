import * as React from 'react';

import { Box } from './box';
import { Bounds } from '../../api';

interface VideoBoxProps {
    video: React.RefObject<HTMLVideoElement>
    layout: Bounds
    videoWidth: number
    videoHeight: number

    onSubmit: (bounds: Bounds) => void
}

export function VideoBox(props: VideoBoxProps) {
    let canvas = React.createRef<HTMLCanvasElement>()

    // bind play/pause to our state variable
    let [ isPlaying, setIsPlaying ] = React.useState(false);
    React.useEffect(() => {
        if (!props.video.current) {
            return;
        }
        function onPause() {
            setIsPlaying(false);
        }
        function onPlay() {
            setIsPlaying(true);
        }
        props.video.current.addEventListener("play", onPlay);
        props.video.current.addEventListener("pause", onPause);
        return () => {
            props.video.current?.removeEventListener("play", onPlay);
            props.video.current?.removeEventListener("pause", onPause);
        }
    }, [props.video.current, canvas.current])

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
            box.release(props.onSubmit, () => props.video.current?.play());
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
            box.release(props.onSubmit, () => props.video.current?.play());
        }
    }, [canvas.current, props.layout, props.video.current])

    return <canvas
        ref={canvas}
        width={props.videoWidth}
        height={props.videoHeight}
        style={{
            ...props.layout,
            position: "absolute",
            display: isPlaying ? "none" : "block"
        }}
    > </canvas>
}