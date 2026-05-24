import os
import re
import subprocess
import xml.etree.ElementTree as ET

# Configuration
INKSCAPE_PATH = "/Applications/Inkscape.app/Contents/MacOS/inkscape"
SRC_FILE = "src/core-architecture-shapes-master.svg"
EXPORT_FILE = "src/core-architecture-shapes-export.svg"
DIST_DIR = "dist"

# Namespaces
NS = {
    "svg": "http://www.w3.org/2000/svg",
    "inkscape": "http://www.inkscape.org/namespaces/inkscape",
    "sodipodi": "http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd",
}


def slugify(text):
    text = text.lower()
    text = re.sub(r"[^a-z0-9]+", "-", text)
    return text.strip("-")


def export_pages():
    if not os.path.exists(DIST_DIR):
        os.makedirs(DIST_DIR)

    # Use the export file if it exists (manual flattening), otherwise use master
    input_file = EXPORT_FILE if os.path.exists(EXPORT_FILE) else SRC_FILE
    print(f"Using input file: {input_file}")

    tree = ET.parse(input_file)
    root = tree.getroot()

    namedview = root.find(".//sodipodi:namedview", NS)
    if namedview is None:
        print("Could not find sodipodi:namedview")
        return

    pages = namedview.findall(".//inkscape:page", NS)
    print(f"Found {len(pages)} pages.")

    for i, page in enumerate(pages):
        label = page.get("{http://www.inkscape.org/namespaces/inkscape}label")
        if not label:
            label = f"page-{i}"

        filename = slugify(label) + ".svg"
        dest_path = os.path.join(DIST_DIR, filename)

        print(f"Exporting {label} (Page {i + 1}) to {dest_path}...")

        # Inkscape actions:
        # 1. select-all:all (selects all objects on the current page context)
        # 2. object-to-path (ensures shapes like circles/rects are paths)
        # 3. stroke-to-path (CRITICAL: converts the 2px strokes into shapes so they can be unioned)
        # 4. selection-union (merges everything into a single production-ready path)
        # 5. export-do (executes the export)
        actions = (
            "select-all:all;object-to-path;stroke-to-path;selection-union;export-do"
        )

        cmd = [
            INKSCAPE_PATH,
            input_file,
            f"--export-page={i + 1}",
            "--export-type=svg",
            "--export-plain-svg",
            f"--actions={actions}",
            "--export-filename=" + dest_path,
        ]

        try:
            subprocess.run(cmd, check=True, capture_output=True)
        except subprocess.CalledProcessError as _:
            # If complex actions fail, fallback to a simpler export to ensure we get something
            print(
                f"Warning: Complex actions failed for {label}, retrying basic export..."
            )
            cmd_basic = [
                INKSCAPE_PATH,
                input_file,
                f"--export-page={i + 1}",
                "--export-type=svg",
                "--export-plain-svg",
                "--actions=select-all:all;object-to-path;export-do",
                "--export-filename=" + dest_path,
            ]
            subprocess.run(cmd_basic, check=True, capture_output=True)


if __name__ == "__main__":
    export_pages()
