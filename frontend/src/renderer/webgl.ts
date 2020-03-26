import { Quaternion } from "../math";
import { skateboardGeometry } from "./deck";

const DECK_GRAPHIC = 0;
const GRIP_TAPE = 1;
const BOARD_RAIL = 2;

export class SkateboardRenderer {

    private canvas: HTMLCanvasElement;
    private gl: WebGL2RenderingContext;

    private boardElmtsSize: number = 0;
    private boardVertsBuffer: WebGLBuffer = 0;
    private boardElmtsBuffer: WebGLBuffer = 0;
    private boardPermBufferSize: number = 0;
    private boardPermBuffer: WebGLBuffer = 0;

    private program: WebGLProgram = 0;
    private materialLoc: WebGLUniformLocation = 0;
    private rotationLoc: WebGLUniformLocation = 0;

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
        // setup the desired rendering size
        this.canvas.width = box[2] - box[0];
        this.canvas.height = box[3] - box[1];
        this.gl.viewport(0, 0, this.canvas.width, this.canvas.height);
        this.gl.clear(this.gl.COLOR_BUFFER_BIT);
        this.gl.clear(this.gl.DEPTH_BUFFER_BIT);

        // set the rotation quaternion
        this.gl.uniform4fv(this.rotationLoc, rotation);

        // draw the deck
        this.gl.bindBuffer(this.gl.ARRAY_BUFFER, this.boardVertsBuffer);
        this.gl.bindBuffer(this.gl.ELEMENT_ARRAY_BUFFER, this.boardElmtsBuffer);
        let pos = this.gl.getAttribLocation(this.program, "vertPos");
        this.gl.vertexAttribPointer(pos, 3, this.gl.FLOAT, false, 6 * 4, 0);
        this.gl.enableVertexAttribArray(pos);
        let norm = this.gl.getAttribLocation(this.program, "vertNorm");
        this.gl.vertexAttribPointer(norm, 3, this.gl.FLOAT, true, 6 * 4, 3 * 4);
        this.gl.enableVertexAttribArray(norm);

        // bottom right
        this.gl.uniform1i(this.materialLoc, DECK_GRAPHIC);
        this.gl.drawElements(this.gl.TRIANGLE_STRIP, this.boardElmtsSize, this.gl.UNSIGNED_SHORT, 0);
        // top right
        this.gl.uniform1i(this.materialLoc, GRIP_TAPE);
        this.gl.drawElements(this.gl.TRIANGLE_STRIP, this.boardElmtsSize, this.gl.UNSIGNED_SHORT, 0);
        // set the rotation quaternion to the other side
        this.gl.uniform4fv(this.rotationLoc, [-rotation[3], rotation[2], -rotation[1], rotation[0]]);
        // bottom left
        this.gl.uniform1i(this.materialLoc, DECK_GRAPHIC);
        this.gl.drawElements(this.gl.TRIANGLE_STRIP, this.boardElmtsSize, this.gl.UNSIGNED_SHORT, 0);
        // top left
        this.gl.uniform1i(this.materialLoc, GRIP_TAPE);
        this.gl.drawElements(this.gl.TRIANGLE_STRIP, this.boardElmtsSize, this.gl.UNSIGNED_SHORT, 0);

        // draw the perimeter
        this.gl.bindBuffer(this.gl.ARRAY_BUFFER, this.boardPermBuffer);
        this.gl.bindBuffer(this.gl.ELEMENT_ARRAY_BUFFER, null);
        pos = this.gl.getAttribLocation(this.program, "vertPos");
        this.gl.vertexAttribPointer(pos, 3, this.gl.FLOAT, false, 6 * 4, 0);
        this.gl.enableVertexAttribArray(pos);
        norm = this.gl.getAttribLocation(this.program, "vertNorm");
        this.gl.vertexAttribPointer(norm, 3, this.gl.FLOAT, true, 6 * 4, 3 * 4);
        this.gl.enableVertexAttribArray(norm);
        this.gl.bindBuffer(this.gl.ARRAY_BUFFER, null);

        this.gl.uniform1i(this.materialLoc, BOARD_RAIL);
        this.gl.drawArrays(this.gl.TRIANGLE_STRIP, 0, this.boardPermBufferSize);
        // set the rotation quaternion to the other side
        this.gl.uniform4fv(this.rotationLoc, rotation);
        this.gl.drawArrays(this.gl.TRIANGLE_STRIP, 0, this.boardPermBufferSize);

        ctx.drawImage(this.canvas, box[0], box[1]);
    }

    private setupBuffers() {
        [ 
            this.boardVertsBuffer, 
            this.boardElmtsBuffer, 
            this.boardElmtsSize,
            this.boardPermBuffer,
            this.boardPermBufferSize
        ] = skateboardGeometry(this.gl);

    }

    private getUniforms() {
        let material = this.gl.getUniformLocation(this.program, "materialType");
        if (!material) {
            throw new Error("Could not get rotation uniform");
        }
        this.materialLoc = material;
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
precision lowp int;

const float n = -8.0;
const float f = -10.0;
const mat4 projection = mat4(
    vec4(n, 0, 0, 0),
    vec4(0, n, 0, 0),
    vec4(0, 0, -(f + n) / (f - n), -1),
    vec4(0, 0, 2.0 * f * n / (f - n), 0)
);

uniform int materialType;
uniform vec4 rotation;

in vec3 vertPos;
in vec3 vertNorm;

out vec3 norm;
out vec3 pos;

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
            vec4(0.0, vector * vec3(1, -1, -1))
        ), 
        vec4(rotation.x, -rotation.yzw)
    ).yzw * vec3(1, -1, -1);
}

void main() {
    vec3 translation = vec3(0, 0, 0);
    if (materialType == ${GRIP_TAPE}) {
        translation = vec3(0, 0, 0.01);
    } else if (materialType == ${DECK_GRAPHIC}) {
        translation = vec3(0, 0, -0.01);
    }

    pos = rotate(vertPos + translation);
    norm = rotate(vertNorm);
    pos.z += (n + f) / 2.0;
    gl_Position = projection * vec4(pos, 1.0);
}
`);
        const fragmentShader = this.loadShader(this.gl.FRAGMENT_SHADER, `#version 300 es
precision highp float;
precision lowp int;

uniform int materialType;

const float n = -8.0;
const float f = -10.0;

const vec3 lightPos = vec3(0, 0, 4);

in vec3 pos;
in vec3 norm;

out vec4 color;

void main() {
    vec3 camera = normalize(-pos);
    vec3 light = normalize(lightPos - pos);

    float diffuse = max(dot(norm, light), 0.0);

    vec3 materialColor = vec3(0, 0, 0);
    if (materialType == ${DECK_GRAPHIC}) {
        materialColor = vec3(diffuse * .4 + .6, 0, 0);
    } else if (materialType == ${BOARD_RAIL}) {
        materialColor = vec3(0.97, 0.87, 0.4) * (diffuse + .5);

    } else if (materialType == ${GRIP_TAPE}) {
        vec3 reflection = -light - 2.0 * dot(-light, norm) * norm;
        float base = dot(reflection, camera);
        if (base > 0.0) base = 0.0;
        float specularColor = pow(-base, 90.0);

        materialColor = vec3(1, 1, 1) * specularColor;
    }

    color = vec4(materialColor, 1);
}
        `);

        if (typeof vertexShader === "string" || typeof fragmentShader === "string") {
            throw new Error(`Error Compiling shader
vertex shader: ${vertexShader}
fragment shader: ${fragmentShader}
            `)
        }

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
