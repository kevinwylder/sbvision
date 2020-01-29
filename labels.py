import os
import random
import json

label_dir = os.environ["DATA_DIR"] + "/labels"

def get_label_ids():
    return map(lambda name: int(name.split(".")[0]), os.listdir(label_dir))

def label_file_name(label_id):
    return label_dir + "/" + str(label_id) + ".json"

def save_label(label):
    next_id = max(get_label_ids()) + 1
    with open(label_file_name(next_id), "w") as f:
        json.dump(label, f)

def get_label(label_id):
    with open(label_file_name(label_id), "r") as f:
        return f.read()

def delete_label(label_id):
    os.remove(label_file_name(label_id))