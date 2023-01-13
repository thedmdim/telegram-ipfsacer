# telegram-ipfsacer
В данной ветке представлен код, который запущен на моём Orange Pi Zero, обслуживающего канал [MOV3371](https://t.me/mov3371)

![](https://cloudflare-ipfs.com/ipfs/bafkreihzbgxq2q7fpvfvuc2o6li33jzd77vil4qhaswzvn4r7mpatelydy)

[Архив](https://cloudflare-ipfs.com/ipns/k51qzi5uqu5di6sixp2l59em0ajrgzakb7p52s8qdgq5j1dolz4aubvdx869a0/) видео c канала

## Компиляция
- Linux
```bash
env GOOS=linux GOARCH=arm GOARM=7 go build -trimpath -ldflags="-s -w" .
```
- Windows
```cmd
set GOOS=linux; set GOARCH=arm; set GOARM=7; go build -trimpath -ldflags="-s -w" .
```