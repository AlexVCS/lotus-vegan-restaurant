"""Extract the restaurant's accordion menu into a small, safe JS data file."""

import html
import json
import re
from html.parser import HTMLParser
from pathlib import Path


class MenuParser(HTMLParser):
    def __init__(self):
        super().__init__(convert_charrefs=True)
        self.in_item = False
        self.item_depth = 0
        self.mode = None
        self.ignore_depth = 0
        self.title = []
        self.content = []
        self.items = []

    def handle_starttag(self, tag, attrs):
        attrs = dict(attrs)
        classes = attrs.get("class", "").split()
        if tag in {"script", "style"}:
            self.ignore_depth += 1
            return
        if self.ignore_depth:
            return
        if tag == "div" and "accordion__item" in classes:
            self.in_item = True
            self.item_depth = 1
            self.mode = None
            self.title = []
            self.content = []
            return
        if self.in_item and tag == "div":
            self.item_depth += 1
            if "accordion__title" in classes:
                self.mode = "title"
            elif "accordion__content" in classes:
                self.mode = "content"
        if self.in_item and self.mode == "content" and tag in {"br", "li", "p"}:
            self.content.append("\n")

    def handle_endtag(self, tag):
        if tag in {"script", "style"} and self.ignore_depth:
            self.ignore_depth -= 1
            return
        if self.ignore_depth:
            return
        if self.in_item and tag == "div":
            self.item_depth -= 1
            if self.item_depth == 0:
                title = clean("".join(self.title)).strip(" :")
                lines = [clean(line) for line in "".join(self.content).splitlines()]
                lines = [line for line in lines if line]
                if title and lines:
                    self.items.append({"title": title.title(), "lines": lines})
                self.in_item = False
                self.mode = None

    def handle_data(self, data):
        if self.ignore_depth or not self.in_item:
            return
        if self.mode == "title":
            self.title.append(data)
        elif self.mode == "content":
            self.content.append(data)


def clean(value):
    value = html.unescape(value).replace("\u200b", " ").replace("\xa0", " ")
    return re.sub(r"\s+", " ", value).strip()


source = Path("/tmp/lotus-original.html")
target = Path("public/menu-data.js")
parser = MenuParser()
parser.feed(source.read_text(encoding="utf-8", errors="ignore"))
target.write_text("window.LOTUS_MENU = " + json.dumps(parser.items, ensure_ascii=False) + ";\n", encoding="utf-8")
print(f"wrote {len(parser.items)} categories to {target}")
