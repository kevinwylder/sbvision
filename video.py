import subprocess
import re
import os

video_dir = os.environ["DATA_DIR"] + "/videos"

def get_all_videos():
    return [ Video(filename) for filename in os.listdir(video_dir) ]

class Video(object):

    def __init__(self, filename):
        self.filename = video_dir + "/" + filename
        title_id = filename.split(".")[0].split("-")
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
            print(len(data), length)
            return data, start, full_size

def youtube_dl(link):
    return subprocess.Popen(["youtube-dl", link], cwd=video_dir).wait()