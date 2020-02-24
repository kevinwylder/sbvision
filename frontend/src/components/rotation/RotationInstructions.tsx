import * as React from 'react';

import './rotation.css';

interface RotationInstructionProps {
    isTouchEvent: boolean;
    onClosed?: () => void
}

export function RotationInstructions({onClosed, isTouchEvent}: RotationInstructionProps) {
    let move = (isTouchEvent) ? "Drag" : "Move";
    let scroll = (isTouchEvent) ? "Two Finger Drag" : "Scroll Up/Down"
    let accept = (isTouchEvent) ? "Tap Next Button" : "Click"
    let reject = "Delete Key"
    return <div className="rotation-instructions">
        <div></div><div></div><div></div><div onClick={onClosed}>Start</div>
        <div></div><div style={{gridColumn: "2 / span 2"}}>Rotation Controls</div><div></div>
        <div></div><div>{move} Up/Down</div><div>"Kickflip" Axis of rotation</div><div></div>
        <div></div><div>{move} Left/Right</div><div>"Shuv-It" Axis of rotation</div><div></div>
        <div></div><div>{scroll}</div><div>"Impossible" Axis of rotation</div><div></div>
        <div></div><div>{accept}</div><div>Submit rotation</div><div></div>
        <div></div><div>{reject}</div><div>Mark rotation as "bad"</div><div></div>
    </div>
}