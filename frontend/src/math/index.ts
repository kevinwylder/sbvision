
export type Vector = [ number, number, number ];

export type Quaternion = [ number, number, number, number ];

export type EulerAngle = [ number, number, number, number ];

export type Vec3 = [ number, number, number ];

// Graphable is a function of x and y which returns Z and a norm
export type Graphable213 = (x: number, y: number) => [number, Vec3];

// Graphable111 is a function of x that returns y and y'
export type Graphable111 = (x: number) => [ number, number ];

// Graphable14 is a parameterized function of quaternions
export type Graphable14 = (t: number) => Quaternion;

export * from './functions';
export * from './quaternions';