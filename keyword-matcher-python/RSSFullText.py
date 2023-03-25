import time

import requests

import Config
from model.RSSFullTextResponse import RSSFullTextResponse


class FullTextRss:
    def __init__(self):
        while True:
            try:
                requests.get(Config.FULLTEXTRSS_SERVER + "/extract.php")
                break
            except:
                print("Waiting for", Config.FULLTEXTRSS_SERVER, "to come up...")
                time.sleep(5)


    def retrieve_full_text(self, url) -> RSSFullTextResponse:
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
    ftr = FullTextRss()
    print(ftr.retrieve_full_text(
        "https://www.welt.de/kultur/article239951133/Documenta-Generaldirektorin-Schormann-legt-nach-Antisemitismus-Eklat-Amt-nieder.html"))
