import { Rotation } from "./frame";
import { Quaternion } from "../math";

export class VisualizationManager {

    private ws: WebSocket;

    constructor(onImage: (elem: HTMLImageElement, rot: Quaternion) => void) {
        this.ws = new WebSocket(`ws${window.location.protocol == "https:" ? "s" : ""}://${window.location.host}/app/visualization`);
        this.ws.onmessage = (ev: MessageEvent) => {
            let img = new Image();
            let { s, r } = JSON.parse(ev.data)
            img.src = s;
            img.onload = () => {
                onImage(img, r);
            }
        }
        this.ws.onerror = () => {
            window.location.reload();
        }
        this.ws.onopen = () => {
            this.setRotation([1, 0, 0, 0]);
        }
    }

    public setRotation(r: Quaternion)  {
        if (this.ws.readyState == this.ws.OPEN) {
            this.ws.send(JSON.stringify({
                r: r[0],
                i: r[1],
                j: r[2],
                k: r[3]
            }))
        }
    }

    public disconnect() {
        this.ws.close();
    }

}