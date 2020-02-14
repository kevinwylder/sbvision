import * as React from "react";
import * as ReactDOM from "react-dom";
import {
    BrowserRouter as Router,
    Switch,
    Route,
} from 'react-router-dom';

import { Listing } from './components/list/Listing';
import { VideoDisplay } from './components/video/VideoDisplay';
import { RotationInstructions } from "./components/rotation/RotationInstructions";
import { RotationMatcher } from "./components/rotation/RotationMatcher";

export function App() {

    return <Router>
        <Switch>
            <Route path="/video/:id">
                <VideoDisplay />
            </Route>
            <Route exact path="/">
                <Listing /> 
            </Route>
            <Route path="/rotations">
                <RotationMatcher/>
            </Route>
        </Switch>
    </Router>
}

ReactDOM.render(
    <App />,
    document.getElementById("root")
)