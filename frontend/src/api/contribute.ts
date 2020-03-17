import { session } from './session';
import { Bound, Frame, Rotation } from './frame';

export function uploadFrame(videoId: number, elem: HTMLVideoElement): Promise<Frame> {
    let frame = Math.ceil(elem.currentTime * 1000);
    let canvas = document.createElement("canvas");
    canvas.width = elem.videoWidth;
    canvas.height = elem.videoHeight;
    let ctx = canvas.getContext("2d");
    if (!ctx) {
        return Promise.reject("Cannot acquire canvas context");
    }
    ctx.drawImage(elem, 0, 0);
    // check if the frame is all black
    let okay = false
    for (let i = 0; !okay && i < 10; i++) {
        let pixel = ctx.getImageData(Math.floor(Math.random() * elem.videoWidth), Math.floor(Math.random() * elem.videoHeight), 1, 1);
        if (pixel.data[0] + pixel.data[1] + pixel.data[2] > 20) {
            okay = true;
        }
    }
    if (!okay) {
        return Promise.reject("Frame is too dark to upload");
    }
    return fetch(canvas.toDataURL("image/png")) // convert data to array
    .then(res => res.arrayBuffer())
    .then(data => {
        // convert array to form data body
        let form = new FormData();
        form.append("image", new Blob([new Uint8Array(data)]));
        return form;
    })
    .then(body => fetch(`/app/contribute/frame?video=${videoId}&frame=${frame}`,
    { // send form data to the server
        method: "POST",
        body,
        headers: session,
    }))
    .then(res => res.json())
}

export function addBounds(frameId: number, bounds: {x: number, y: number, width: number, height: number}): Promise<Bound> {
    return fetch(`/app/contribute/bounds?frame=${frameId}`, {
        method: "POST",
        body: JSON.stringify({
            x: Math.round(bounds.x),
            y: Math.round(bounds.y),
            width: Math.round(bounds.width),
            height: Math.round(bounds.height),
        }),
        headers: session,
    })
    .then(res => res.json())
}

export function addRotation(boundId: number, r: Rotation): Promise<Rotation[]> {
    return fetch(`/app/contribute/rotation?bound=${boundId}`,{
        method: "POST",
        headers: session,
        body: JSON.stringify(r),
    })
    .then(res => res.json())
}