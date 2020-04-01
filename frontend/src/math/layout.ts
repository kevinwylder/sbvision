
type Box = [ number, number, number, number ];

export const describeBox = ([x0, x1, x2, x3]: Box) => ({
    left: Math.min(x0, x2),
    top: Math.min(x1, x3),
    width: Math.abs(x0 - x2),
    height: Math.abs(x1 - x3)
})