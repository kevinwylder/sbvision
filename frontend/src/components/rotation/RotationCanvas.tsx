import * as React from 'react';

import { renderSkateboard, Quaternion, rotateSkateboard, tiltSkateboard } from './skateboard';
import { Frame, addRotation } from '../../api';

interface RotationCanvasProps {
    frame: Frame
    style: React.CSSProperties
    onFrameComplete: () => void
}

export function RotationCanvas({frame, style, onFrameComplete}: RotationCanvasProps) {
    
    let [ quaternion, setQuaternion ] = React.useState<Quaternion>([0, Math.SQRT1_2, Math.SQRT1_2, 0]);
    let [ bound, setBound ] = React.useState(0)
    let [ image, setImage ] = React.useState<HTMLImageElement>()
    React.useEffect(() => {
        let i = document.createElement("img");
        i.src = `/images/frame/${frame.id}.png`;
        i.onload = () => {
            setImage(i);
        }
    }, [frame.id]);

    let canvas = React.createRef<HTMLCanvasElement>()
    React.useEffect(() => {
        if (!canvas.current || !image) return;
        let ctx = canvas.current.getContext("2d");
        if (!ctx) return;
        let { x, y, width, height } = frame.bounds[bound];
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
    }, [canvas.current, image, bound, quaternion])

    React.useEffect(() => {
        function deleteBounds(e: KeyboardEvent) {
            if (e.key == "Backspace") {
            }
        }
        window.addEventListener("keypress", deleteBounds);
        return () => window.removeEventListener("keypress", deleteBounds);
    }, []);

    const handleClick = () => {
        let [ r, i, j, k ] = quaternion;
        addRotation(frame.bounds[bound].id, {r, i, j, k})
        .then(() => {
            if (bound >= frame.bounds.length - 1) {
                onFrameComplete();
                return
            }
            setBound(bound + 1);
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
    
    return <canvas
        onClick={handleClick}
        ref={canvas}
        onMouseMove={handleMove}
        onWheel={handleWheel}
        style={{
            display: (image) ? "block" : "none",
            objectFit: "contain",
            ...style,
        }}
        width={image?.width}
        height={image?.height} />

}