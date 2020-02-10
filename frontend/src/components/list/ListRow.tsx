import * as React from 'react';

import { Video } from '../../api';

import { Redirect } from 'react-router-dom';

interface VideoRowProps {
    video: Video
    selected: boolean
    onSelect: (video: Video) => void
}

function timeString(duration: number): string {
    if (duration > 60) {
        return `${Math.floor(duration / 60)}:${duration % 60 < 10 ? "0" : ""}${duration % 60}`
    }
    return duration + "s";
}

export function ListRow({video, selected, onSelect}: VideoRowProps) {
    if (selected) {
        return <Redirect push to={`/video/${video.id}`}/>
    }
    return (
        <div className="list-row" onClick={() => onSelect(video)}>
            <img className="list-row-image"
                src={`/image/${video.thumbnail}`} />
            <div className="list-row-text">
                <h3 className="list-row-title" style={{
                    color: selected ? "red" : "black"
                }}> {video.title} </h3>
                <div className="list-row-stats">
                    <div> { timeString(video.duration) } </div> 
                    <div> Analyzed {video.clips} frames </div>
                    <div> Youtube </div>
                </div>
            </div>
        </div>
    )
}