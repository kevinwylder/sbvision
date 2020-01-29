import * as React from 'react';
import { Label } from '../model';
import { renderSkateboard } from '../skateboard';

interface APIClipResponse {
    labels: number[]
}

export function ClipViewer() {
    let [ clips, setClips ] = React.useState<APIClipResponse>();

    React.useEffect(() => {
        fetch("/skateboards")
        .then(res => res.json())
        .then(clips => setClips(clips));
    }, []);

    return (<div className="clip-container">
        {clips?.labels.sort().map(labelId => 
            <ClipCell 
                key={labelId}
                size={100}
                id={labelId} 
            />
        )}
    </div>)
}

interface ClipCellProps {
    size: number
    id: number
}
function ClipCell(props: ClipCellProps) {

    let [ label, setLabel ] = React.useState<Label>()
    React.useEffect(() => {
        fetch(`/skateboards?id=${props.id}`)
        .then(res => res.json())
        .then(label => setLabel(label));
    }, []);

    let [ deleted, setDeleted ] = React.useState(false);

    let canvas = React.createRef<HTMLCanvasElement>()
    React.useEffect(() => {
        let ctx = canvas.current?.getContext("2d");
        let width = canvas.current?.width;
        let height = canvas.current?.height;
        if (!ctx || !width || !height) {
            return
        }
        ctx.fillStyle = "white";
        ctx.fillRect(0, 0, width, height)
        ctx.textAlign = "center";
        ctx.fillStyle = 'black';
        if (!label) {
            ctx.fillText("Loading...", width / 2, height / 2);
        } else if (deleted) {
            ctx.fillText("Deleted", width / 2, height / 2);
        } else if (label.output.isSkateboard) {
            renderSkateboard(ctx, label.output.rotation, [width * .2, height * .2, width * .8, height * .8]);
        } else {
            ctx.fillText("Not a Skateboard", width / 2, height / 2);
        }
    }, [label, canvas.current])

    return <div className="clip">
        <div>
            <h3>Clip {props.id}</h3>
            <button onClick={() => {
                fetch(`/skateboards?id=${props.id}`, {
                    method: "DELETE"
                })
                .then(res => res.json())
                .then(_ => setDeleted(true));
            }}>Delete this clip</button>
        </div>
        <img src={label?.input.data} 
            style={{
                objectFit: "contain",
                width: props.size * 2,
                height: props.size,
            }}
        />
        <canvas 
            ref={canvas} 
            width={props.size * 2} 
            height={props.size} 
            style={{
                width: props.size * 2,
                height: props.size,
            }}
        />
    </div>
}