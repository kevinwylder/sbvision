import subprocess
import re
import json
import os

video_dir = os.environ["DATA_DIR"] + "/videos"

def get_all_videos():
    return [ Video(filename) for filename in os.listdir(video_dir) ]

class Video(object):

    def __init__(self, filename):
        self.filename = video_dir + "/" + filename
        title_id = ".".join(filename.split(".")[:-1]).split("-")
        self.title = "".join(title_id[:-1])
        self.id = title_id[-1]

    def get_video_range(self, section=None):
        full_size = os.stat(self.filename).st_size
        start = 0
        length = 1024 * 100

        if section is not None:
            start, end = re.search(r"(\d+)-(\d*)", section).groups()
            if start is None or end is None:
                start = 0
                length = 1024 * 100
            else:
                start = max(min(int(start), full_size), 0)
                if end is "":
                    length = min(full_size - start, 1024 * 100)
                else:
                    end = max(min(int(end), full_size), 0)

        with open(self.filename, 'rb') as f:
            f.seek(start)
            data = f.read(length)
            return data, start, full_size
    
    def delete(self):
        os.remove(self.filename)

# reads up to carriage return
def read_line(reader):
    output = ""
    char = reader.read(1).decode('utf8')
    while len(char) != 0:
        if char == '\r' or char == '\n':
            break
        output += char
        char = reader.read(1).decode('utf8')
    return output

destination_matcher = re.compile(r'Destination: (.*.mp4)')
progress_matcher = re.compile(r'download]\s+(\d+)\.\d% of ')
def youtube_dl(link):
    process = subprocess.Popen(["youtube-dl", link], stdout=subprocess.PIPE, cwd=video_dir)
    video = None
    progress = 0
    while True:
        output = read_line(process.stdout)
        if output == '' and process.poll() is not None:
            return
        if output:
            output = str(output)
            if "has already been downloaded" in output:
                yield json.dumps({
                    'error': "already downloaded"
                })
                return
            if video is None:
                destination_match = destination_matcher.search(output)
                if destination_match is not None:
                    video = Video(destination_match.group(1))
            else:
                progress_match = progress_matcher.search(output)
                if progress_match is not None:
                    progress = int(progress_match.group(1))
            if video is not None:
                yield json.dumps({
                    'title': video.title,
                    'id': video.id,
                    'progress': progress
                })
    yield json.dumps({
        'title': video.title,
        'id': video.id,
        'progress': 100
    })