import * as React from "react";
import * as ReactDOM from "react-dom";

import { VideoController } from './components/VideoController'
import { Video, VideoListing } from './components/VideoList';
import { RotationMatcher }  from './components/RotationMatcher'

export function PanelLayout() {

    let [ videoId, setVideoId ] = React.useState<string>()
    let [ videos, setVideos ] = React.useState<Video[]>();
    let [ clip, setClip ] = React.useState<ImageData>();

    React.useEffect(() => {
        fetch("/videos")
        .then(res => res.json())
        .then(res => setVideos(res.videos));
    }, [])

    return <div>
        <VideoListing 
            callback={setVideoId}
            videos={videos ? videos : []}
        />
        { videoId && <VideoController 
            width={700}
            url={`/video?id=${videoId}`}
            callback={(data) => setClip(data)}
        /> }
        <br/>
        <br/>
        <br/>
        { clip && <RotationMatcher
            input={clip}
            onOutputConfirmed={(output) => {
                if (!clip) return;
                let data = btoa(String.fromCharCode.apply(null, Array.from(clip.data.filter((_, idx) => idx % 4 != 3))));
                let body = JSON.stringify({
                        input: {
                            width: clip?.width,
                            height: clip?.height,
                            data: data
                        },
                        output: output
                    })
                console.log(body);
                fetch("/skateboards", { method: "POST", body })
                .then( res => res.json())
                .then(() => setClip(undefined));
            }}
        /> }
    </div>
}

ReactDOM.render(
    <PanelLayout />,
    document.getElementById("root")
)