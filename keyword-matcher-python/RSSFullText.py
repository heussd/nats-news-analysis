import requests

import Config
from model.RSSFullTextResponse import RSSFullTextResponse


def retrieve_full_text(url) -> RSSFullTextResponse:
    response = requests.get(Config.FULLTEXTRSS_SERVER + "/extract.php",
                            params={
                                "url": url
                            })

    if response.text == "Invalid URL supplied":
        raise Exception("Invalid URL supplied", url)

    try:
        json = response.json()
        response = RSSFullTextResponse.parse_obj(json)
    except:
        raise Exception("Failed to parse JSON:", response.content)
    return response


if __name__ == "__main__":
    print(retrieve_full_text(
        "https://www.welt.de/kultur/article239951133/Documenta-Generaldirektorin-Schormann-legt-nach-Antisemitismus-Eklat-Amt-nieder.html"))
