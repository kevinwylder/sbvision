from flask import Flask, request, Response, send_file
import os
import path
import json 

from video import Video, get_all_videos, youtube_dl
from labels import save_label

app = Flask(__name__)

@app.route("/videos", methods=["GET"])
def list_videos():
    return json.dumps({"videos": list(map(lambda v: {"title": v.title, "id": v.id}, get_all_videos()))})

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

@app.route("/video", methods=["POST"])
def download_video():
    if "link" not in request.json:
        return "looking for a {\"link\":\"http://www.youtube.com...\"}"
    if youtube_dl(request.json["link"]) == 0:
        return "Success"
    else:
        return Response("Failed to download video", status=400)

@app.route("/skateboards", methods=["POST"])
def add_label():
    try:
        save_label(request.json)
        return json.dumps({
            "success": True
        })
    except Exception as e:
        print(e)
        return Response(
            response=json.dumps({
                "success": False,
                "error": "there was an error storing the label"
            }), status=400
        )

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