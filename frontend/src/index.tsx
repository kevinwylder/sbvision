import * as React from "react";
import * as ReactDOM from "react-dom";

import { VideoController } from './components/VideoController'
import { Video, VideoListing } from './components/VideoList';

export function PanelLayout() {

    let [ videoId, setVideoId ] = React.useState<string>()
    let [ videos, setVideos ] = React.useState<Video[]>();

    React.useEffect(() => {
        fetch("/videos")
        .then(res => res.json())
        .then(res => setVideos(res.videos));
    }, [])

    if (videoId) {
        return <VideoController 
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
        />
    } else {
        return <VideoListing 
                callback={setVideoId}
                videos={videos ? videos : []}
            />
    }
}

ReactDOM.render(
    <PanelLayout />,
    document.getElementById("root")
)