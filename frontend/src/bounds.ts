
export type Box = [number, number, number, number]

export const BOX_BORDER = 12;

interface BoxDescription {
    top: number,
    bottom: number,
    left: number, 
    right: number, 
}

export function describeBox(box: Box): BoxDescription {
    let [ x0, y0, x1, y1 ] = box;
    return {
        top: Math.min(y0, y1), 
        bottom: Math.max(y0, y1), 
        left: Math.min(x0, x1), 
        right: Math.max(x0, x1),
    }
}

export function drawBoundingBox(ctx: CanvasRenderingContext2D, width: number, height: number, box: Box, ) {
    let { left, top, bottom, right } = describeBox(box);
    ctx.fillStyle = "rgba(0, 0, 0, .5)";
    ctx.fillRect(0, 0, left, bottom);
    ctx.fillRect(left, 0, width, top);
    ctx.fillRect(right, top, width, height);
    ctx.fillRect(0, bottom, right, height);
    ctx.fillStyle = "#33b5e5";
    ctx.fillRect(left - BOX_BORDER, top, BOX_BORDER, bottom - top);
    ctx.fillRect(left, top - BOX_BORDER, right - left, BOX_BORDER);
    ctx.fillRect(right, top, BOX_BORDER, bottom - top);
    ctx.fillRect(left, bottom, right - left, BOX_BORDER);
    ctx.fillStyle = "#ace1f4";
    ctx.fillRect(left - BOX_BORDER, top - BOX_BORDER, BOX_BORDER, BOX_BORDER);
    ctx.fillRect(right, top - BOX_BORDER, BOX_BORDER, BOX_BORDER);
    ctx.fillRect(left - BOX_BORDER, bottom, BOX_BORDER, BOX_BORDER);
    ctx.fillRect(right, bottom, BOX_BORDER, BOX_BORDER);
}

export function dragBoxBounds(box: Box, dx: number, dy: number, indexes: number[] ): Box {
    let [ x0, y0, x1, y1 ] = box;
    if (indexes.indexOf(0) != -1) {
        x0 += dx;
    }
    if (indexes.indexOf(1) != -1) { 
        y0 += dy;
    }
    if (indexes.indexOf(2) != -1) {
        x1 += dx;
    }
    if (indexes.indexOf(3) != -1) {
        y1 += dy;
    }
    return [ Math.floor(x0), Math.floor(y0), Math.floor(x1), Math.floor(y1) ];
}

export function chooseDraggingIndexes(box: Box, x: number, y: number): number[] {
    if (!isInsideBox(box, x, y)) {
        return [2, 3];
    }
    let { top, bottom, left, right } = describeBox(box);
    let collisions: number[] = [];

    if (x + BOX_BORDER > left && x < left) {
        collisions.push(left == box[0] ? 0 : 2);
    } 
    if (y + BOX_BORDER > top && y < top) {
        collisions.push(top == box[1] ? 1 : 3);
    }
    if (x > right && x - right < BOX_BORDER) {
        collisions.push(right == box[2] ? 2 : 0);
    }
    if (y > bottom && y - bottom < BOX_BORDER) {
        collisions.push(bottom == box[3] ? 3 : 1);
    }
    return collisions;
}

export function isInsideBox(box: Box, x: number, y: number): boolean { 
    let { top, bottom, left, right } = describeBox(box);
    return x + BOX_BORDER > left && 
           y + BOX_BORDER > top &&
           x - BOX_BORDER < right &&
           y - BOX_BORDER < bottom;

}