import * as React from 'react';

interface VideoControlProps {
    video: React.RefObject<HTMLVideoElement>
    bounds: {
        top: number,
        left: number,
        width: number
    }
}

export function VideoScrubber({ bounds, video }: VideoControlProps) {

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
        if (!video.current) {
            return;
        }
        function onTimeUpdate(this: HTMLVideoElement) {
            setTime(this.currentTime);
            setBuffer(this.buffered);
            setDuration(this.duration);
        }
        video.current.addEventListener("timeupdate", onTimeUpdate);
        return () => video.current?.removeEventListener("timeupdate", onTimeUpdate);
    }, [video.current])

    const scrub = (e: {clientX: number}) => {
        render();
        if (!canvas.current || !video.current) return
        let { x } = canvas.current.getBoundingClientRect();
        let time = video.current.duration * (e.clientX - x) / bounds.width;
        video.current.currentTime = time;
        setTime(time);
    }

    const render = () => {
        let ctx = canvas.current?.getContext("2d");
        if (!ctx) {
            return;
        }
        let height = (new Date().getTime() - lastInteraction.getTime()) > WIDE_SCRUBBER_DURATION ? 5 : 30;
        let y = (30 - height) / 2
        let pixelsPerSecond = bounds.width / duration;
        ctx.clearRect(0, 0, bounds.width, 30);
        if (buffer) {
            ctx.fillStyle = "grey";
            for (let i = 0; i < buffer.length; i++) {
                ctx.fillRect(buffer.start(i) * pixelsPerSecond, y, buffer.end(i) * pixelsPerSecond, height);
            }
        }
        ctx.fillStyle = "#33b5e5";
        ctx.fillRect(0, y, time * pixelsPerSecond, height);
    }

    return <canvas
        ref={canvas}
        width={bounds.width}
        onMouseDown={_ => setDragging(true)}
        onMouseMove={e => {
            setLastInteraction(new Date());
            if (dragging) {
                scrub(e);
            }
        }}
        onMouseUp={_ => setDragging(false)}
        onMouseLeave={_ => setDragging(false)}
        onMouseOut={_ => setDragging(false)}
        onTouchStart={_ => setLastInteraction(new Date())}
        onTouchMove={e => {
            setLastInteraction(new Date());
            scrub(e.targetTouches[0])
        }}
        height={30}
        style={{
            ...bounds,
            position: "absolute",
            height: 30
        }}
    ></canvas>

}