import json

def read_data(filename):
    f = open(filename, 'r')
    isHeader = True
    header = {}
    contents = []
    while True:
        line = f.readline()
        if not line:
            break
        line = line.strip()
        if isHeader:
            header = json.loads(line)
            isHeader = False
        else:
            contents.append(json.loads(line))
    return header, contents