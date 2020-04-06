# SkateboardVision

Identify skateboards in images, and track their rotation. 

This is the data api for annotating videos with bounding rectangles and orientation quaternions.

See the frontend in action at https://sbvision.kwylder.com or checkout the [code](https://github.com/kevinwylder/sbvision-frontend)

### Development Environment

```
docker-compose up
```

The docker-compose.yml is an environment that includes some helpful dev features.
* mysql database with the latest schema (you can check it out at port 3369 if you're into that)
* image asset data storage in ./data directory
* filesystem watcher thanks to github.com/codegangsta/gin 




