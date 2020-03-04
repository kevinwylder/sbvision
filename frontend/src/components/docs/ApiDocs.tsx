import * as React from 'react';

import "./docs.css";

export function ApiDocs() {

    return <div className="docs">
        <h2> API Documentation </h2>

        <p style={{textIndent: "40px"}}> The api allows you to get all the annotations for a video, and image resources associated with each frame. We recommend using the video viewer to review the content of a video before using it in the dataset, as this is a crowd sourced project</p> 

        <div className="docs-endpoint">
            <h3><b> GET /api/frames</b>?video=<span className="highlight">:id</span>&amp;offset=<span className="highlight">:offset</span> </h3>
            <div className="fields">
                <div className="key"><span>video <i>(required)</i></span></div> <div><p>
                    Returns frames for the given video ID. Get this from the desired video's URL <span className="code">/videos/:ID</span>
                </p></div>
                <div className="key"><span>offset <i>(optional)</i></span></div> <div><p>
                    This optional field is used to pagenate the output from this endpoint. If the result from hitting the endpoint has the <span className="code">isTruncated</span> flag set to true, then hit the endpoint again with this parameter from the <span className="code">nextOffset</span> field.
                </p></div>
            </div>
            <div className="example">
                <pre> $ curl https://sbvision.kwylder.com/api/frames?video=1 | jq </pre>
                <pre> {`
{
    "isTruncated": false,
    "nextOffset": 0,
    "frames": [
    {
        "id": 2395,
        "time": 903,
        "bounds": [
        {
            "id": 1892,
            "width": 121,
            "height": 55,
            "x": 356,
            "y": 192,
            "rotations": [
            {
                "r": 0.5323442542091413,
                "i": -0.24405891699300233,
                "j": -0.020520261420443997,
                "k": -0.8103232434763015
            }
            ]
        } 
        ]
    ]
}
`}
                </pre>
            </div>

            <div className="fields">
                <div className="key"><span>isTruncated</span></div> <div>
                    Set to true if there are more frames available for this video. Use the <span className="code">nextOffset</span> to get the rest of the frames.
                </div>
                <div className="key"><span>nextOffset</span></div> <div>
                    if <span className="code">isTruncated</span>, use this value as the <span className="code">offset</span> for the next request to this video endpoint
                </div>
                <div className="key"><span>frames <i>(array)</i></span></div> <div>
                    an array of objects in time sorted order. Each object represents a frame that has been uploaded for this video. When the video is paused by the application player, the frame is uploaded and tagged. Duplicate frames are avoided by hashing the image data. If no frames are available, this field is <span className="code">null</span>
                
                <div className="fields">
                    <div className="key"><span>id</span></div>
                    <div>
                        A unique frame identifier for this video and timestamp. Use this id to get the image at the <span className="code">/api/images?frame=</span> endpoint.
                    </div>
                    <div className="key"><span>time</span></div>
                    <div>
                        The number of milliseconds into the video this frame was captured at.
                    </div>
                    <div className="key"><span>bounds <i>(array)</i></span></div>
                    <div>
                        All the bounding boxes that have been identified in this frame. Boxes are described in terms of pixels relative to the image resource.

                    <div className="fields">
                        <div className="key"><span>id</span></div>
                        <div>
                            The unique identifier of this bounding box. You can pass this id to the <span className="code">/api/images?bound=</span> endpoint to get a cropped image of this bounding box.
                        </div>
                        <div className="key"><span>x</span></div>
                        <div>
                            The x offset from the upper left corner this bounding box starts at, in pixels
                        </div>
                        <div className="key"><span>y</span></div>
                        <div>
                            The y offset from the upper left corner this bounding box starts at, in pixels
                        </div>
                        <div className="key"><span>width</span></div>
                        <div>
                            The width in pixels of the bounding box
                        </div>
                        <div className="key"><span>height</span></div>
                        <div>
                            The height in pixels of the bounding box
                        </div>
                        <div className="key"><span>rotations <i>(array)</i></span></div>
                        <div>
                            All the skateboard orientations that have been provided for this bounding box, represented as a <a href="https://eater.net/quaternions">unit quaternion</a>. If no rotations have been provided, this array is null.
                            
                            <div className="fields">
                                <div className="key"><span>r</span></div>
                                <div>
                                    the real component of the quaternion
                                </div>
                                <div className="key"><span>i</span></div>
                                <div>
                                    the i component of the quaternion
                                </div>
                                <div className="key"><span>j</span></div>
                                <div>
                                    the j component of the quaternion
                                </div>
                                <div className="key"><span>k</span></div>
                                <div>
                                    the k component of the quaternion
                                </div>
                            </div>
                        </div></div>
                    </div></div>
                </div></div>
            </div>


    <div className="docs-endpoint">
        <h3><b> GET /api/image</b>?frame=<span className="highlight">:id</span>&amp;bound=<span className="highlight">:id</span> </h3>
        This endpoint gets a PNG image of the entire frame, or bounded section of the frame. The two url parameters are mutually exclusive.
        <div className="fields">
            <div className="key"><span>frame <i>(optional)</i></span></div> <div><p>
                If provided, the response will be an image of the entire frame
            </p></div>
            <div className="key"><span>bound <i>(optional)</i></span></div> <div><p>
                If provided, only the area described by this bound will be returned
            </p></div>
        </div>
    </div>

    <h1> Footnote on Quaternions </h1>

    <div>
        This project uses quaternions to describe skateboard orientations because they have a uniform covering of the rotation space (meaning the norm of their difference is a good error function), and the web application does not suffer from Gimble Lock. 
        Quaternions, however, do take some special care to use in the application. Please note:

        <ol>
            <li> Each rotation has two unit quaternions that describe the position due to the "double covering" of this space </li>
            <li> A 180 degree rotation along the medium axis of rotation (a pop-shuvit) reveals that the skateboard has an axis of symmetry, leading to another "double covering" of the space</li>
            <li> A <a href="https://math.stackexchange.com/questions/90081/quaternion-distance">Euclidean norm</a> will only yield a small error for one of these quaternions</li>
        </ol>

        I recommend <a href="https://eater.net/quaternions">Ben Eater and Grant Sanderson's description web demo</a> of their properties to get a better understanding of how quaternions describe rotations.
    </div>


    </div>
}