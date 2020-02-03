# SkateboardVision

Identify skateboards in images, and track their rotation.

## Data Collection

The underlying solution to this problem is a convolutional neural network. Before we train the CNN we must collect a dataset that is representitive of the types of images it 
will be used for.

Our dataset will come from youtube. We use youtube-dl to store the videos on the server, then stream the video to the users to pause and clip frames from. The video must be streamed from our HTTP origin because youtube's preferred embedding method (iframes) prevents extracting frames for security reasons.

There are ethical concerns with this tool that I would like to acknowledge, but justify by this project being personal, and not for any kind of commercial application.

### Server

to run the development environment, use

```
docker-compose up
```

