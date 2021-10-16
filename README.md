# NhenDownloader

![Go Report Card](https://goreportcard.com/badge/github.com/tianyagk/NhenDownloader) ![Lines of code](https://img.shields.io/tokei/lines/github/tianyagk/NhenDownloader)    


Nentai Downloader for Windows/Linux.

Nhentai下载器，用于批量下载指定漫画。

原项目由Python编写，考虑到并发性能与可移植性使用Golang进行重构。


- [X] 代理模式  Proxy
- [X] 并发下载  Goroutine
- [X] 选择语言  Language
- [X] 下载失败重试  Retry
- [ ] 根据tags 筛选  Filter by Tags



```shell
NhentaiDownloader.exe
Launch Nhentai-Downloader, entry 'help' for more information.
>> help
help                    | Query help
exit                    | Exit NhenDownloader
lang    -language       | Chose from ['chinese', 'japanese', 'english',etc.]
proxy   -host:port      | Like 'http://localhost:7890'
maxTryTimes     -tryNum | Default max try num: 5
recent                  | Show recent popular manga with id
download        -id     | 376478 -> 'https://nhentai.net/g/376478/'

>> lang -english
Setting language: english

>> proxy -http://localhost:8000
Setting proxy: http://localhost:8000

>> recent
0 | [HUACA] DEC20 [English] |
1 | [HUACA] JAN21 [English] |
....
23 | [Yami Books] No Alert [English] [FC] |
24 | [Tanabe Kyou] Love Petit Gate Ch.1-4 [English] {Mistvern + Bigk40k} |

>> download -376453

```





## Python-version

