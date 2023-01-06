package main

import (
	"context"
	"flag"
	"fmt"
	"telegram-ipfsacer/telegram"
	"telegram-ipfsacer/video"
	"log"
	"strings"
	"time"

	"telegram-ipfsacer/ipfs"
)

var (
	tgBotApiToken = flag.String("token", "", "specify the token for Telegram bot API")
	tgChannel   = flag.String("channel", "", "specify the target Telegram cahnnel, example @channelname")
	ipfsNodeUrl = flag.String("url", "localhost:5001", "specify the IPFS node URL, default localhost:5001")
	mfsDirName = flag.String("storage", "storage", "specify the target Telegram cahnnel (optional)")
	keyPath = flag.String("key", "", "specify the path of IPNS key (optional)")
)

const tgBotApiUrl string = "https://api.telegram.org/bot"

var vid video.Client = video.Client{}
var tg *telegram.Client
var c *ipfs.Client

func main(){
	var err error

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintln(w, "Usage of go-ipfs-arch")
		flag.PrintDefaults()
	}
	flag.Parse()

	// some checks before start
	if *tgBotApiToken == "" || *tgChannel == "" {
		fmt.Println(tgBotApiToken, *tgChannel)
		log.Fatalln("Set TGBOT_TOKEN and TG_CHANNEL env variables")
	}

	tg = telegram.NewClient(tgBotApiUrl, *tgBotApiToken)

	c, err = ipfs.NewIPFSClient(*ipfsNodeUrl, *mfsDirName, *keyPath)
	if err != nil {
		log.Fatalln(fmt.Errorf("cannot create IPFS client: %w", err))
	}
	
	// event loop

	for ;; {
		updates, err := tg.Updates()
		time.Sleep(time.Second)

		if err != nil {
			log.Println(err)
		}
		for _, update := range updates {
			log.Println("Process update")
			processUpdate(&update)
			tg.Offset = update.UpdateId + 1
			
			
		}
	}
}


func processUpdate(result *telegram.Result){
	
	if result == nil {
		return
	}

	defer func() {
        if r := recover(); r != nil {
            log.Println("Recovered: %", r)
        }
    }()


	lines := strings.Split(result.Post.Text, "\n")
	ytLink := lines[len(lines)-1]
	v := strings.Split(ytLink, "?v=")
	videoID := v[len(v)-1]
	if len(videoID) != 11 {
		log.Printf("cannot parse video id in %s", ytLink)
		return
	}

	log.Printf("try get stream of %s", videoID)
	vid, err := vid.Stream(videoID)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("adding stream to ipfs %s", videoID)
	cid, err := c.AddVideo(context.TODO(), vid)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("added %s -> %s", vid.Filename, cid)

	lines[len(lines)-1] = fmt.Sprintf("[youtube](https://www.youtube.com/watch?v=%s) | [ipfs](https://ipfs.io/ipfs/%s/%s)", videoID, cid, vid.Filename)

	message := telegram.EditedPost{
		Id: result.Post.Id,
		ChatId: *tgChannel,
		Text: strings.Join(lines, "\n"),
		ParseMode: "Markdown",
	}

	tg.EditMessage(message)

	go func() {
		_, err = c.Sh.PublishWithDetails(cid, c.KeyName, 0, 0, false)
		if err != nil {
			log.Println(err)
		}
		log.Println("IPNS updated with %s", cid)
	}()

}

