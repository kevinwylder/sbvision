
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
    let cx0 = x - box[0] > -BOX_BORDER && x - box[0] < 0;
    let cy0 = y - box[1] > -BOX_BORDER && y - box[1] < 0;
    let cx1 = x - box[2] > 0 && x - box[2] < BOX_BORDER;
    let cy1 = y - box[3] > 0 && y - box[3] < BOX_BORDER;
    let collisions: number[] = [];
    if (cx0) {
        collisions.push(0);
    } 
    if (cy0) {
        collisions.push(1);
    }
    if (cx1) {
        collisions.push(2);
    }
    if (cy1) {
        collisions.push(3);
    }
    return collisions;
}

export function isInsideBox(box: Box, x: number, y: number): boolean { 
    return x > box[0] - BOX_BORDER && 
                 y > box[1] - BOX_BORDER &&
                 x < box[2] + BOX_BORDER &&
                 y < box[3] + BOX_BORDER;
}