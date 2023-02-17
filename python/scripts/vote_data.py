import io
import zipfile

import requests

res = requests.get("https://data.cec.gov.tw/選舉資料庫/votedata.zip")

with zipfile.ZipFile(io.BytesIO(res.content), "r", metadata_encoding="big5") as zip_ref:
    zip_ref.extractall("../data")
