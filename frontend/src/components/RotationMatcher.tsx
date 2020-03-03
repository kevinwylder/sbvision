import * as React from 'react';
import { Frame, RotationlessFramePager } from '../api';
import { Quaternion, renderSkateboard, rotateSkateboard, tiltSkateboard } from '../renderer/skateboard';
import { Link } from 'react-router-dom';


export function RotationMatcher() {

    // compute the remaining height of the window
    let instructions = React.createRef<HTMLDivElement>();
    let [ height, setHeight ] = React.useState(0);
    React.useEffect(() => {
        if (!instructions.current) {
            return
        }
        let { bottom } = instructions.current.getBoundingClientRect();
        setHeight(window.innerHeight - bottom);
    }, [instructions.current])

    // bound paging
    let [ frameManager, _ ] = React.useState(new RotationlessFramePager());
    let [ frame, setFrame ] = React.useState<Frame|undefined>();
    let [ boundIdx, setBoundIdx ] = React.useState(0)
    let [ hasFrames, setHasFrames ] = React.useState(true);
    React.useEffect(() => {
        if (!frame || boundIdx >= frame.bounds.length) {
            frameManager.nextFrame()
            .then(frame => {
                if (frame) {
                    setBoundIdx(0);
                    setFrame(frame);
                } else {
                    setHasFrames(false);
                }
            })
            .catch(err => console.log(err));
        }
    }, [frame, boundIdx]);
    
    let [ image, setImage ] = React.useState<HTMLImageElement>()
    React.useEffect(() => {
        if (!frame) return;
        let i = document.createElement("img");
        i.src = `/api/image?frame=${frame.id}`;
        i.onload = () => {
            setImage(i);
        }
    }, [frame]);

    let [ quaternion, setQuaternion ] = React.useState<Quaternion>([0, Math.SQRT1_2, Math.SQRT1_2, 0]);
    let canvas = React.createRef<HTMLCanvasElement>()
    React.useEffect(() => {
        if (!canvas.current || !image || !frame || frame.bounds.length <= boundIdx) return;
        let ctx = canvas.current.getContext("2d");
        if (!ctx) return;
        let { x, y, width, height } = frame.bounds[boundIdx];
        ctx.drawImage(image, 0, 0);
        ctx.strokeStyle = "black";
        ctx.lineWidth = 5;
        ctx.strokeRect(x, y, width, height);
        ctx.strokeStyle = "#33b5e5";
        ctx.lineWidth = 3;
        ctx.strokeRect(x, y, width, height);
        let long = Math.min(image.width, image.height) * .25;
        let offsetX = (x + x + width < image.width) ? width : -long;
        let offsetY = (y + y + height < image.height) ? 0 : height - long;
        renderSkateboard(ctx, quaternion, [ x + offsetX, y + offsetY, x + offsetX + long, y + offsetY + long]);
    }, [canvas.current, image, boundIdx, quaternion])

    React.useEffect(() => {
        function deleteBounds(e: KeyboardEvent) {
            if (e.key == "Backspace") {
                setBoundIdx(boundIdx + 1);
                frameManager.skipRotation();
            }
        }
        window.addEventListener("keypress", deleteBounds);
        return () => window.removeEventListener("keypress", deleteBounds);
    }, []);


    const handleClick = () => {
        if (!frame) return;
        let [ r, i, j, k ] = quaternion;
        frameManager.addRotation(frame.bounds[boundIdx].id, {r, i, j, k})
        .then(() => {
            setBoundIdx(boundIdx + 1);
        })
        .catch(err => console.log(err));
    }

    let [ lastPosition, setLastPosition ] = React.useState<[number, number]>([0, 0])
    const handleMove = (e: React.MouseEvent) => {
        let [ lastX, lastY ] = lastPosition;
        if (lastX != 0 || lastY != 0) {
            setQuaternion(rotateSkateboard(e.clientX - lastX, e.clientY - lastY, quaternion));
        }
        setLastPosition([e.clientX, e.clientY]);
    }

    const handleWheel = (e: React.WheelEvent) => { 
        setQuaternion(tiltSkateboard(e.deltaY, quaternion));
    }

    if (!hasFrames) {
        return <div style={{margin: "10px"}}> 
            <h1> No more rotations to match!</h1>
            <p> Add new bounds <Link to="/videos"> by pausing videos </Link> to make roations show up here </p>
        </div>
    }

    return <div style={{
        display: "flex",
        flexDirection: "column",
        width: "100%",
    }} >
        <div ref={instructions}>
            Move the mouse and scroll to find the matching skateboard orientation. Click to accept the rotation, or press backspace to mark the bounds as "not a skateboard".
        </div>
        <canvas
            onClick={handleClick}
            ref={canvas}
            onMouseMove={handleMove}
            onWheel={handleWheel}
            style={{
                height,
                display: (image) ? "block" : "none",
                objectFit: "contain",
                width: "100%",
                flexGrow: 1
            }}
            width={image?.width}
            height={image?.height} />
    </div>
}