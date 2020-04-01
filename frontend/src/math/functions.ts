
import { Graphable111, Graphable213 } from '.';

export const one: Graphable111 = () => [1, 0];

export const identity: Graphable111 = (x) => [x, 1];

export const absolute: Graphable111 = (x) => [Math.abs(x), Math.sign(x)];

// this function is basically x^-1 + x, but translated to the origin
export const rationalEase: Graphable111 = (x) => {
    return [ (1 / (x + 1) + x - 1), (-1 / ((x + 1) * (x + 1)) + 1) ];
}

export function bezier(b0: number, b1: number, b2: number, b3: number): Graphable111 {
    return (t) => {
        let ti = (1 - t);
        return [
            ti * ti * ti * b0 + 3 * ti * ti * t * b1 + 3 * ti * t * t * b2 + t * t * t * b3,
            3 * ti * ti * (b1 - b0) + 6 * t * ti * (b2 - b1) + 3 * t * t * (b3 - b2)
        ];
    }
}

export function exponentNorm(exponent: number): Graphable111 {
    return (t) => {
        let p = Math.pow(t, exponent);
        return [
            Math.pow( 1 - p, 1 / exponent),
            Math.pow( 1 - p, 1 / exponent - 1) * Math.pow(t, exponent - 1),
        ]
    }
}

export function translate(f: Graphable111, tx: number): Graphable111 {
    return (x) => {
        let [ y, dy ] = f(x + tx);
        return [y, dy]
    }
}

export function scale(base: Graphable111, s: number): Graphable111 {
    return (x) => {
        let [ y, dy ] = base(x);
        return [ y * s, dy * s];
    }
}

export function compose(f: Graphable111, g: Graphable111): Graphable111 {
    return (x) => {
        let [ x1, dx1 ] = g(x);
        let [ x2, dx2 ] = f(x1);
        return [ x2, dx2 * dx1 ];
    }
}

export function cross(fx: Graphable111, fy: Graphable111): Graphable213 {
    return (x, y) => {
        let [ zx, dx ] = fx(x);
        let [ zy, dy ] = fy(y);
        dy *= zx;
        dx *= zy;
        let n = Math.sqrt(1 + dx * dx + dy * dy);
        return [ zx * zy, [ -dx / n, -dy / n, 1.0 / n] ];
    }
}