import { Quaternion } from ".";

// Skateboard Geometry

const tailIncline = Math.PI / 4;

// vertPos (float 3), vertNorm (float 3)
const verts = [
    2, .5, 0,           0, 0, 1,
    2, -.5, 0,          0, 0, 1,
    2.5, 0, .25,        Math.sin(-tailIncline), 0, Math.cos(-tailIncline),
    -2, .5, 0,          0, 0, 1,
    -2, -.5, 0,         0, 0, 1,
    -2.5, 0, .25,       Math.sin(tailIncline), 0, Math.cos(tailIncline),
    1.75, 0, 0,         1, 0, 0,
    1.75, -.5, -.5,     1, 0, 0,
    1.75, .5, -.5,      1, 0, 0,
    -1.75, 0, 0,        -1, 0, 0,
    -1.75, -.5, -.5,    -1, 0, 0,
    -1.75, .5, -.5,     -1, 0, 0,
]

const idxs = [
    0, 1, 2, // board
    0, 1, 3,
    1, 3, 4,
    3, 4, 5,
    6, 7, 8, // trucks
    9, 10, 11,
]

export class SkateboardRenderer {

    private canvas: HTMLCanvasElement;
    private gl: WebGL2RenderingContext;

    private program: WebGLProgram = 0;
    private rotationLoc: WebGLUniformLocation = 0;
    private vertexBuffer: WebGLBuffer = 0;
    private indexBuffer: WebGLBuffer = 0;

    constructor() {
        this.canvas = document.createElement("canvas");
        if (!this.canvas) {
            throw new Error("Could not create canvas");
        }
        let gl = this.canvas.getContext("webgl2");
        if (!gl) {
            throw new Error("Cannot create webgl2 context");
        }
        this.gl = gl;
        gl.enable(gl.DEPTH_TEST);
        gl.enable(gl.DITHER);
        gl.clearColor(1, 1, 1, 0);
        gl.depthMask(true);
        gl.depthFunc(gl.LEQUAL);
        gl.depthRange(0.0, 1.0);

        this.compileShader();
        this.getUniforms();
        this.setupBuffers();
    }

    public drawSkateboard(ctx: CanvasRenderingContext2D, rotation: Quaternion, box: [number, number, number, number]) {
        this.canvas.width = box[2] - box[0];
        this.canvas.height = box[3] - box[1];
        this.gl.viewport(0, 0, this.canvas.width, this.canvas.height);
        this.gl.clear(this.gl.COLOR_BUFFER_BIT);
        this.gl.clear(this.gl.DEPTH_BUFFER_BIT);

        this.gl.uniform4fv(this.rotationLoc, rotation);

        // this.gl.bindBuffer(this.gl.ARRAY_BUFFER, this.vertexBuffer);
        // this.gl.bindBuffer(this.gl.ELEMENT_ARRAY_BUFFER, this.indexBuffer);
        this.gl.drawElements(this.gl.TRIANGLES, idxs.length, this.gl.UNSIGNED_SHORT, 0);
        this.gl.finish();
        ctx.drawImage(this.canvas, box[0], box[1]);
    }

    private setupBuffers() {
        let buffer = this.gl.createBuffer();
        if (!buffer) {
            throw new Error("Could not create vertex buffer");
        }
        this.gl.bindBuffer(this.gl.ARRAY_BUFFER, buffer);
        this.gl.bufferData(this.gl.ARRAY_BUFFER, new Float32Array(verts), this.gl.STATIC_DRAW);
        //this.gl.bindBuffer(this.gl.ARRAY_BUFFER, null);
        this.vertexBuffer = buffer;

        let indexes = this.gl.createBuffer();
        if (!indexes) {
            throw new Error("Could not create element buffer");
        }
        this.gl.bindBuffer(this.gl.ELEMENT_ARRAY_BUFFER, indexes);
        this.gl.bufferData(this.gl.ELEMENT_ARRAY_BUFFER, new Uint16Array(idxs), this.gl.STATIC_DRAW);
        //this.gl.bindBuffer(this.gl.ELEMENT_ARRAY_BUFFER, null);
        this.indexBuffer = indexes;

        let pos = this.gl.getAttribLocation(this.program, "vertPos");
        this.gl.vertexAttribPointer(pos, 3, this.gl.FLOAT, false, 6 * 4, 0);
        this.gl.enableVertexAttribArray(pos);
        
        let norm = this.gl.getAttribLocation(this.program, "vertNorm");
        this.gl.vertexAttribPointer(norm, 3, this.gl.FLOAT, false, 6 * 4, 3 * 4);
        this.gl.enableVertexAttribArray(norm);
    }

    private getUniforms() {
        let rotation = this.gl.getUniformLocation(this.program, "rotation");
        if (!rotation) {
            throw new Error("Could not get rotation uniform");
        }
        this.rotationLoc = rotation;
    }
    
    private loadShader(type: number, source: string) {
        const shader = this.gl.createShader(type);
        if (!shader) {
            return `Could not create shader type ${type}`;
        }
    
        this.gl.shaderSource(shader, source);
    
        this.gl.compileShader(shader);
    
        if (!this.gl.getShaderParameter(shader, this.gl.COMPILE_STATUS)) {
            let log = this.gl.getShaderInfoLog(shader) as string;
            this.gl.deleteShader(shader);
            return log;
        }

        return shader;
    }

    private compileShader() {
        const vertexShader = this.loadShader(this.gl.VERTEX_SHADER, `#version 300 es
precision highp float;

uniform vec4 rotation;

in vec3 vertPos;
in vec3 vertNorm;
out vec3 norm;

vec4 qMultiply(vec4 a, vec4 b) {
    return vec4(
        a.x * b.x - a.y * b.y - a.z * b.z - a.w * b.w,
        a.x * b.y + a.y * b.x + a.z * b.w - a.w * b.z,
        a.x * b.z - a.y * b.w + a.z * b.x + a.w * b.y,
        a.x * b.w + a.y * b.z - a.z * b.y + a.w * b.x
    );
}

vec3 rotate(vec3 vector) {
    return qMultiply(
        qMultiply(
            rotation, 
            vec4(0.0, vector)
        ), 
        vec4(rotation.x, -rotation.yzw)
    ).yzw;
}

void main() {
    vec3 rotated = 0.4 * rotate(vertPos);
    norm = rotate(vertNorm);
    gl_Position = vec4(rotated, 1.0);
}
`);
        const fragmentShader = this.loadShader(this.gl.FRAGMENT_SHADER, `#version 300 es
precision highp float;

in vec3 norm;
out vec4 color;

void main() {
    if (norm.z > 0.0) {
        color = vec4(1, 0, 0, 1);
    } else {
        color = vec4(0, 0, 0, 1);
    }
}
        `);

        if (typeof vertexShader === "string" || typeof fragmentShader === "string") {
            throw new Error(`Error Compiling shader
vertex shader: ${vertexShader}
fragment shader: ${fragmentShader}
            `)
        }
        console.log(vertexShader, fragmentShader)

        let program = this.gl.createProgram();
        if (!program) {
            throw new Error("Could not create program")
        }
        this.gl.attachShader(program, vertexShader);
        this.gl.attachShader(program, fragmentShader);
        this.gl.linkProgram(program);
    
        if (!this.gl.getProgramParameter(program, this.gl.LINK_STATUS)) {
            throw new Error(`Linker error
${this.gl.getProgramInfoLog(program)}
`);
        }
        this.gl.useProgram(program);
        this.program = program;
    }
}
