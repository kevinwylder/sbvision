import * as React from 'react';

import './list.css';

import { ListRow } from './ListRow'
import { Video, getVideos } from '../../api';
import { ListSidebar } from './ListSidebar';
import { ListPagenation } from './ListPagenation';
import { getStats } from '../../api/data';
import { CollectionStatistics } from '../CollectionStats';

export function Listing() {

    const PAGE_SIZE = 7;
    // fields for video select
    let [ offset, setOffset ] = React.useState(0);
    let [ maxVideos, setMaxVideos ] = React.useState(0);
    let [ videos, setVideos ] = React.useState<Video[]>([]);
    let [ stats, setStats ] = React.useState<CollectionStatistics>();

    React.useEffect(() => {
        Promise.all([getVideos(offset, PAGE_SIZE), getStats()])
        .then(([{ total, videos }, stats]) => {
            setVideos(videos); 
            setMaxVideos(total);
            setStats(stats);
        })
        .catch(err => {
            console.log(err)
        });
    }, [offset]);

    return (
    <div className="listing">
        <ListSidebar stats={stats} />
        <div className="list-container">
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
        </div>
    </div>)


}