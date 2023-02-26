import json
import logging
from dataclasses import dataclass
from typing import List, Optional

import requests
from prefect import flow, task


@dataclass
class Legislator:
    areaName: str
    committee: str
    degree: str
    ename: str
    experience: str
    leaveDate: str
    leaveFlag: str
    leaveReason: str
    name: str
    onboardDate: str
    party: str
    partyGroup: str
    picUrl: str
    sex: str
    term: str
    addr: Optional[str] = None
    fax: Optional[str] = None
    tel: Optional[str] = None


def get(url: str) -> bytes:
    resp = requests.get(url)
    resp.raise_for_status()
    return resp.content


def parse_legislator(body: bytes) -> List[Legislator]:
    result = json.loads(body)
    legislator_list = [Legislator(**data) for data in result["dataList"]]
    return legislator_list


@task
def get_current_legislator_info() -> List[Legislator]:
    url = "https://data.ly.gov.tw/odw/ID9Action.action?fileType=json"
    logging.info(f"get_current_legislator_info: {url}")
    return parse_legislator(get(url))


@task
def get_history_legislator_info() -> List[Legislator]:
    url = "https://data.ly.gov.tw/odw/ID16Action.action?fileType=json"
    logging.info(f"get_history_legislator_info: {url}")
    return parse_legislator(get(url))


@task
def to_file(file_path: str, records: List[Legislator]):
    with open(file_path, "a") as f:
        for record in records:
            data = json.dumps(record.__dict__, ensure_ascii=False)
            f.write(f"{data}\n")


@flow
def crawl_legislator():
    records = get_current_legislator_info() + get_history_legislator_info()
    to_file("data/legislators.jsonl", records)


if __name__ == "__main__":
    crawl_legislator()  # type: ignore
