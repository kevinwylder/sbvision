import * as React from 'react';

import { renderSkateboard } from '../skateboard';

import { Quaternion, ModelInput, ModelOutput, eToQ, qMultiply, qRotate } from '../model';

interface RotationMatcherProps{
    input: ModelInput
    output?: ModelOutput
    onOutputConfirmed: (output: ModelOutput) => void
}

export function RotationMatcher(props: RotationMatcherProps) {

    let [ isSkateboard, setIsSkateboard ] = React.useState(1);
    let [ quaternion, setQuaternion ] = React.useState<Quaternion>([Math.SQRT1_2, Math.SQRT1_2, 0, 0]);
    let [ lastMousePosition, setLastMousePosition ] = React.useState<[number, number]>([0, 0]);
    let [ isOffAxisRotation, setIsOffAxisRotation ] = React.useState(false);
    let [ isDragging, setIsDragging ] = React.useState(false);

    let canvasRef = React.createRef<HTMLCanvasElement>();

    React.useEffect(() => {
        if (!canvasRef.current) {
            return;
        }
        let { width, height } = canvasRef.current;
        let ctx = canvasRef.current.getContext("2d");
        if (!ctx) { 
            return;
        }
        ctx.clearRect(0, 0, width, height);
        if (props.input.width > props.input.height) {
            ctx.putImageData(props.input, 0, (props.input.width - props.input.height) / 2);
        } else {
            ctx.putImageData(props.input, (props.input.height - props.input.width) / 2, 0);
        }

        renderSkateboard(ctx, quaternion, [width / 2, 0, width, height]);
    }, [ canvasRef.current, props.input, quaternion ]);

    let maxDimension = Math.max(props.input.width, props.input.height);

    return (
    <div>
        <canvas 
            width={maxDimension * 2}
            height={maxDimension}
            style={{
                width: Math.max(maxDimension * 2, 400),
                height: Math.max(maxDimension, 200),
                border: "1px solid black"
            }}
            contentEditable={true}
            onKeyDown={() => { setIsOffAxisRotation(true) }}
            onKeyUp={() => { setIsOffAxisRotation(false) }}
            onMouseDown={(e) => {
                setIsDragging(true);
                setIsOffAxisRotation(false)
                let { clientX, clientY } = e;
                setLastMousePosition([clientX, clientY]);
            }}
            onMouseMove={(e) => {
                if (!isDragging) {
                    return;
                }
                e.preventDefault();
                let { clientX, clientY } = e;
                let dx = lastMousePosition[0] - clientX;
                let dy = lastMousePosition[1] - clientY;
                setLastMousePosition([clientX, clientY]);
                let m = Math.sqrt(dx * dx + dy * dy);

                let delta: Quaternion;
                if (!isOffAxisRotation) {
                    delta = eToQ([ dy / m, -dx / m, 0, m / 100]);
                } else {
                    delta = eToQ([ 0, 0, 1, dx / 100]);
                }

                let newQuaternion = qMultiply(quaternion, delta);
                let [ x ] = qRotate([1, 0, 0], newQuaternion);
                if (x > 0) {
                    setQuaternion(newQuaternion);
                }
            }}
            onMouseUp={() => setIsDragging(false) }
            onMouseLeave={() => setIsDragging(false)}
            onMouseOut={() => setIsDragging(false)}
            ref={canvasRef} >
        </canvas>
        <br/>
        <button onClick={() => setIsSkateboard((isSkateboard > .5) ? 0 : 1)}>
            { ((isSkateboard > .5) ? "This is a skateboard" : "This is not a skateboard") + " (click to invert)" }
        </button>
        <button onClick={() => props.onOutputConfirmed({
            rotation: quaternion,
            isSkateboard: (isSkateboard > .5) ? 1 : 0
        })} style={{
            color: (isSkateboard > .5) ? "green" : "red"
        }}> Confirm Clip </button>
    </div>) 
}