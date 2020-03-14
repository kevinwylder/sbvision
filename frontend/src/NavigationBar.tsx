import * as React from 'react';
import {
    BrowserRouter as Router,
    Switch,
    Route,
    Link,
    useParams,
} from 'react-router-dom';

import './style.css';
import { VideoDisplay } from './components/VideoManager';
import { Listing } from './components/list/Listing';
import { RotationMatcher } from './components/verify/RotationMatcher';
import { RotatingSkateboard } from './components/RotatingSkateboard';
import { ApiDocs } from './components/docs/ApiDocs';
import { AboutPage } from './components/AboutPage';
import { DataVisualization } from './components/Visualization';

export function Main() {

    return <Router>
        <Switch>
            <Route exact path="/" >
                <NavigationBar />
                <AboutPage />
            </Route>
            <Route path="/video/:id">
                <NavigationBar />
                <VideoDisplay />
            </Route>
            <Route exact path="/videos">
                <NavigationBar />
                <Listing /> 
            </Route>
            <Route path="/rotation/:id">
                <NavigationBar />
                <RotationMatcher/>
            </Route>
            <Route exact path="/explore">
                <NavigationBar />
                <DataVisualization />
            </Route>
            <Route exact path="/api-docs">
                <NavigationBar />
                <ApiDocs/>
            </Route>
        </Switch>
    </Router>
}

function NavigationBar() {

    // load the renderer and get it's size
    let params: { id?: string } = useParams();
    console.log(params);

    return <div className="nav-bar">
        <Link to={  params.id ? `/rotation/${params.id}` : "/" } className="nav-bar-tab"><RotatingSkateboard/></Link>
        <Link to="/" className="nav-bar-tab"><div>About</div></Link>
        <Link to="/videos" className="nav-bar-tab"><div>Videos</div></Link>
        <Link to="/explore" className="nav-bar-tab"><div>Explore</div></Link>
        <Link to="/api-docs" className="nav-bar-tab"><div>API</div></Link>
    </div>
}