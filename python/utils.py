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

# Define a custom formatting function
def human_readable_formatter(x, pos):
    """
    Custom formatting function to convert tick values to human-readable names.
    For example, 1000 will be converted to '1K', and 1500000 will be converted to '1.5M'.
    """
    if x >= 1e6:
        return f"{x/1e6:.1f}M"
    elif x >= 1e3:
        return f"{x/1e3:.0f}K"
    else:
        return f"{x:.0f}"