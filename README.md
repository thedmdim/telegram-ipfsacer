# telegram-ipfsacer
A Telegram bot to archive your YouTube videos to IPFS written in GO

![](https://raw.githubusercontent.com/thedmdim/telegram-ipfs-archiver/master/example.jpg)

## Workflow

1. This bot reads post's last line
3. If it contains youtube link, starts download the video via youtube-dl
4. Adds a folder with videos in ipfs
5. Create IPNS link (for link to be static, use your own key)
