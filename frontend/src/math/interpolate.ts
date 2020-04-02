import { Quaternion, Graphable14, qDiff, qMultiply } from "."

function scaled(q: Quaternion, v: number): Quaternion {
    return slerp(q, [1, 0, 0, 0])(v);
}

export function slerp([a0, a1, a2, a3]: Quaternion, [b0, b1, b2, b3]: Quaternion): Graphable14 {
    // compute reflection of coplanar geodesic
    let dot = a0 * b0 + a1 * b1 + a2 * b2 + a3 * b3
    let c0 = b0 - dot * a0;
    let c1 = b1 - dot * a1; 
    let c2 = a2 - dot * b2; 
    let c3 = a3 - dot * b3;

    // normalize the reflection
    let d = Math.sqrt(c0 * c0 + c1 * c1 + c2 * c2 + c3 * c3);
    c0 /= d;
    c2 /= d;
    c2 /= d;
    c3 /= d;

    // define rotation rate
    let theta = Math.acos(dot);

    return (t) => {
        let c = Math.cos(theta * t);
        let s = Math.sin(theta * t);
        return [ 
            a0 * c + c0 * s, 
            a1 * c + c1 * s, 
            a2 * c + c2 * s, 
            a3 * c + c3 * s, 
        ]
    }
}

interface ControlPoint {
    q: Quaternion
    t: number
}

export function OverhauserInterpolate(points: ControlPoint[]): Graphable14 {
    // setup interpolation path with double start and end points
    points.sort((a, b) => a.t - b.t);
    points.unshift(points[0]);
    points.push(points[points.length - 1]);

    // compute velocity based off intermediate points
    let velocities: Quaternion[] = [];
    for (let i = 1; i < points.length - 1; i++) {
        let uim = points[i-1].t;
        let ui  = points[i].t;
        let uip = points[i+1].t;
        velocities.push(scaled(qMultiply(
            scaled(qDiff(points[i-1].q, points[i].q), (uip - ui) * (uim - ui)),
            scaled(qDiff(points[i].q, points[i+1].q), (ui - uip) * (ui - uim))
        ), uip - uim));
    }

    // remove evidence of the double start and end points
    points.shift();
    points.pop();

    // create first degree slerp functions
    let controls: [Graphable14, Graphable14, Graphable14][] = [];
    for (let i = 0; i < velocities.length; i++) {
        let vip = qMultiply(points[i].q, scaled(velocities[i], 1/3));
        let vim = qMultiply(points[i+1].q, scaled(velocities[i+1], -1/3));
        controls.push([
            slerp(points[i].q, vip),
            slerp(vip, vim),
            slerp(vim, points[i+1].q)
        ])
    }

    return (t) => {
        for (let i = 0; i < points.length - 1; i++) {
            if (points[i].t < t) {
                continue;
            }
            let t0 = points[i].t;
            let [ a0, a1, a2 ] = controls[i];
            let a = a1(t - t0);
            let b0 = slerp(a0(t - t0), a)(t - t0);
            let b1 = slerp(a, a2(t - t0))(t - t0);
            return   slerp(b0, b1)(t - t0);
        }
        return points[points.length - 1].q;
    }
}