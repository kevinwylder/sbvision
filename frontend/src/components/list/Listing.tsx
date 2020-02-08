import * as React from 'react';

import '../../styles/list.css';

import { ListRow } from './ListRow'
import { Video, getVideos } from '../../api';
import { ListSidebar } from './ListSidebar';
import { ListPagenation } from './ListPagenation';

interface VideoList {
    onVideoSelected: (video: Video) => void
}

export function Listing(props: VideoList) {

    const PAGE_SIZE = 7;
    // fields for video select
    let [ offset, setOffset ] = React.useState(0);
    let [ maxVideos, setMaxVideos ] = React.useState(0);
    let [ selectedVideo, setSelectedVideo ] = React.useState<Video>()
    let [ videos, setVideos ] = React.useState<Video[]>([]);

    React.useEffect(() => {
        getVideos(offset, PAGE_SIZE)
        .then(({ total, videos}) => {
            setVideos(videos); 
            setMaxVideos(total);
        })
        .catch(err => {
            console.log(err)
        });
    }, [offset])

    // Pass selected video up the tree
    React.useEffect(() => selectedVideo && props.onVideoSelected(selectedVideo), [selectedVideo]);

    return (
    <div className="listing">
        <ListSidebar 
            onVideoAdded={setSelectedVideo}
        />
        <div className="list-container">
            {videos.map((v, i) => 
                <ListRow
                    selected={selectedVideo ? v.id == selectedVideo.id : false}
                    onSelect={setSelectedVideo}
                    key={i}
                    video={v}
                />
            )}
            <ListPagenation
                pageSize={PAGE_SIZE}
                maxTotal={maxVideos}
                start={offset}
                end={offset + videos.length}
                onPageSelected={setOffset}
            />
        </div>
    </div>)


}