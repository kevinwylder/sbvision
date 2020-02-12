import * as React from 'react';

import { Redirect } from 'react-router-dom';

import { CollectionStats } from '../CollectionStats';

interface VideoControlProps {
    bounds: {
        left: number,
        top: number,
        width: number,
        height: number
    }
    stats: {
        frames: number
        bounds: number
        rotations: number
    }
}

export function VideoControls(props: VideoControlProps) {

    let [ goBack, setGoBack ] = React.useState(false);

    if (goBack) {
        return <Redirect to="/" />
    }

    return <div
        className="video-controls"
        style={ props.bounds }>
        <div 
            onClick={() => setGoBack(true)}
            className="video-back">
            Back to Videos
        </div>
        <CollectionStats { ...props.stats } />
    </div>
}