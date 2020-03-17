import * as React from 'react';
import './verify.css';

import { useParams } from 'react-router-dom';
import { ImageWrapperFrameStatus, ImageWrapper } from './imagewrapper';
import { BoundsList } from './BoundRow';

export function RotationMatcher() {

    // setup the imagewrapper and callbacks for the current frame
    let canvas = React.createRef<HTMLCanvasElement>()
    let [ status, setStatus ] = React.useState<ImageWrapperFrameStatus>();
    let [ renderer, setRenderer ] = React.useState<ImageWrapper>();
    let params: { id?: string } = useParams();
    React.useEffect(() => {
        if (!canvas.current) {
            return
        }
        if (!params.id) {
            console.log("no frameid")
            return
        }
        let renderer = new ImageWrapper(parseInt(params.id), canvas.current, setStatus);
        setRenderer(renderer)

        canvas.current.onclick = (e) => {
            renderer.click(e);
            e.preventDefault();
        }
        canvas.current.onmousemove = (e) => { 
            renderer.move(e);
            e.preventDefault();
        };
        canvas.current.onwheel = (e) => { 
            renderer.scroll(e.deltaY);
            e.preventDefault();
        };

        function keyListener(e: KeyboardEvent) {
            switch (e.key) {
                case "Backspace":
                    renderer.remove();
                    break
                case "ArrowRight":
                    renderer.next();
                    break
                case "ArrowLeft":
                    renderer.prev();
                    break
            }
        }
        window.addEventListener('keydown', keyListener);

        return () => {
            renderer.destroy();
            window.removeEventListener('keydown', keyListener);
        }
    }, [canvas.current]);

    return <div style={{
        display: "grid",
        gridTemplateColumns: "1fr 400px",
        width: "100%",
    }} >
        <canvas
            ref={canvas}
            style={{
                position: "sticky",
                top: "0px",
                objectFit: "contain",
                width: "100%",
                flexGrow: 1
            }} />
        <div >
            <BoundsList
                bounds={status?.frame?.bounds || []} 
                selected={status?.selectedBound || 0}
                onclick={(id) => { renderer?.setBound(id)} }/>
        </div>
    </div>
}