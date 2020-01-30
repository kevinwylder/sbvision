import * as React from 'react';

export interface Video {
    id: string
    title: string
    progress: number
}

interface VideoList {
    callback: (id: string|undefined) => void
}

export function VideoListing(props: VideoList) {

    let inputURL = React.createRef<HTMLInputElement>()
    let [ isDownloading, setIsDownloading ] = React.useState(false);

    // fields for video select
    let [ selectedVideo, setSelectedVideo ] = React.useState<Video>()
    let [ videos, setVideos ] = React.useState<Video[]>([]);
    React.useEffect(() => {
        fetch("/videos")
        .then(res => res.json())
        .then(res => setVideos(res.videos));
    }, [])

    return (<div className="title-container">
        <div className="title-cell">
        <input type="url" ref={inputURL} placeholder="youtube url"></input>
        <button 
        disabled={isDownloading}
        onClick={() => {
            if (isDownloading) {
                return;
            }
            let link = inputURL.current?.value
            if (!/https:\/\/www.youtube.com\/watch\?v=[a-zA-Z0-9]+/.test(link || "")) {
                return;
            }
            setIsDownloading(true);
            let ws = new WebSocket(`ws${window.location.protocol === "https" ? "s" : ""}://${window.location.host}/video-download`);
            ws.onopen = () => {
                console.log("sending", link);
                ws.send(JSON.stringify({ link }))
            }
            ws.onmessage = function(this: WebSocket, event: MessageEvent) {
                let video: Video = JSON.parse(event.data);
                if ("error" in video) {
                    setSelectedVideo(undefined);
                    alert(video["error"]);
                } else if (!selectedVideo || selectedVideo.id != video.id) {
                    setVideos([video, ...videos]);
                    setSelectedVideo(video)
                }
            }
            ws.onclose = () => {
                if (selectedVideo) {
                    props.callback(selectedVideo.id);
                }
                setIsDownloading(false);
            }
        }}>Download video</button>
        </div>
        {videos.map(({id, title, progress}, i) => 
            <div key={i} className="title-cell">
                <h3
                    onClick={() => {
                        props.callback(id);
                        setSelectedVideo({id, title, progress});
                    }}
                    style={{
                        color: (selectedVideo && id == selectedVideo.id) ? "red" : "black"
                    }}
                > {title} {(progress != 100) ? ` ${progress}% downloaded` : "" } </h3>
                <button onClick={() => {
                    fetch(`/video?id=${id}`, {
                        method: "DELETE"
                    })
                    .then(_ => setVideos(videos.filter(video => id != video.id)))
                }}>Delete</button>
            </div>
        )}
    </div>)


}