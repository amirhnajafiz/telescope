import requests
from bs4 import BeautifulSoup
import os
from urllib.parse import urljoin

BASE_URL = "https://ftp.itec.aau.at/datasets/DASHDataset2014/BigBuckBunny/6sec/"
OUTPUT_DIR = "bunny_dash"
SELECTED_BITRATES = ["bunny_128256bps", "bunny_217651bps"]

def download_file(url, output_path):
    r = requests.get(url)
    r.raise_for_status()
    with open(output_path, "wb") as f:
        f.write(r.content)

def download_segments_from_folder(folder_name):
    folder_url = urljoin(BASE_URL, folder_name + "/")
    print(f"\n📁 Scanning folder: {folder_url}")

    resp = requests.get(folder_url)
    soup = BeautifulSoup(resp.text, "html.parser")

    output_path = os.path.join(OUTPUT_DIR, folder_name)
    os.makedirs(output_path, exist_ok=True)

    for link in soup.find_all("a"):
        href = link.get("href")
        if href and href.endswith(".m4s"):
            full_url = urljoin(folder_url, href)
            save_path = os.path.join(output_path, href)
            print(f"⬇️ Downloading: {href}")
            download_file(full_url, save_path)

def download_manifest():
    manifest_name = "BigBuckBunny_6s_onDemand_2014_05_09.mpd"
    full_url = urljoin(BASE_URL, manifest_name)
    save_path = os.path.join(OUTPUT_DIR, manifest_name)

    print(f"\n📄 Downloading MPD: {manifest_name}")
    download_file(full_url, save_path)

if __name__ == "__main__":
    os.makedirs(OUTPUT_DIR, exist_ok=True)

    download_manifest()

    for folder in SELECTED_BITRATES:
        download_segments_from_folder(folder)

    print("\n✅ All done!")
