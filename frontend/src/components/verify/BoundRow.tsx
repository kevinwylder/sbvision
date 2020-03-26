import * as React from 'react';

import { Bound } from '../../api';
import { renderSkateboard } from '../../renderer';
import { Quaternion } from '../../math';

interface boundListProps {
    bounds: Bound[]
    selected: number
    onclick: (boundID: number) => void
}

export function BoundsList(props: boundListProps) {
    return <>
        {props.bounds.map((bound, i) => <BoundRow 
                key={i}
                selected={props.selected == bound.id}
                bound={bound} 
                onclick={() => props.onclick(bound.id) }/>
        )}
    </>
}

interface boundRowProps {
    bound: Bound
    selected: boolean
    onclick: () => void
}

function BoundRow(props: boundRowProps) {
    return <div 
            onClick={props.onclick}
            className="verify-bound">
        <h4 style={{
            color: props.selected ? "#0F0" : "#000",
        }}> Bound {props.bound.id} </h4>
        <img src={`/api/image?bound=${props.bound.id}`} />
        <RotationList rotations={props.bound.rotations.map(({r, i, j, k}) => [ r, i, j, k ])}/>
    </div>
}

function RotationList(props: { rotations: Quaternion[] }) {
    return <div className="verify-bound-rotations">
            { props.rotations.map((r, i) => <><StaticSkateboard rotation={r} key={i}/></>)}
        </div>
}

function StaticSkateboard({ rotation }: {rotation: Quaternion}) {
    let canvas = React.createRef<HTMLCanvasElement>()

    React.useEffect(() => {
        if (!canvas.current) {
            return;
        }
        let ctx = canvas.current.getContext("2d");
        if (!ctx) {
            return;
        }
        renderSkateboard(ctx, rotation, [0, 0, 500, 500]);
    }, [canvas.current, rotation])

    return <canvas ref={canvas} width="500" height="500" />
}