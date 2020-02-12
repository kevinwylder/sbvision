import * as React from 'react';

interface VideoControlProps {
    video: React.RefObject<HTMLVideoElement>
    bounds: {
        top: number,
        left: number,
        width: number
    }
    height: number
}

export function VideoScrubber(props: VideoControlProps) {

    const WIDE_SCRUBBER_DURATION = 2000;

    let [ time, setTime ] = React.useState(0);
    let [ buffer, setBuffer ] = React.useState<TimeRanges>();
    let [ duration, setDuration ] = React.useState(0);
    let [ dragging, setDragging ] = React.useState(false);
    let [ lastInteraction, setLastInteraction ] = React.useState(new Date())
    let canvas = React.createRef<HTMLCanvasElement>()
    React.useEffect(() => {
        render();
    }, [time, buffer]);

    React.useEffect(() => {
        render();
        let timeout = window.setTimeout(() => render(), WIDE_SCRUBBER_DURATION);
        return () => window.clearTimeout(timeout);
    }, [lastInteraction]);

    React.useEffect(() => {
        if (!props.video.current) {
            return;
        }
        function onTimeUpdate(this: HTMLVideoElement) {
            setTime(this.currentTime);
            setBuffer(this.buffered);
            setDuration(this.duration);
        }
        props.video.current.addEventListener("timeupdate", onTimeUpdate);
        return () => props.video.current?.removeEventListener("timeupdate", onTimeUpdate);
    }, [props.video.current])

    const scrub = (e: {clientX: number}) => {
        render();
        if (!canvas.current || !props.video.current) return
        let { x } = canvas.current.getBoundingClientRect();
        let time = props.video.current.duration * (e.clientX - x) / props.bounds.width;
        props.video.current.currentTime = time;
        setTime(time);
    }

    const render = () => {
        let ctx = canvas.current?.getContext("2d");
        if (!ctx) {
            return;
        }
        let height = (new Date().getTime() - lastInteraction.getTime()) > WIDE_SCRUBBER_DURATION ? 5 : props.height;
        let pixelsPerSecond = props.bounds.width / duration;
        ctx.clearRect(0, 0, props.bounds.width, props.height);
        if (buffer) {
            ctx.fillStyle = "grey";
            for (let i = 0; i < buffer.length; i++) {
                ctx.fillRect(buffer.start(i) * pixelsPerSecond, 0, buffer.end(i) * pixelsPerSecond, height);
            }
        }
        ctx.fillStyle = "#33b5e5";
        ctx.fillRect(0, 0, time * pixelsPerSecond, height);
    }

    return <canvas
        ref={canvas}
        onMouseDown={e => {
            setDragging(true);
            setLastInteraction(new Date());
            scrub(e);
        }}
        onMouseLeave={_ => setDragging(false)}
        onMouseOut={_ => setDragging(false)}
        onMouseMove={e => {
            setLastInteraction(new Date());
            if (dragging) {
                scrub(e);
            }
        }}
        onMouseUp={_ => setDragging(false)}
        onTouchStart={_ => setLastInteraction(new Date())}
        onTouchMove={e => {
            setLastInteraction(new Date());
            scrub(e.targetTouches[0])
        }}
        width={props.bounds.width}
        height={props.height}
        style={{
            ...props.bounds,
            position: "absolute",
            height: props.height
        }}
    ></canvas>

}