import * as React from "react";
import * as ReactDOM from "react-dom";

import { VideoController } from './components/VideoController'

ReactDOM.render(
    <VideoController
        width={600}
        callback={(image: ImageData) => {
            
        }}
        url={"/video?id=biSqKzgnjV8"}
    ></VideoController>,
    document.getElementById("root")
)