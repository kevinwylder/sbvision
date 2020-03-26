
import { Graphable111 } from '.';

export function bezier(b0: number, b1: number, b2: number, b3: number): Graphable111 {
    return (t) => {
        let ti = (1 - t);
        return [
            ti * ti * ti * b0 + 3 * ti * ti * t * b1 + 3 * ti * t * t * b2 + t * t * t * b3,
            3 * ti * ti * (b1 - b0) + 6 * t * ti * (b2 - b1) + 3 * t * t * (b3 - b2)
        ];
    }
}

export const rationalEase: Graphable111 = (x) => {
    return [ (1 / (x + 1) + x - 1), (-1 / ((x + 1) * (x + 1)) + 1) ];
}
