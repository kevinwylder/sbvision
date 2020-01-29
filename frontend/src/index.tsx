import * as React from "react";
import * as ReactDOM from "react-dom";

import { VideoController } from './components/VideoController'
import { Video, VideoListing } from './components/VideoList';
import { ClipViewer } from "./components/ClipViewer";

export function PanelLayout() {

    let [ videoId, setVideoId ] = React.useState<string>()
    let [ videos, setVideos ] = React.useState<Video[]>();

    React.useEffect(() => {
        fetch("/videos")
        .then(res => res.json())
        .then(res => setVideos(res.videos));
    }, [])

    return <div className="main-layout">
        <VideoListing 
            callback={setVideoId}
            videos={videos ? videos : []}
        />
        {(videoId) ? <VideoController 
            width={700}
            url={`/video?id=${videoId}`}
            callback={(data) => {
                fetch("/skateboards", {
                    headers: {
                        "Content-Type": "application/json"
                    },
                    method: "POST",
                    body: JSON.stringify(data)
                })
            }}
        /> : <ClipViewer/>}
        
    </div>
}

ReactDOM.render(
    <PanelLayout />,
    document.getElementById("root")
)