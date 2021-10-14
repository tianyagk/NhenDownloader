import os
import argparse
import json
import multiprocessing
from multiprocessing import Pool

import requests
from bs4 import BeautifulSoup

# parser = argparse.ArgumentParser(description='nhentai')
# parser.add_argument("--id", help="nhentai manka id")


def process_bar(percent, start_str='', end_str='', total_length=0):
    """
    进度条打印函数
    """
    bar = ''.join(["\033[31m%s\033[0m"%'   '] * int(percent * total_length)) + ''
    bar = '\r' + start_str + bar.ljust(total_length) + ' {:0>4.1f}%|'.format(percent*100) + end_str
    print(bar, end='', flush=True)


def spider(manka_id, proxies):
    """
    Nhentai 内容抓取函数
    :param
        manka_id: (int) 漫画编号
    """
    home_url = f"https://nhentai.net/g/{manka_id}/"
    r = requests.get(home_url)
    soup = BeautifulSoup(r.text, 'html.parser')

    galleries = soup.head.find(attrs={"itemprop":"image"})["content"].split("/")[-2]

    content = soup.find(id="content")
    info_block = content.find(id="info-block")

    title = info_block.find(class_="title").find(class_="pretty").text.replace(" ", "_")

    tags = info_block.find(name="section", id="tags")
    # pages = len(content.find(class_="thumbs").find_all())

    pages = int(tags.find_all(class_="tag-container field-name")[-2].find(name="span", class_="name").text)
    # datetime = tags[5].time["datetime"]

    # Ready for download
    path = f"./download/{title}"
    if not os.path.exists(path):
        os.mkdir(path)
    
    print(f"Start download {title}")
    
    try:
        for i in range(1, pages):
            file_path = "/".join((path, f"{i}.jpg"))

            if not os.path.exists(file_path):
                img = requests.get(f"https://i.nhentai.net/galleries/{galleries}/{i}.jpg", proxies=proxies)
                with open(file_path, "wb") as f:
                    f.write(img.content)
                del img
    except Exception as e:
        print(e)

        # process_bar(i/pages, start_str=f"{manka_id}: ", end_str=f"{pages} Pages", total_length=pages)

    print(f"Finish {title}")


if __name__ == "__main__":
    # args = parser.parse_args()

    proxies={'http':'127.0.0.1:7890', 'https':'127.0.0.1:7890'}

    ids = []
    with open("./manga.json", "r") as f:
        ids = json.loads(f.read())["ids"]

    pool =  Pool(multiprocessing.cpu_count())
    for manka_id in ids:
        pool.apply_async(spider, (manka_id, proxies, ))

    pool.close()
    pool.join()