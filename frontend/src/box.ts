
export class Box {

    private grabbed: number[];
    private coordinates: number[];

    constructor(private border: number) {
        // these coordinates represent two points A and B. [x, y]_A and [x, y]_B
        this.coordinates = [0, 0, 0, 0];
        this.grabbed = [];
    }

    public describe() {
        let [ x0, y0, x1, y1 ] = this.coordinates;
        return {
            top: Math.min(y0, y1), 
            bottom: Math.max(y0, y1), 
            left: Math.min(x0, x1), 
            right: Math.max(x0, x1),
        }
    }

    public isEmpty() {
        let { top, bottom, left, right } = this.describe();
        return right == left || top == bottom;
    }

    /**
     * @param param0 a point [x, y] in the box's coordinate space
     */
    public grab([x, y]: [number, number]): boolean {
        let { top, bottom, left, right } = this.describe();
        let isInsideX = left - this.border <= x && x <= right + this.border;
        let isInsideY = top - this.border <= y && y <= bottom + this.border;
        let [ xa, ya, xb, yb ] = this.coordinates;
        this.grabbed = [];
        if (Math.abs(x - xa) < this.border && isInsideY) {
            this.grabbed.push(0);
        }
        if (Math.abs(y - ya) < this.border && isInsideX) {
            this.grabbed.push(1);
        }
        if (Math.abs(x - xb) < this.border && isInsideY) {
            this.grabbed.push(2);
        }
        if (Math.abs(y - yb) < this.border && isInsideX) {
            this.grabbed.push(3);
        }
        return isInsideX && isInsideY;
    }

    /**
     * @param param0 a point [x, y] in the box's coordinate space
     */
    public drag([x, y]: [number, number], indices: number[]) { 
        if (indices.indexOf(0) != -1) {
            this.coordinates[0] = x;
        }
        if (indices.indexOf(1) != -1) {
            this.coordinates[1] = y;
        }
        if (indices.indexOf(2) != -1) {
            this.coordinates[2] = x;
        }
        if (indices.indexOf(3) != -1) {
            this.coordinates[3] = y;
        }
    }

}

class BoxRenderer {

}