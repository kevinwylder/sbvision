import * as React from "react";
import * as ReactDOM from "react-dom";

import { Listing } from './components/list/Listing';
import { VideoDisplay } from './components/video/VideoDisplay';
import { Label } from "./model";

export function PanelLayout() {

    let [ videoId, setVideoId ] = React.useState<number|undefined>();
    let [ labelId, setLabelId ] = React.useState<number>(0);
    let [ lastClip, setLastClip ] = React.useState<Label>();

    return <div className="main-layout">
        {/* <Listing 
            callback={() => {}}
            /> */}
        <VideoDisplay
            topBarSize={40}
            video={{
                clips: 0,
                duration: 90,
                format: "video/mp4",
                fps: 60,
                id: 2,
                thumbnail: "/images/thumbnail-2.jpg",
                title: "This video doesn't actually exist!!",
                type: 3,
            }} />
    </div>
}

ReactDOM.render(
    <PanelLayout />,
    document.getElementById("root")
)