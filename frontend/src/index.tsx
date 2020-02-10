import * as React from "react";
import * as ReactDOM from "react-dom";
import {
    BrowserRouter as Router,
    Switch,
    Route,
} from 'react-router-dom';

import { Listing } from './components/list/Listing';
import { VideoDisplay } from './components/video/VideoDisplay';
import { Video } from './api';

export function App() {

    let [ video, setVideo ] = React.useState<Video>();

    return <Router>
        <Switch>
            <Route path="/video/:id">
                <VideoDisplay video={video} />
            </Route>
            <Route exact path="/">
                <Listing onVideoSelected={(v) => { 
                    setVideo(v);
                }} /> 
            </Route>
        </Switch>
    </Router>
}

ReactDOM.render(
    <App />,
    document.getElementById("root")
)