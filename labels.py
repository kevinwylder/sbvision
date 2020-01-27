import os

label_dir = os.environ["DATA_DIR"] + "/labels"

def is_good_label(label): 
    try:
        return len(label["input"]["data"]) == label["input"]["width"] * label["input"]["height"] * 4 and \
           len(label["output"]["rotation"]) == 4 and \
               "isSkateboard" in label["output"]
    except:
        return False

def save_label(label):
    filename = str(label["input"]["width"]) + "x" + str(label["input"]["height"]) + "-" + ",".join(map(str, label["output"]["rotation"])) + ".json"
    with open(label_dir + "/" + filename, "w") as f:
        json.dump(label, f)