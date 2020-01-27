import os
import slack
import time
from typing import Union, Dict, Tuple
from asyncio import Future
from slack.web.base_client import SlackResponse
from datetime import datetime, timedelta


def valid_response(response: Union[Future, SlackResponse]):
    assert response["ok"] is True, f"Response Error -> {response}"
    return response


def user_list(client:slack.WebClient) -> Dict:
    response = client.users_list(limit=1000)
    response = valid_response(response)

    users = {}
    for user in response["members"]:

        users[user["id"]] = {"name": user["name"],
                             "display_name": user["profile"]["display_name"],
                             "favicon_img": [user["profile"]["image_24"],
                                             user["profile"]["image_32"],
                                             user["profile"]["image_72"],
                                             user["profile"]["image_192"],
                                             user["profile"]["image_512"]]}
    return users


def channel_list(client:slack.WebClient) -> Dict:
    response = client.channels_list(exclude_archived=1)
    response = valid_response(response)

    channels = {}
    for channel in response["channels"]:
        channels[channel["name"]] = channel["id"]
    return channels


def get_history(client:slack.WebClient, channel_id: str, start: str, end: str):
    response = client.channels_history(channel=channel_id,
                                       latest=start,
                                       oldest=end,
                                       count=1000,
                                       inclusive=0)
    response = valid_response(response)

    messages = []
    for message in response["messages"]:

        if is_thread(message):
            continue

        messages.append({"user_id": message["user"],
                         "text": message["text"],
                         "timestamp": message["ts"]})

    return messages


def is_thread(message):
    keys = message.keys()
    return True if "parent_user_id" in keys else False


def get_history_days() -> Tuple[str, str]:
    today = str(time.mktime(datetime.today().timetuple()))
    yesterday = str(time.mktime((datetime.today() + timedelta(days=-1)).timetuple()))
    return today, yesterday


if __name__ == "__main__":
    client = slack.WebClient(token=os.environ.get("DFAB_BOT"), ssl=False)

    channels = channel_list(client=client)
    users = user_list(client)

    channel_id = channels["daily-logs"]
    start, end = get_history_days()

    response = get_history(client, channel_id, start, end)
    print(response)
