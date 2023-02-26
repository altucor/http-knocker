import re

bad_symbols = [
    "-", "/", "."
]

def sanitize_string(input):
    for item in bad_symbols:
        input = input.replace(item, "_")
    input = input.replace("+", "Plus")
    return input

def process_line(line):
    # print(f"Processing line {line}")
    if line.startswith("#"):
        return
    matches = re.findall("^(.+)	(\d+)", line)
    text_name = sanitize_string(matches[0][0])
    print(f"{text_name.upper()}: \"{matches[0][0]}\",")
    # print(matches)

def main():
    with open("/etc/protocols") as f:
        for line in f:
            # if not line:
            #     return
            process_line(line)



if __name__ == "__main__":
    main()