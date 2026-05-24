import os

icons_dir = 'dist/optimized'
if not os.path.exists(icons_dir):
    print(f"Directory {icons_dir} not found.")
    exit(1)

icons = sorted([f for f in os.listdir(icons_dir) if f.endswith('.svg')])

def format_name(filename):
    return filename.replace('.svg', '').replace('-', ' ').title()

gallery = "## Icon Gallery\n\n"
gallery += "| Icon | Name | Icon | Name | Icon | Name |\n"
gallery += "| :---: | :--- | :---: | :--- | :---: | :--- |\n"

for i in range(0, len(icons), 3):
    row = "| "
    for j in range(3):
        if i + j < len(icons):
            filename = icons[i+j]
            name = format_name(filename)
            # Use raw github user content or relative path. Relative path works best for GitHub READMEs.
            row += f"<img src='./dist/optimized/{filename}' width='48' height='48' /> | {name} | "
        else:
            row += " | | "
    gallery += row + "\n"

with open('README.md', 'r') as f:
    readme = f.read()

marker = '## Development Workflow'
if marker in readme and '## Icon Gallery' not in readme:
    new_readme = readme.replace(marker, gallery + '\n' + marker)
    with open('README.md', 'w') as f:
        f.write(new_readme)
    print("Gallery generated and inserted into README.md")
elif '## Icon Gallery' in readme:
    # Update existing gallery if it exists (very basic replacement)
    import re
    pattern = re.compile(r'## Icon Gallery.*?## Development Workflow', re.DOTALL)
    new_readme = pattern.sub(gallery + '\n## Development Workflow', readme)
    with open('README.md', 'w') as f:
        f.write(new_readme)
    print("Gallery updated in README.md")
else:
    print("Marker not found in README.md")
