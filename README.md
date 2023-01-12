# telegram-ipfsacer
A Telegram bot to archive your YouTube videos to IPFS written in GO

![](https://cloudflare-ipfs.com/ipfs/bafybeifk6hallazcdbgimuwie47uzukm56ltljdgqdsx2poh74xrfm37wu)

## Workflow

1. This bot reads post's last line
3. If it contains youtube link, starts download the video via [kkdai youtube downloader](https://github.com/kkdai/youtube)
4. Writes a fideo to /storage/<video_id>.mp4 in [MFS](https://docs.ipfs.tech/concepts/file-systems/#add-a-file-to-mfs) of your IPFS node
5. Edits post with link to added video
6. Creates IPNS link ([example](https://cloudflare-ipfs.com/ipns/k51qzi5uqu5di6sixp2l59em0ajrgzakb7p52s8qdgq5j1dolz4aubvdx869a0/)) for whole storage to them all (for link to be static, use your own key)

## Host project is [MOV3371](https://t.me/mov3371)
_a place where forgotten videos find a second life_
