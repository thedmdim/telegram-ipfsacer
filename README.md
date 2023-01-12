# telegram-ipfsacer
A Telegram bot to archive your YouTube videos to IPFS written in GO

![](https://cloudflare-ipfs.com/ipfs/bafybeifk6hallazcdbgimuwie47uzukm56ltljdgqdsx2poh74xrfm37wu)

## Workflow

1. This bot reads post's last line
3. If it contains youtube link, starts download the video via [kkdai youtube downloader](https://github.com/kkdai/youtube)
4. Writes a fideo to /storage/<video_id>.mp4 in [MFS](https://docs.ipfs.tech/concepts/file-systems/#add-a-file-to-mfs) of your IPFS node
5. Edits post with link to added video
6. Creates IPNS link for whole storage to them all (for link to be static, use your own key)
