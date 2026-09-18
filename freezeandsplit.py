import pathlib
import sys

import ruamel.yaml

yaml = ruamel.yaml.YAML()
yaml.indent(sequence=4, offset=2)


def save(obj, path):
    print(f"Writing {path}")
    with path.open(mode="w", encoding="utf-8") as outfile:
        yaml.dump(obj, outfile)


def main():
    filepath = pathlib.Path(sys.argv[1])
    with filepath.open(encoding="utf-8") as yamlfile:
        data = yaml.load(yamlfile)

    outdir = filepath.parent

    common = data.get(".common")
    save(common, outdir / "common.yaml")

    for key, obj in data["image_types"].items():
        save(obj, outdir / (key + ".yaml"))


if __name__ == "__main__":
    main()
