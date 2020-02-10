# SkateboardVision

Identify skateboards in images, and track their rotation.

## Data Collection

The underlying solution to this problem is a convolutional neural network. Before we train the CNN we must collect a dataset that is representitive of the types of images it 
will be used for.

Our dataset will come from youtube. We use youtube-dl to store the videos on the server, then stream the video to the users to pause and clip frames from. The video must be streamed from our HTTP origin because youtube's preferred embedding method (iframes) prevents extracting frames for security reasons.

We might also expand our dataset to reddit, by modifying the video discovery phase for Gif videos, or other formats

### Data collection flow

1. Video Discovery
    * Videos are discovered from sources like youtube or reddit, and recon is carried out to be able to proxy a stream through our origin
2. User Acquisition
    * Users visit the website and a jwt is generated for their session. 
    * In the future, they might optionally log in to receive some kind of reward
3. Video Selection  
    * The user selects a video
    * This might be automatically selected if they have clicked on a preshared link 
4. Frame identification 
    * The user pauses the video at an interesting time
5. Clipping
    * The user draws a box around the skateboard using the built in image clipper
6. Rotation
    * The user uses the mouse and scroll wheel to align the rotation of the skateboard to the video frame
7. Collection
    * The frame is uploaded to the server with the coordinates of the clip box, and rotation
    * After this, the next frame is automatically identified

### Server

to run the development environment, use

```
docker-compose up
```

The server is a golang server connected to a mysql database. It hosts the frontend, and has various queries 


### Frontend

The frontend is a react application. If you used docker-compose to start the server, then it is statically served at the website root, and running `yarn dev` will update your changes as you save files

