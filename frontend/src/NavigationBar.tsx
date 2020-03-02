import * as React from 'react';
import {
    BrowserRouter as Router,
    Switch,
    Route,
    Link,
} from 'react-router-dom';

import './style.css';
import { VideoDisplay } from './components/video/VideoDisplay';
import { Listing } from './components/list/Listing';
import { RotationMatcher } from './components/rotation/RotationMatcher';
import { RotatingSkateboard } from './components/RotatingSkateboard';
import { ApiDocs } from './components/docs/ApiDocs';
import { AboutPage } from './components/AboutPage';

export function NavigationBar() {

    return <Router>
        <div className="nav-bar">
            <Link to="/" className="nav-bar-tab"><RotatingSkateboard/></Link>
            <Link to="/videos" className="nav-bar-tab"><div>Videos</div></Link>
            <Link to="/rotations" className="nav-bar-tab"><div>Rotations</div></Link>
            <Link to="/verifications" className="nav-bar-tab"><div>Verifications</div></Link>
            <Link to="/api-docs" className="nav-bar-tab"><div>API</div></Link>
        </div>
            <Switch>
                <Route exact path="/" >
                    <AboutPage />
                </Route>
                <Route path="/video/:id">
                    <VideoDisplay />
                </Route>
                <Route exact path="/videos">
                    <Listing /> 
                </Route>
                <Route exact path="/rotations">
                    <RotationMatcher/>
                </Route>
                <Route exact path="/api-docs">
                    <ApiDocs/>
                </Route>
            </Switch>
        </Router>
}