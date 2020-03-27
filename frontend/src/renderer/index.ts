import { SkateboardRenderer } from './webgl';

// setup the renderSkateboard to use '2d' ctx until 'webgl2' gl loads
import { renderSkateboard as renderLegacy } from './2d';
export let renderSkateboard = renderLegacy;
try {
    let renderer = new SkateboardRenderer();
    renderSkateboard = (...params) => renderer.drawSkateboard(...params);
} catch (e) {
    console.log(e);
    console.log("Defaulting to 2d renderer");
}