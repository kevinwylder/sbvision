import * as React from 'react';

export interface Video {
    id: string
    title: string
}

interface VideoList {
    videos: Video[]
    callback: (id: string) => void
}

export function VideoListing(props: VideoList) {

    let inputURL = React.createRef<HTMLInputElement>()
    let [selectedTitle, setSelectedTitle] = React.useState<string>()

    return (<div className="title-container">
        <div className="title-cell">
        <input type="url" ref={inputURL} placeholder="youtube url"></input>
        <button onClick={event => {
            fetch("/video", {
                method: "POST",
                body: JSON.stringify({
                    link: inputURL.current?.value
                }),
                headers: {
                    'Content-Type': 'application/json'
                }
            })
            .then(() => {
                location.reload();
            })
        }}>Download video</button>
        </div>
        {props.videos.map(({id, title}, i) => 
            <div 
                key={i}
                className="title-cell"
                onClick={() => {
                    props.callback(id);
                    setSelectedTitle(title);
                }}
                style={{
                    color: (title == selectedTitle) ? "red" : "black"
                }}>
                    {title}
            </div>
        )}
    </div>)


}