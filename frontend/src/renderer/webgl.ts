import { Quaternion, bezier, exponentNorm } from "../math";
import { skateboardGeometry } from "./deck";
import { wheelGeometry } from "./wheel";
import { describeBox } from "../math/layout";

const DECK_GRAPHIC = 0;
const GRIP_TAPE = 1;
const BOARD_RAIL = 2;
const WHEEL = 3;

export class SkateboardRenderer {

    private deckParameters = {
        thickness: 0.021,
        width: 1.9,
        height: .5,
        perimeterStepFactor: .02,
        perimeterZNormalSize: .2,
        tailEllipseExponent: 2.4,
        tailLiftX: .6,
        tailLiftFactor: .6,
        tailStretchFactor: 1.4,
        railTailTransition: bezier(1, 1, .17, 0),
        railLiftFactor: .34,
        railLiftX: .5,
        wheelX: .52,
        wheelZ: -.12,
        wheelNormTaper: .2,
        wheelScale: 0.05,
        wheelLength: .04,
        wheelRadius: 0.05,
        wheelMeshRes: 12,
    };

    private canvas: HTMLCanvasElement;
    private gl: WebGL2RenderingContext;

    private program: WebGLProgram = 0;
    private materialLoc: WebGLUniformLocation = 0;
    private rotationLoc: WebGLUniformLocation = 0;
    private translationLoc: WebGLUniformLocation = 0;

    private boardVertsBuffer: WebGLBuffer = 0;
    private boardElmtsBuffer: WebGLBuffer = 0;
    private boardElmtsSize: number = 0;

    private boardPermBuffer: WebGLBuffer = 0;
    private boardPermBufferSize: number = 0;

    private wheelVertsBuffer: WebGLBuffer = 0;
    private wheelElmtsBuffer: WebGLBuffer = 0;
    private wheelBufferSize: number = 0;

    private wheelCapBuffer: WebGLBuffer = 0;
    private wheelCapSize: number = 0;

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
        gl.clearColor(1, 1, 1, 1);
        gl.depthMask(true);
        gl.depthFunc(gl.LEQUAL);
        gl.depthRange(0.0, 1.0);

