from htmllaundry import strip_markup
import re


def clean_content(content):
    content = strip_markup(content)
    content = re.sub(r'\n\n+', '\n\n', content)
    return content


if __name__ == '__main__':
    print(clean_content("Hallo welt\n\n\n\naSdasd"))
