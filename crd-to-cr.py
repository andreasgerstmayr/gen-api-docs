#!/usr/bin/env python
"""
This tool reads a Kubernetes Custom Resource Definition (CRD) and outputs a commented, full example Custom Resource (CR).
"""
import sys
import argparse
import yaml

COMMENT_PADDING = 40


def get_default_value(obj):
    type_ = obj.get("type")
    if type_ == "string":
        return '""'
    elif type_ == "boolean":
        return "false"
    elif type_ == "integer":
        return "0"
    elif obj.get("x-kubernetes-int-or-string", False):
        return "1Gi"
    elif obj.get("x-kubernetes-preserve-unknown-fields", False):
        return "{}"
    else:
        return ""


def format_comment(obj):
    return obj.get("description", "").replace("\n", " ")


def write_line(out, line, comment=""):
    padding = COMMENT_PADDING - len(line)
    if comment:
        out.write(f"{line}{' '*padding} # {comment}\n")
    else:
        out.write(f"{line}\n")


def render(out, obj, level, isList=False):
    type_ = obj["type"]
    if type_ == "object" and "properties" in obj:  # struct
        for i, (name, member) in enumerate(obj["properties"].items()):
            indent = f"{'  '*(level-1)}- " if i == 0 and isList else f"{'  '*level}"
            comment = format_comment(member)
            value = get_default_value(member)

            if value:
                write_line(out, f"{indent}{name}: {value}", comment)
            else:
                write_line(out, f"{indent}{name}:", comment)
                render(out, member, level + 1)
    elif type_ == "object" and "additionalProperties" in obj:  # map
        indent = f"{'  '*level}"
        comment = format_comment(obj["additionalProperties"])
        value = get_default_value(obj["additionalProperties"])
        if value:
            write_line(out, f'{indent}"key": {value}', comment)
        else:
            write_line(out, f'{indent}"key":')
            render(out, obj["additionalProperties"], level + 1)
    elif type_ == "array":
        indent = f"{'  '*(level-1)}"
        comment = format_comment(obj["items"])
        value = get_default_value(obj["items"])
        if value:
            write_line(out, f"{indent}- {value}", comment)
        else:
            render(out, obj["items"], level, isList=True)


def run(crdFile, out):
    crd = yaml.safe_load(crdFile)
    spec = crd["spec"]
    group = spec["group"]
    kind = spec["names"]["kind"]

    version = spec["versions"][0]
    version_name = version["name"]
    obj = version["schema"]["openAPIV3Schema"]

    apiVersionProp = obj["properties"].pop("apiVersion")
    kindProp = obj["properties"].pop("kind")
    obj["properties"].pop("metadata")

    # use real values where possible
    write_line(
        out, f"apiVersion: {group}/{version_name}", format_comment(apiVersionProp)
    )
    write_line(out, f"kind: {kind}", format_comment(kindProp))
    write_line(out, f"metadata:")
    write_line(out, f"  name: example")
    render(out, obj, 0)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="This tool reads a Kubernetes Custom Resource Definition (CRD) and outputs a commented, full example Custom Resource (CR)."
    )
    parser.add_argument(
        "crd", nargs="?", type=argparse.FileType("r"), default=sys.stdin
    )
    parser.add_argument(
        "out", nargs="?", type=argparse.FileType("w"), default=sys.stdout
    )
    args = parser.parse_args()
    run(args.crd, args.out)
