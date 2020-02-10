import * as React from 'react';

import { Video, addVideo } from '../../api';

interface ListAddProps {
    onVideoAdded: (video: Video) => void
}

export function ListSidebar(props: ListAddProps) {

    let inputURL = React.createRef<HTMLInputElement>()
    let [ errorText, setErrorText ] = React.useState("");

    const onVideoDiscoveryRequest = () => {
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
        .then(video => props.onVideoAdded(video))
        .catch(err => setErrorText("Server Error - " + err.toString()));

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
        </div>


        <div className="list-add-extra">
            <p>Welcome to the kwylder SkateboardVision project. </p> 
            <p>We are collecting a dataset of skateboards for training a neural network to identify skateboards and their orientation. </p>
            <p>Go ahead and paste a youtube video link in the box above to add it to the collection on the right</p>
            <p>Some things that might go on this sidebar in the future include</p>
            <ol>
                <li>User Login</li>
                <li>Clip Review</li>
                <li>Total Dataset Visualization</li>
            </ol>
        </div>
    </div>
    </>)
}