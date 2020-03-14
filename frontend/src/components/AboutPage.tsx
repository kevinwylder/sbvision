import * as React from 'react';

import { Link } from 'react-router-dom';

export function AboutPage() {
    return <div className="about-page">
        <h2> About the Skateboard Vision Project </h2>
        <p>
            We are collecting a dataset of skateboard images to train a classifier to identify skateboards, and what orientation
            they are in. Once the model is robust at detecting skateboards in images, we will move to video clips to identify
            tricks. To really understand the dataset, check out our <Link to="/explore">dataset visualization</Link>.
        </p>

        <h2> Applications for the Skateboard Visison Project </h2>

        <p> 
            This project, if successful, will be able to extract physical context out of skateboarding videos. With the help of
            <a href="https://medium.com/tensorflow/real-time-human-pose-estimation-in-the-browser-with-tensorflow-js-7dd0bc881cd5"> Pose Estimation </a>
            we could enable a skateboard version of <a href="https://www.youtube.com/watch?v=o_DhNqHazKY"> this project </a> to create a hyper-realistic skateboarding video game.
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
            <li>Rotation Match: Click on the rotating skateboard in the upper left corner of the screen to enter the frame review screen. There, you can add rotations to existing bounds </li>
            <li>Data Verification: Each bounding box and orientation needs to be reviewed for quality. Soon, we'll add 
                functionality to delete bad frames from the dataset with the visualizor </li>
        </ol>

        <p>
            You can also help us out by providing suggestions, and submitting bug reports on <a href="https://github.com/kevinwylder/sbvision/issues">github</a>
        </p>

        <h2> Like this dataset? You can have it! </h2>

        <p>
            The data is publicly available, check out <Link to="/api-docs"> our documentation </Link> on how to get the dataset. 
            We're also working on a python package to wrap the http api and download local datasets.
        </p>

    </div>
}