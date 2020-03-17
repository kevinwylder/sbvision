import * as React from 'react';
import { useParams } from 'react-router-dom';

import { VideoWrapper } from './videowrapper';

export function VideoDisplay() {

    // compute the remaining height of the window
    let container = React.createRef<HTMLDivElement>();
    let [ height, setHeight ] = React.useState(0);
    React.useEffect(() => {
        if (!container.current) {
            return
        }
        let { top } = container.current.getBoundingClientRect();
        setHeight(window.innerHeight - top - 1);
    }, [container.current])

    // load the renderer and get it's size
    let params: { id?: string } = useParams();
    let [ renderer, setRenderer ] = React.useState<VideoWrapper>()
    let [ videoWidth, setVideoWidth ] = React.useState(0);
    let [ videoHeight, setVideoHeight ] = React.useState(0);
    React.useEffect(() => {
        let renderer = new VideoWrapper(parseInt(params.id as string), (width, height) => {
            // relay size information back to state for event positioning
            setVideoWidth(width) 
            setVideoHeight(height);
            setRenderer(renderer);
        });

        const keyListener = function(e: KeyboardEvent) {
            console.log(e.key);
            switch (e.key) {
                case "ArrowLeft":
                    renderer.back();
                    break;
                case "ArrowRight":
                    renderer.forward();
                    break;
            }
        }

        window.addEventListener("keydown", keyListener);

        return () => window.removeEventListener("keydown", keyListener);
    }, [])

    // route events to the renderer
    let canvas = React.createRef<HTMLCanvasElement>();
    React.useEffect(() => {
        if (!canvas.current || !renderer) {
            return;
        }
        const ctx = canvas.current.getContext("2d");
        if (!ctx) {
            return;
        }
        renderer.setContext(ctx);

        const { top } = canvas.current.getBoundingClientRect();
        const convertCoordinates = ({clientX, clientY}: {clientX: number, clientY: number}): [number, number] => {
            let t = Math.max(height - (videoHeight * window.innerWidth) / videoWidth, 0) / 2;
            let l = Math.max(window.innerWidth - (videoWidth * height) / videoHeight, 0) / 2;
            let c = videoHeight / ( height - 2 * t );
            return [
                c * (clientX - l),
                c * (clientY - t - top),
            ];
        }

        canvas.current.onmousedown = (e) => {
            if (renderer?.grab(convertCoordinates(e))) {
                e.preventDefault();
            }
        }
        canvas.current.onmousemove = (e) => {
            if (renderer?.drag(convertCoordinates(e))) {
                e.preventDefault();
            }
        }
        canvas.current.onmouseup = (e) => {
            if (renderer?.up()) {
                e.preventDefault();
            }
        }
        canvas.current.ontouchstart = (e) => {
            if (renderer?.grab(convertCoordinates(e.targetTouches[0]))) {
                e.preventDefault();
            }
        }
        canvas.current.ontouchmove = (e) => {
            if (renderer?.drag(convertCoordinates(e.targetTouches[0]))) {
                e.preventDefault();
            }
        }
        canvas.current.ontouchend = (e) => {
            if (renderer?.up()) {
                e.preventDefault();
            }
        }


    }, [canvas.current, renderer]);

    // remove the renderer from the dom tree when done
    React.useEffect(() => {
        const oldRenderer = renderer;
        return () => {
            if (oldRenderer) {
                oldRenderer.destroy();
            }
        }
    }, [renderer])

    return <div 
        ref={container}
        className="video-container">
        <canvas 
            ref={canvas}
            style={{ 
                height,
                objectFit: "contain",
                width: "100%",
             }}
            width={videoWidth}
            height={videoHeight}
        />
        {videoWidth == 0 ? "Loading..." : ""}
    </div> 
}