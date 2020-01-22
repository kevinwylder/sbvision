from flask import Flask, request, Response, send_file
from video import Video, get_all_videos
import os
import path
import json 

app = Flask(__name__)

@app.route("/videos", methods=["GET"])
def list_videos():
    return json.dumps({"videos": get_all_videos()})

@app.route("/video", methods=["GET"])
def stream_video():
    video_id = request.args.get("id")
    videos = [ v for v in get_all_videos() if v.id == video_id ]
    if len(videos) == 0:
        return Response(response="No video files match that ID", status=404)

    data, start, size = videos[0].get_video_range(section=request.headers.get("Range", None))
    response = Response(
        response=data,
        status=206,
        mimetype="video/mp4",
        direct_passthrough=True
    )
    response.headers.add('Content-Range', 'bytes {0}-{1}/{2}'.format(start, start + len(data) - 1, size))
    return response

@app.after_request
def after_request(response):
    response.headers.add('Accept-Ranges', 'bytes')
    return response


@app.route("/", methods=["GET"])
def index():
    return send_file("frontend/dist/index.html")

@app.route("/style.css")
def style():
    return send_file("frontend/dist/clipper.css", cache_timeout=1)

@app.route("/main.js")
def bundle():
    return send_file("frontend/dist/main.js", cache_timeout=1)

@app.route("/react.js")
def react():
    return send_file("frontend/node_modules/react/umd/react.development.js")

@app.route("/react-dom.js")
def react_dom():
    return send_file("frontend/node_modules/react-dom/umd/react-dom.development.js")

if __name__ == '__main__':
    app.run(threaded=True)