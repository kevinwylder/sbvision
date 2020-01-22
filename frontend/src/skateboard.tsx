import * as React from "react";
import * as ReactDOM from "react-dom";

import { VideoController } from './components/VideoController'

document.domain = "www.youtube.com"

ReactDOM.render(
    <VideoController
        width={600}
        height={450}
        url={"https://www.youtube.com/embed/hODoEs_RNGw?autoplay=1"}
    ></VideoController>,
    document.getElementById("root")
)