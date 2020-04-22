# SkateboardVision

Identify skateboards in images, and track their rotation. 

This is the data api for processing and annotating videos with bounding rectangles and orientation quaternions.

See the frontend in action at https://sbvision.kwylder.com or checkout the [code](https://github.com/kevinwylder/sbvision-frontend)

1. Build sbgetvid to push progress to SNS
2. Build Manager to enqueue batch processes 
7. Setup SNS to route topics back to user websocket

5. Build Batch cluster to handle gsbvid