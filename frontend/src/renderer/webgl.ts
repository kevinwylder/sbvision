

export class SkateboardRenderer {

    private canvas: HTMLCanvasElement;
    private gl: WebGL2RenderingContext;

    constructor() {
        this.canvas = document.createElement("canvas");
        let gl = this.canvas.getContext("webgl2");
        if (!gl) {
            throw new Error("Cannot create webgl2 context");
        }
        this.gl = gl;
        this.compileShader();
    }

    private compileShader() {
        const vertexShader = this.loadShader(this.gl.VERTEX_SHADER, `
attribute vec4 aVertexPosition;

        
        `);
        const fragmentShader = this.loadShader(this.gl.FRAGMENT_SHADER, `
        
        
        `);

        if (typeof vertexShader === "string" || typeof fragmentShader === "string") {
            console.log(`Error Compiling shader
vertex shader: ${vertexShader}
fragment shader: ${fragmentShader}
            `)
            return;
        }

        const shaderProgram = this.gl.createProgram();
        if (!shaderProgram) {
            console.log("Could not create program")
            return;
        }
        this.gl.attachShader(shaderProgram, vertexShader);
        this.gl.attachShader(shaderProgram, fragmentShader);
        this.gl.linkProgram(shaderProgram);
    
        if (!this.gl.getProgramParameter(shaderProgram, this.gl.LINK_STATUS)) {
            console.log(`Linker error
${this.gl.getProgramInfoLog(shaderProgram)}
`);
        }
        return shaderProgram;
    }
    
    //
    // creates a shader of the given type, uploads the source and
    // compiles it.
    //
    private loadShader(type: number, source: string) {
        const shader = this.gl.createShader(type);
        if (!shader) {
            return `Could not create shader type ${type}`;
        }
    
        this.gl.shaderSource(shader, source);
    
        this.gl.compileShader(shader);
    
        if (!this.gl.getShaderParameter(shader, this.gl.COMPILE_STATUS)) {
            this.gl.deleteShader(shader);
            return this.gl.getShaderInfoLog(shader) as string;
        }

        return shader;
    }
}