        this.compileShader();
        this.getUniforms();
        this.setupBuffers();
    }

    public drawSkateboard(ctx: CanvasRenderingContext2D, rotation: Quaternion, [x0, x1, x2, x3]: [number, number, number, number]) {
        let { left, top, width, height } = describeBox([x0, x1, x2, x3]);
        let long = Math.max(height, width);

        // setup the desired rendering size
        this.canvas.width = long;
        this.canvas.height = long;
        this.gl.viewport(0, 0, long, long);
        this.gl.clear(this.gl.COLOR_BUFFER_BIT);
        this.gl.clear(this.gl.DEPTH_BUFFER_BIT);

        // set the rotation quaternion
        this.gl.uniform4fv(this.rotationLoc, rotation);

        // draw the deck
        this.bindBuffer(this.boardVertsBuffer, this.boardElmtsBuffer);
        this.gl.uniform1i(this.materialLoc, DECK_GRAPHIC);
        this.gl.drawElements(this.gl.TRIANGLE_STRIP, this.boardElmtsSize, this.gl.UNSIGNED_SHORT, 0);
        this.gl.uniform1i(this.materialLoc, GRIP_TAPE);
        this.gl.drawElements(this.gl.TRIANGLE_STRIP, this.boardElmtsSize, this.gl.UNSIGNED_SHORT, 0);
        this.gl.uniform4fv(this.rotationLoc, [-rotation[3], rotation[2], -rotation[1], rotation[0]]);
        this.gl.uniform1i(this.materialLoc, DECK_GRAPHIC);
        this.gl.drawElements(this.gl.TRIANGLE_STRIP, this.boardElmtsSize, this.gl.UNSIGNED_SHORT, 0);
        this.gl.uniform1i(this.materialLoc, GRIP_TAPE);
        this.gl.drawElements(this.gl.TRIANGLE_STRIP, this.boardElmtsSize, this.gl.UNSIGNED_SHORT, 0);

        // draw the perimeter
        this.bindBuffer(this.boardPermBuffer, null);
        this.gl.uniform1i(this.materialLoc, BOARD_RAIL);
        this.gl.drawArrays(this.gl.TRIANGLE_STRIP, 0, this.boardPermBufferSize);
        this.gl.uniform4fv(this.rotationLoc, rotation);
        this.gl.drawArrays(this.gl.TRIANGLE_STRIP, 0, this.boardPermBufferSize);

        // draw the wheels
        let { wheelX, wheelLength, wheelZ } = this.deckParameters;
        let wheelY = this.deckParameters.height / 2 - wheelLength;
        this.bindBuffer(this.wheelVertsBuffer, this.wheelElmtsBuffer);
        this.gl.uniform1i(this.materialLoc, WHEEL);
        this.gl.uniform3fv(this.translationLoc, [-wheelX, -wheelY, wheelZ]);
        this.gl.drawElements(this.gl.TRIANGLE_STRIP, this.wheelBufferSize, this.gl.UNSIGNED_SHORT, 0);
        this.gl.uniform3fv(this.translationLoc, [ wheelX,  wheelY, wheelZ]);
        this.gl.drawElements(this.gl.TRIANGLE_STRIP, this.wheelBufferSize, this.gl.UNSIGNED_SHORT, 0);
        this.gl.uniform3fv(this.translationLoc, [ wheelX, -wheelY, wheelZ]);
        this.gl.drawElements(this.gl.TRIANGLE_STRIP, this.wheelBufferSize, this.gl.UNSIGNED_SHORT, 0);
        this.gl.uniform3fv(this.translationLoc, [-wheelX,  wheelY, wheelZ]);
        this.gl.drawElements(this.gl.TRIANGLE_STRIP, this.wheelBufferSize, this.gl.UNSIGNED_SHORT, 0);

        // draw wheel caps
        this.bindBuffer(this.wheelCapBuffer, null);
        let outer = this.deckParameters.height / 2;
        let inner = this.deckParameters.height / 2 - 2 * wheelLength;
        this.gl.uniform3fv(this.translationLoc, [-wheelX, -outer, wheelZ]);
        this.gl.drawArrays(this.gl.TRIANGLE_FAN, this.wheelCapSize, this.wheelCapSize);
        this.gl.uniform3fv(this.translationLoc, [-wheelX, -inner, wheelZ]);
        this.gl.drawArrays(this.gl.TRIANGLE_FAN, 0, this.wheelCapSize);
        this.gl.uniform3fv(this.translationLoc, [-wheelX, outer, wheelZ]);
        this.gl.drawArrays(this.gl.TRIANGLE_FAN, 0, this.wheelCapSize);
        this.gl.uniform3fv(this.translationLoc, [-wheelX, inner, wheelZ]);
        this.gl.drawArrays(this.gl.TRIANGLE_FAN, this.wheelCapSize, this.wheelCapSize);
        this.gl.uniform3fv(this.translationLoc, [wheelX, outer, wheelZ]);
        this.gl.drawArrays(this.gl.TRIANGLE_FAN, 0, this.wheelCapSize);
        this.gl.uniform3fv(this.translationLoc, [wheelX, inner, wheelZ]);
        this.gl.drawArrays(this.gl.TRIANGLE_FAN, this.wheelCapSize, this.wheelCapSize);
        this.gl.uniform3fv(this.translationLoc, [wheelX, -outer, wheelZ]);
        this.gl.drawArrays(this.gl.TRIANGLE_FAN, this.wheelCapSize, this.wheelCapSize);
        this.gl.uniform3fv(this.translationLoc, [wheelX, -inner, wheelZ]);
        this.gl.drawArrays(this.gl.TRIANGLE_FAN, 0, this.wheelCapSize);

        // transfer to 2d context
        let x = Math.max(0, (height - width) / 2);
        let y = Math.max(0, (width - height) / 2);
        ctx.drawImage(this.canvas, x, y, width, height, left, top, width, height);
    }

    private setupBuffers() {
        [ 
            this.boardVertsBuffer, 
            this.boardElmtsBuffer, 
            this.boardElmtsSize,
            this.boardPermBuffer,
            this.boardPermBufferSize
        ] = skateboardGeometry(this.gl, this.deckParameters);

        [
            this.wheelVertsBuffer,
            this.wheelElmtsBuffer,
            this.wheelBufferSize,
            this.wheelCapBuffer,
            this.wheelCapSize,
        ] = wheelGeometry(this.gl, this.deckParameters);
    }
    
    private bindBuffer(vert: WebGLBuffer, index: WebGLBuffer|null) {
        this.gl.bindBuffer(this.gl.ARRAY_BUFFER, vert);
        this.gl.bindBuffer(this.gl.ELEMENT_ARRAY_BUFFER, index);
        let pos = this.gl.getAttribLocation(this.program, "vertPos");
        this.gl.vertexAttribPointer(pos, 3, this.gl.FLOAT, false, 6 * 4, 0);
        this.gl.enableVertexAttribArray(pos);
        let norm = this.gl.getAttribLocation(this.program, "vertNorm");
        this.gl.vertexAttribPointer(norm, 3, this.gl.FLOAT, true, 6 * 4, 3 * 4);
        this.gl.enableVertexAttribArray(norm);
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
        let translation = this.gl.getUniformLocation(this.program, "translationAdditional");
        if (!translation) {
            throw new Error("Could not get translation uniform");
        }
        this.translationLoc = translation;
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
uniform vec3 translationAdditional;

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
        translation = vec3(0, 0, ${this.deckParameters.thickness / 2});
    } else if (materialType == ${DECK_GRAPHIC}) {
        translation = vec3(0, 0, -${this.deckParameters.thickness / 2});
    } else if (materialType == ${WHEEL}) {
        translation = translationAdditional;
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

const vec3 lightPos = vec3(1, 0, -4);

in vec3 pos;
in vec3 norm;

out vec4 color;

void main() {
    vec3 camera = normalize(vec3(0, 0, n) - pos);
    vec3 light = normalize(lightPos);

    float diffuse = max(dot(norm, light), 0.0);
    vec3 reflection = -light - 2.0 * dot(-light, norm) * norm;

    vec3 materialColor = vec3(0, 0, 0);
    if (materialType == ${DECK_GRAPHIC}) {
        float base = dot(reflection, camera);
        if (base > 0.0) base = 0.0;
        float specularColor = pow(-base, 90.0);
        materialColor = vec3(min(diffuse + specularColor, 1.0) * .4 + .6, specularColor, specularColor);

    } else if (materialType == ${BOARD_RAIL}) {
        materialColor = vec3(0.97, 0.87, 0.4) * (diffuse + .5);

    } else if (materialType == ${GRIP_TAPE}) {
        float base = dot(reflection, camera);
        if (base > 0.0) base = 0.0;
        float specularColor = pow(-base, 90.0);
        materialColor = vec3(1, 1, 1) * specularColor;

    } else if (materialType == ${WHEEL}) {
        materialColor = vec3(1, 1, 1) * (0.5 * diffuse + 0.5);

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
