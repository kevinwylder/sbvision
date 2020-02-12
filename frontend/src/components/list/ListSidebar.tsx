import * as React from 'react';

import { Video, addVideo } from '../../api';
import { CollectionStats, CollectionStatistics } from '../CollectionStats';
import { Redirect } from 'react-router-dom';

interface ListSidebarProps {
    stats?: CollectionStatistics
}

export function ListSidebar(props: ListSidebarProps) {

    let inputURL = React.createRef<HTMLInputElement>()
    let [ errorText, setErrorText ] = React.useState("");
    let [ video, setVideo ] = React.useState<Video>()

    const onVideoDiscoveryRequest = () => {
        setErrorText("loading...")
        if (!inputURL.current) {
            return;
        }
        if (!inputURL.current.value) {
            setErrorText("Error - Please put a URL in the box")
            return;
        }
        if (!inputURL.current.value.startsWith("https://www.youtube.com")) {
            setErrorText("Error - Please use a full youtube URL")
            return;
        }
        if (!inputURL.current.value.startsWith("https://www.youtube.com/watch?v=")) {
            setErrorText("Error - Please use the \"/watch?v= url\" format")
            return;
        }
        addVideo(inputURL.current.value, 1)
        .then(setVideo)
        .catch(err => setErrorText("Server Error - " + err.toString()));
    }

    if (video) {
        return <Redirect to={`/video/${video.id}`} />
    }

    return (<>
    <div className="list-add-space"/>
    <div className="list-add">
        <div className="list-add-always">
            <h3>Add Video</h3>
            <input 
                type="url" 
                ref={inputURL} 
                placeholder="https://www.youtube.com/watch?v="
                className="list-url-textbox"
            />
            <button onClick={onVideoDiscoveryRequest}> Go</button>
            <div className="list-add-error">
                {errorText}
            </div>
            <h3>The dataset so far</h3>
            <div style={{
                margin: "10px",
            }}>
                { props.stats && <CollectionStats {...props.stats} /> } 
            </div>

        </div>


        <div className="list-add-extra">
            <p>Welcome to the kwylder SkateboardVision project. </p> 
            <p>We are collecting a dataset to learn how to identify skateboard orientations from video frames </p>
        </div>
    </div>
    </>)
}