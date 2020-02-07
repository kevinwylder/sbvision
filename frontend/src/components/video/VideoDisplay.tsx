import * as React from 'react';

import { Video } from '../../api';

interface VideoDisplayProps {
    video: Video
}

export function VideoDisplay({video}: VideoDisplayProps) {
    return <video
        controls
    >
        <source
            src={`/video?type=${video.type}&id=${video.id}`}
            type={video.format}
        ></source>
    </video> 
}