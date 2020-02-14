# SkateboardVision

Identify skateboards in images, and track their rotation. See it in action at https://sbvision.kwylder.com

## Data Collection

The underlying solution to this problem is a convolutional neural network. Before we train the CNN we must collect a dataset that is representitive of the types of images it 
will be used for.

Our dataset will come from youtube. We use youtube-dl to store the videos on the server, then stream the video to the users to pause and clip frames from. The video must be streamed from our HTTP origin because youtube's preferred embedding method (iframes) prevents extracting frames for security reasons.

We might also expand our dataset to reddit, by modifying the video discovery phase for Gif videos, or other formats

### Data collection flow

1. Video Discovery
    * Videos are discovered from sources like youtube or reddit, and recon is carried out to be able to proxy a stream through our origin
2. Video Selection  
    * The user selects a video
    * This might be automatically selected if they have clicked on a preshared link 
3. Frame identification 
    * The user pauses the video at an interesting time
    * The frame is uploaded and indexed on the server, returning a Frame ID
4. Clipping
    * The user draws a box around the skateboard using the built in image clipper
    * The bounds of this area are associated with the frame ID and stored in the database
    * After clipping, the user advances to the next frame (goto step 4)
5. Rotation
    * Once many frames and bounds have been extracted from the video, we find the rotation of individual bounds
    * The user aligns a 3D skateboard rendering with the skateboard in the image

### Server

to run the development environment, use

```
docker-compose up
```

The server is a golang server connected to a mysql database. It hosts the frontend, and has various queries 


### Frontend

The frontend is a react application. If you used docker-compose to start the server, then it is statically served at the website root, and running `yarn dev` will update your changes as you save files

