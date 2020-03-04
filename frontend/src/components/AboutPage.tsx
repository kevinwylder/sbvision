import * as React from 'react';

import { Link } from 'react-router-dom';

export function AboutPage() {
    return <div className="about-page">
        <h2> About the Skateboard Visision Project </h2>
        <p>
            3D Pose Estimation is a well studied field of Modern Image Processing. Using revolutionary 
            Neural Network technology and a dataset painstakingly collected on thousands of images, it 
            is now possible to very quickly detect physical wireframes of humans in images. 
            <a href="https://medium.com/tensorflow/real-time-human-pose-estimation-in-the-browser-with-tensorflow-js-7dd0bc881cd5"> Here </a> 
            is a good example of what this looks like.
        </p>

        <p>
            The Skateboard Vision Project hopes to extend this pose estimation process to the world of
            flip tricks. Using a similar technique, we are trying to detect the <b> Position </b> and
            <b> Orientation </b> of skateboards, given an image and a wireframe model. {/*To get a better
            understanding of what this means, check out our <Link to="/dataset"> Dataset Visualization </Link>. */}
        </p>

        <h2> Applications for the Skateboard Visison Project </h2>

        <p> 
            This project will be able to extract physical descriptions from real videos. With enough data,
            we can enable a skateboard version of <a href="https://www.youtube.com/watch?v=o_DhNqHazKY"> this project </a>
            to create a hyper-realistic skateboarding video game.
        </p>

        <p>
            Another application of the Skateboard Vision Project is as a learner's guide. This would work
            by having a camera focused on somebody doing skateboard tricks. It could be able to detect
            which tricks were done, and whether they were landed or not. This would allow for high level
            feedback, or even suggest tips to improve the chances of landing the trick.
        </p>

        <h2> How to Help </h2>

        <p>
            We need to review skateboarding clips to collect a dataset for this to be possible. There are
            3 steps to adding to this dataset, all of which you can help with.
        </p>

        <ol>
            <li>Bounding Boxes: Pause any <Link to="/videos"> video </Link>to draw a box around skateboards. </li>
            <li>Rotation Match: Head over to <Link to="/rotations"> the rotation matcher </Link> to mark the rotation of the skateboard </li>
            <li>Data Verification: (Coming Soon!) Each bounding box and orientation needs to be reviewed for quality </li>
        </ol>

        <p>
            You can also help us out by providing suggestions, and submitting bug reports on <a href="https://github.com/kevinwylder/sbvision/issues">github</a>
        </p>

        <h2> Like this dataset? You can have it! </h2>

        <p>
            The data is publicly available, check out <Link to="/api-docs"> our documentation </Link> on how to get the dataset. We're also working on a python package to wrap the http api and download local datasets.
        </p>

    </div>
}