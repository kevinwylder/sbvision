import os
import random
import json

label_dir = os.environ["DATA_DIR"] + "/labels"

def get_label_ids():
    return map(lambda name: int(name.split(".")[0]), os.listdir(label_dir))

def save_label(label):
    filename = str(max(get_label_ids()) + 1) + ".json"
    with open(label_dir + "/" + filename, "w") as f:
        json.dump(label, f)