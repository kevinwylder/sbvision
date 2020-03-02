import * as React from 'react';

import './list.css';

import { ListRow } from './ListRow'
import { Video, getVideos } from '../../api';
import { ListPagenation } from './ListPagenation';

export function Listing() {

    const PAGE_SIZE = 7;
    // fields for video select
    let [ offset, setOffset ] = React.useState(0);
    let [ maxVideos, setMaxVideos ] = React.useState(0);
    let [ videos, setVideos ] = React.useState<Video[]>([]);

    React.useEffect(() => {
        getVideos(offset, PAGE_SIZE)
        .then(({ total, videos }) => {
            setVideos(videos); 
            setMaxVideos(total);
        })
        .catch(err => {
            console.log(err)
        });
    }, [offset]);

    return (
    <div className="listing">
        {videos.map((v, i) => 
            <ListRow key={i} video={v} />
        )}
        <ListPagenation
            pageSize={PAGE_SIZE}
            maxTotal={maxVideos}
            start={offset}
            end={offset + videos.length}
            onPageSelected={setOffset}
        />
    </div>)


}