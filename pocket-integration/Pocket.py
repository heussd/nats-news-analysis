from pocket import Pocket, PocketException
import Config

pocket = Pocket(
    consumer_key=Config.POCKET_CONSUMER_KEY,
    access_token=Config.POCKET_ACCESS_TOKEN
)


def add_to_pocket(url, tags):
    try:
        response = pocket.add(url=url, tags=([tags] + ["spoiler"]))
        item_id = response.get("item").get("item_id")
        assert item_id is not None
        print("Added to Pocket", url)
    except PocketException:
        print("ERROR: Could not add", url)


if __name__ == "__main__":
    add_to_pocket('https://www.tagesschau.de/newsticker/liveblog-coronavirus-sonntag-361.html', ['spoiler', 'mycustomtag'])
