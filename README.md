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



**使用help指令获取功能详情**

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
```

**使用lang指令设置默认语言**

```go
>> lang -english
Setting language: english
```

**使用recent指令查看最近流行**

```go
>> recent
0 | [Puranpuman] Muteikou-san ni wa Kiken ga Ippai! [Chinese] [羅莎莉亞漢化] [Digital] |
  sole male, ahegao, bbw, big ass, big breasts, big nipples, blowjob, blowjob face, collar, exhibitionism, hair buns, hotpants, huge breasts, milf, sole female, sweating, mosaic censorship,
1 | [Takashi] Koi | 恋 (COMIC X-EROS #87) [Chinese] [暴碧汉化组] [Digital] |
  big breasts, blowjob, cheating, deepthroat, milf, nakadashi, netorare, rape, sole female, group,
2 | [Jewelry Box (Tamaki Nao)] Chikako-san to Issho ni! 1 | 和千伽子小姐一起! 1 [Chinese] [橄榄汉化组] [Incomplete] |
  doujinshi, glasses, incomplete,
3 | [Uno Ryoku]关于我性转成魅魔以后建立幸福家庭这件事1 [Chinese] [Aelit个人汉化] |
  gender bender, transformation, corruption, crotch tattoo, demon girl, horns, lolicon, monster girl, tail,
......
20 | [Small Marron (アサクラククリ)] 性交秘話〜彼氏持ちの私が年下のオタクに堕とされるまで〜[Chinese]【枫原万叶汉化】 |
  bbm, big penis, condom, glasses, sole male, beauty mark, big breasts,
```

**使用download -id指令下载指定**

```go
>> download -376453

- (C78)_[Jouji_Mujoh_(Shinozuka_George)]_Ura_Mugi_(K-ON!)_[Chinese]_[final個人漢化]_[Decensored]  in download queue.
- download finish (C78)_[Jouji_Mujoh_(Shinozuka_George)]_Ura_Mugi_(K-ON!)_[Chinese]_[final個人漢化]_[Decensored] .
```



## Python-version

