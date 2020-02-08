import * as React from "react";
import * as ReactDOM from "react-dom";

import { Listing } from './components/list/Listing';
import { VideoDisplay } from './components/video/VideoDisplay';
import { Video } from './api';

export function PanelLayout() {

    let [ video, setVideo ] = React.useState<Video>();

    return <div className="main-layout">
        { video ? 
            <VideoDisplay video={video} /> :
            <Listing onVideoSelected={(v) => { setVideo(v); }} /> 
        }
    </div>
}

ReactDOM.render(
    <PanelLayout />,
    document.getElementById("root")
)