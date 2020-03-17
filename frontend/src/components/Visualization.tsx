import * as React from 'react';
import { VisualizationManager } from '../api';
import { Quaternion, rotateSkateboard, tiltSkateboard, renderSkateboard } from '../renderer';

export function DataVisualization() {
    // map canvas mouse events to the websocket
    let canvas = React.createRef<HTMLCanvasElement>();
    React.useEffect(() => {
        const ctx = canvas.current?.getContext('2d');
        if (!ctx || !canvas.current) return

        const { top } = canvas.current.getBoundingClientRect();
        canvas.current.height = window.innerHeight - top;

        let img: HTMLImageElement|undefined;
        let imgQ: Quaternion = [0, 0, 0, 0];
        let quaternion: Quaternion = [1, 0, 0, 0];
        const draw = () => {
            if (img) {
                let width = Math.min(img.width / img.height, 1) * 500;
                let height = Math.min(img.height / img.width, 1) * 500;
                let top = (500 - height) / 2;
                let left = (500 - width) / 2;
                ctx.clearRect(0, 0, 500, 500);
                ctx.drawImage(img, left, top, width, height);
                let err = (imgQ[0] - quaternion[0]) * (imgQ[0] - quaternion[0]) + 
                          (imgQ[1] - quaternion[1]) * (imgQ[1] - quaternion[1]) + 
                          (imgQ[2] - quaternion[2]) * (imgQ[2] - quaternion[2]) + 
                          (imgQ[3] - quaternion[3]) * (imgQ[3] - quaternion[3]);
                ctx.fillStyle = `rgba(255, 255, 255, ${err})`;
                ctx.fillRect(0, 0, 500, 500);
            }
            renderSkateboard(ctx, quaternion, [500, 0, 1000, 500]);
        }
        draw();

        let websocket = new VisualizationManager((image, imageQ) => {
            img = image;
            imgQ = imageQ;
            draw();
        });


        let lastPosition: [number, number] = [0, 0];
        canvas.current.onmousemove = (e: MouseEvent) => {
            let [ lastX, lastY ] = lastPosition;
            if (lastX != 0 || lastY != 0) {
                quaternion = rotateSkateboard(e.clientX - lastX, e.clientY - lastY, quaternion);
                websocket.setRotation(quaternion);
                draw();
            }
            lastPosition = [e.clientX, e.clientY];
        }

        canvas.current.onwheel = (e: WheelEvent) => { 
            e.preventDefault();
            quaternion = tiltSkateboard(e.deltaY, quaternion);
            websocket.setRotation(quaternion);
            draw();
        }

        canvas.current.ontouchstart = () => {
            lastPosition = [0, 0];
        }
        canvas.current.ontouchmove = (e: TouchEvent) => {
            e.preventDefault();
            if (e.touches.length > 1) {
                let y = (e.touches[0].clientY + e.touches[1].clientY) / 2;
                if (lastPosition[1] != 0) {
                    let dy = y - lastPosition[1];
                    quaternion = tiltSkateboard(dy, quaternion);
                    websocket.setRotation(quaternion);
                    draw();
                } 
                lastPosition = [0, y];
            } else {
                if (lastPosition[1] != 0) {
                    let [ x, y ] = lastPosition;
                    let dx = e.touches[0].clientX - x;
                    let dy = e.touches[0].clientY - y;
                    quaternion = rotateSkateboard(dx, dy, quaternion);
                    websocket.setRotation(quaternion);
                    draw();
                }
                lastPosition = [e.touches[0].clientX, e.touches[0].clientY];
            }
        }

        return () => {
            websocket.disconnect();
        }
    }, [canvas.current]);

    return <canvas
            ref={canvas}
            style={{
                objectFit: "contain",
                width: "100%",
            }}
            width={1000}
            height={500} />
}