import * as React from 'react';

export interface CollectionStatistics {
    bounds: number
    frames: number
    rotations: number
}

export function CollectionStats({ frames, bounds, rotations }: CollectionStatistics) {
    const titleStyle = {
        margin: "3px",
        fontSize: "8px",
    }
    const numberStyle = {
        margin: "4px",
        fontSize: "20px",
    }
    return <div style={{
        display: "grid",
        gridTemplateColumns: "1fr 1fr 1fr",
        gridTemplateRows: "1fr 3fr",
        textAlign: "center",
        verticalAlign: "middle",
        height: 40,
    }} >
        <div style={titleStyle}>Frames</div>
        <div style={titleStyle}>Bounds</div>
        <div style={titleStyle}>Rotations</div>
        <div style={numberStyle}>{frames}</div>
        <div style={numberStyle}>{bounds}</div>
        <div style={numberStyle}>{rotations}</div>
    </div>
}