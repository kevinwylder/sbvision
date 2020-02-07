import * as React from "react";
import * as ReactDOM from "react-dom";

import { Listing } from './components/list/Listing';
import { Label } from "./model";

export function PanelLayout() {

    let [ videoId, setVideoId ] = React.useState<number|undefined>();
    let [ labelId, setLabelId ] = React.useState<number>(0);
    let [ lastClip, setLastClip ] = React.useState<Label>();

    return <div className="main-layout">
        <Listing 
            callback={() => {}}
            />
    </div>
}

ReactDOM.render(
    <PanelLayout />,
    document.getElementById("root")
)