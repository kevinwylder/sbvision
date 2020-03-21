
/**
 * This class manages the box being drawn over the video
 */
export class Box {

    private helpTextInterval = 0;
    private empty: boolean = true;
    private dragDistance: number = 0;
    private lastPosition: number[] = [0, 0];
    private grabbed: number[] = [];
    private lastCoordinates = [0, 0, 0, 0];
    private coordinates = [0, 0, 0, 0];
    private dragging: boolean = false;

    constructor(
        private border: number, 
        private areaWidth: number,
        private areaHeight: number,
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

    public bounds() {
        let [ xa, ya, xb, yb ] = this.coordinates;
        return {
            x: Math.min(xa, xb),
            y: Math.min(ya, yb),
            width: Math.abs(xa - xb),
            height: Math.abs(ya - yb)
        }
    }

    public reset() {
        this.empty = true;
        this.coordinates = [0, 0, 0, 0];
        window.clearInterval(this.helpTextInterval);
    }

    public draw(ctx: CanvasRenderingContext2D) {
        ctx.fillStyle = "rgba(0, 0, 0, .4)";
        if (this.empty) {
            ctx.fillRect(0, 0, this.areaWidth, this.areaHeight);
            ctx.font = (this.areaHeight * .05) + "px Arial";
            ctx.fillStyle = "white";
            ctx.fillText("box the skateboard", 50, this.areaHeight * .85);
            return;
        }
        let { top, bottom, left, right } = this.describe();
        ctx.fillRect(0, 0, left, bottom);
        ctx.fillRect(left, 0, this.areaWidth, top);
        ctx.fillRect(right, top, this.areaWidth, this.areaHeight);
        ctx.fillRect(0, bottom, right, this.areaHeight);
        ctx.lineCap = "round";
        ctx.lineWidth = this.border / 2;
        ctx.strokeStyle = "#33b5e5";
        ctx.beginPath();
        ctx.moveTo(left, top);
        ctx.lineTo(left, bottom);
        ctx.lineTo(right, bottom);
        ctx.lineTo(right, top);
        ctx.closePath();
        ctx.stroke();
    }

    public drawHelpDelayed(ctx: CanvasRenderingContext2D) {
        let text = "click to add";
        this.helpTextInterval = window.setTimeout(() => {
            ctx.font = "12px Arial";
            let { width, actualBoundingBoxDescent, actualBoundingBoxAscent } = ctx.measureText(text);
            let height = actualBoundingBoxAscent + actualBoundingBoxDescent;
            let { left, top, bottom, right } = this.describe();
            let fontSize = 12 * Math.min((bottom - top - 10) / height, (right - left - 10) / width);
            ctx.font = fontSize + "px Arial";
            ctx.strokeStyle = "black"
            ctx.lineWidth = 5;
            ctx.strokeText(text, left, (top + bottom + fontSize) / 2);
            ctx.fillStyle = "white"
            ctx.fillText(text, left, (top + bottom + fontSize) / 2);
        }, 1500);
    }

    /**
     * @param param0 a point [x, y] in the box's coordinate space
     */
    public grab([x, y]: [number, number]): boolean {
        this.dragging = true;
        this.lastCoordinates = this.coordinates;
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
            return false;
        }
        this.empty = false;

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

        return true;
    }

    // release calls onSubmit if the stroke was an 
    public release(): [boolean, boolean, boolean] {
        this.dragging = false;

        const wasTap = this.dragDistance < 5;
        let wasInside = false;

        let wasEmpty = this.empty;
        if (wasTap) {
            this.coordinates = this.lastCoordinates;
            this.lastCoordinates = this.coordinates;

            const [ x, y ] = this.lastPosition;
            const { top, bottom, left, right } = this.describe();
            wasInside = left - this.border < x && x < right + this.border && top - this.border < y && y < bottom + this.border;
        }

        return [ wasTap, wasInside, wasEmpty ];
    }

}