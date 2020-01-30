import * as React from "react";
import * as ReactDOM from "react-dom";

import { VideoController } from './components/VideoController'
import { VideoListing } from './components/VideoList';
import { ClipViewer, ClipCell } from "./components/ClipViewer";
import { Label } from "./model";

export function PanelLayout() {

    let [ videoId, setVideoId ] = React.useState<string|undefined>();
    let [ labelId, setLabelId ] = React.useState<number>(0);
    let [ lastClip, setLastClip ] = React.useState<Label>();

    return <div className="main-layout">
        <VideoListing 
            callback={(videoId) => {
                setVideoId(undefined);
                setTimeout(() => setVideoId(videoId));
            }}
        />
        {(videoId) ? <VideoController 
            width={700}
            url={`/video?id=${videoId}`}
            callback={(data: Label) => {
                setLastClip(undefined);
                fetch("/skateboards", {
                    headers: {
                        "Content-Type": "application/json"
                    },
                    method: "POST",
                    body: JSON.stringify(data)
                })
                .then(res => res.json())
                .then(({id}) => {
                    setLastClip(data);
                    setLabelId(id);
                });
                
            }}
        /> : <ClipViewer/>}
        {lastClip && <ClipCell size={100} label={{...lastClip, id: labelId}}/>}
    </div>
}

ReactDOM.render(
    <PanelLayout />,
    document.getElementById("root")
)