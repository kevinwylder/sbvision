import * as React from 'react';
import './rotation.css';
import { Frame, getRotations } from '../../api';
import { RotationInstructions } from './RotationInstructions';
import { RotationCanvas } from './RotationCanvas';


export function RotationMatcher() {

    let [ showingInstructions, setShowingInstructions ] = React.useState(true);
    let [ frames, setFrames ] = React.useState<Frame[]>();
    let [ page, setPage ] = React.useState(0);
    let [ index, setIndex ] = React.useState(0);
    React.useEffect(() => {
        getRotations()
        .then(rotations => rotations.length && setFrames(rotations))
        .then(_ => setIndex(0));
    }, [page])

    const next = () => {
            if (frames && (index + 1) < frames?.length) {
                setIndex(index+1);
            } else {
                setPage(page+1);
            }
    }

    if (!frames || showingInstructions) {
        return <RotationInstructions 
            isTouchEvent={false} 
            onClosed={() => setShowingInstructions(false)} />
    }

    return <RotationCanvas 
        onFrameComplete={next}
        frame={frames[index]}
        style={{
            width: "100%",
            height: "100%",
        }}
    />
}