from flask import Flask, request, Response, send_file
from flask_sockets import Sockets

import os
import path
import json 

from video import Video, get_all_videos, youtube_dl
from labels import save_label, get_label_ids, get_label, delete_label

app = Flask(__name__)
sockets = Sockets(app)

@app.route("/videos", methods=["GET"])
def list_videos():
    return json.dumps({"videos": list(map(lambda v: {"title": v.title, "id": v.id, "progress": 100}, get_all_videos()))})

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

@app.route("/video", methods=["DELETE"])
def delete_video():
    video_id = request.args.get("id")
    videos = [ v for v in get_all_videos() if v.id == video_id]
    if len(videos) == 0:
        return Response(response="No videos match that id", status=404)
    for video in videos:
        video.delete()
    return "\{\}" 

@sockets.route("/video-download")
def download_video(ws):
    link = ws.receive()
    try:
        link = json.loads(link)
        for progress in youtube_dl(link["link"]):
            ws.send(progress)
        ws.close()
    except:
        ws.send('{"error":"Could not download"}')
        ws.close()

@app.route("/skateboards", methods=["POST"])
def add_label():
    try:
        new_id = save_label(request.json)
        return json.dumps({
            "success": True,
            "id": new_id
        })
    except Exception as e:
        print(e)
        return Response(
            response=json.dumps({
                "success": False,
                "error": "there was an error storing the label"
            }), status=400
        )

@app.route("/skateboards", methods=["GET"])
def get_labels():
    label_id = request.args.get("id")
    if label_id is None:
        return json.dumps({
            "labels": list(sorted(get_label_ids()))
        })
    else: 
        return get_label(label_id)

@app.route("/skateboards", methods=["DELETE"])
def label_delete():
    delete_label(request.args.get("id"))
    return '{"success": true}'
    

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

if __name__ == "__main__":
    from gevent import pywsgi
    from geventwebsocket.handler import WebSocketHandler
    server = pywsgi.WSGIServer(('', 5000), app, handler_class=WebSocketHandler)
    server.serve_forever()