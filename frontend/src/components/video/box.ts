import { Bounds } from '../../api';

export class Box {

    private empty: boolean = true;
    private helpTextDelay: number = 1000;
    private dragStartTime: number = 0;
    private helpTextInterval: number = 0;
    private dragDistance: number = 0;
    private lastPosition: number[] = [0, 0];
    private grabbed: number[] = [];
    private lastCoordinates = [0, 0, 0, 0];
    private coordinates = [0, 0, 0, 0];
    private type = "click/tap";
    private dragging: boolean = false;

    constructor(
        private areaWidth: number,
        private areaHeight: number,
        private border: number, 
        private ctx: CanvasRenderingContext2D
    ) {
        this.reset();
     }

    // describe puts semantic meaning behind coordinates
    private describe() {
        let [ x0, y0, x1, y1 ] = this.coordinates;
        return {
            top: Math.min(y0, y1), 
            bottom: Math.max(y0, y1), 
            left: Math.min(x0, x1), 
            right: Math.max(x0, x1),
        }
    }

    private bounds() {
        let [ xa, ya, xb, yb ] = this.coordinates;
        return {
            left: Math.min(xa, xb),
            top: Math.min(ya, yb),
            width: Math.abs(xa - xb),
            height: Math.abs(ya - yb)
        }
    }

    private reset() {
        this.empty = true;
        this.coordinates = [0, 0, 0, 0];
        this.draw();
    }

    private draw() {
        this.ctx.clearRect(0,  0, this.areaWidth, this.areaHeight);
        this.ctx.fillStyle = "rgba(0, 0, 0, .7)";
        if (this.empty) {
            this.ctx.fillRect(0, 0, this.areaWidth, this.areaHeight);
            this.ctx.font = "12px Arial";
            this.ctx.fillStyle = "white";
            this.ctx.fillText("Draw a box around the skateboard", 100, 50);
            return;
        }
        let { top, bottom, left, right } = this.describe();
        this.ctx.fillRect(0, 0, left, bottom);
        this.ctx.fillRect(left, 0, this.areaWidth, top);
        this.ctx.fillRect(right, top, this.areaWidth, this.areaHeight);
        this.ctx.fillRect(0, bottom, right, this.areaHeight);
        this.ctx.lineCap = "round";
        this.ctx.lineWidth = this.border / 2;
        this.ctx.strokeStyle = "#33b5e5";
        this.ctx.beginPath();
        this.ctx.moveTo(left, top);
        this.ctx.lineTo(left, bottom);
        this.ctx.lineTo(right, bottom);
        this.ctx.lineTo(right, top);
        this.ctx.closePath();
        this.ctx.stroke();
    }

    /**
     * @param param0 a point [x, y] in the box's coordinate space
     */
    public grab([x, y]: [number, number], type: string): boolean {
        if (this.helpTextInterval) {
            window.clearInterval(this.helpTextInterval);
        }
        this.empty = false;
        this.type = type;
        this.dragging = true;
        this.lastCoordinates = this.coordinates;
        this.dragStartTime = new Date().getTime();
        this.dragDistance = 0
        this.lastPosition = [x, y];

        let { top, bottom, left, right } = this.describe();
        let isInsideX = left - this.border <= x && x <= right + this.border;
        let isInsideY = top - this.border <= y && y <= bottom + this.border;
        let [ xa, ya, xb, yb ] = this.coordinates;
        this.grabbed = [];
        if (isInsideY) {
            let da = Math.abs(x - xa);
            let db = Math.abs(x - xb);
            if (da < Math.min(this.border, db)) {
                this.grabbed.push(0);
            } else if (db < this.border) {
                this.grabbed.push(2);
            }

        }
        if (isInsideX) {
            let da = Math.abs(y - ya);
            let db = Math.abs(y - yb);
            if (da < Math.min(this.border, db)) {
                this.grabbed.push(1);
            } else if (db < this.border) {
                this.grabbed.push(3);
            }
        }
        if (this.grabbed.length == 0) {
            this.coordinates = [x, y, x, y];
            this.grabbed = [0, 1];
        }
        return isInsideX && isInsideY;
    }

    /**
     * @param param0 a point [x, y] in the box's coordinate space
     */
    public drag([x, y]: [number, number]) { 
        if (!this.dragging) {
            return;
        }

        x = Math.max(0, Math.min(x, this.areaWidth));
        y = Math.max(0, Math.min(y, this.areaHeight));

        let [ lastX, lastY ] = this.lastPosition;
        this.dragDistance += Math.sqrt((x - lastX) * (x - lastX) + (y - lastY) * (y - lastY))
        this.lastPosition = [x, y];

        if (this.grabbed.indexOf(0) != -1) {
            this.coordinates[0] = x;
        }
        if (this.grabbed.indexOf(1) != -1) {
            this.coordinates[1] = y;
        }
        if (this.grabbed.indexOf(2) != -1) {
            this.coordinates[2] = x;
        }
        if (this.grabbed.indexOf(3) != -1) {
            this.coordinates[3] = y;
        }
        this.draw();
    }

    // release calls onSubmit if the stroke was an 
    public release(onSubmit: (b: Bounds) => void, onRefuse: () => void) {
        this.dragging = false;

        const wasTap = this.dragDistance < 5 && (new Date().getTime() - this.dragStartTime) < 200;

        if (wasTap) {
            this.coordinates = this.lastCoordinates;
            const [ x, y ] = this.lastPosition;
            const { top, bottom, left, right } = this.describe();
            const wasInside = left - this.border < x && x < right + this.border 
                            && top - this.border < y && y < bottom + this.border;
            if (wasInside) {
                onSubmit(this.bounds());
            } else {
                onRefuse();
            }
            this.reset();
        } else {
            let text = this.type + " to submit";
            this.helpTextInterval = window.setTimeout(() => {
                this.ctx.font = "12px Arial";
                let { width, actualBoundingBoxDescent, actualBoundingBoxAscent } = this.ctx.measureText(text);
                let height = actualBoundingBoxAscent + actualBoundingBoxDescent;
                let { left, top, bottom, right } = this.describe();
                let fontSize = 12 * Math.min((bottom - top - 10) / height, (right - left - 10) / width);
                this.ctx.font = fontSize + "px Arial";
                this.ctx.strokeStyle = "black"
                this.ctx.lineWidth = 5;
                this.ctx.strokeText(text, left, (top + bottom + fontSize) / 2);
                this.ctx.fillStyle = "white"
                this.ctx.fillText(text, left, (top + bottom + fontSize) / 2);
                this.helpTextDelay = Math.min(this.helpTextDelay + 300, 5000);
            }, this.helpTextDelay);
        }
    }

}