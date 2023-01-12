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
	tgBotApiToken = flag.String("token", "", "the token for Telegram bot API")
	tgChannel   = flag.String("channel", "", "the target Telegram cahnnel, example: @channelname")
	ipfsNodeUrl = flag.String("url", "localhost:5001", "the IPFS node URL")
	mfsDirName = flag.String("storage", "storage", "specify the dir name for all videos (optional)")
	keyPath = flag.String("key", "", "specify the path of IPNS key (optional)")
	ipfsGateway = flag.String("ipfs-gateway", "ipfs.io", "specify public IPFS gateway (optional)")
	ipnsUpdate = flag.Int("ipns-update", 24, "specify period in hours when IPNS will be updated")
)

const tgBotApiUrl string = "https://api.telegram.org/bot"

var vid video.Client = video.Client{}
var tg *telegram.Client
var c *ipfs.Client

func main(){
	var err error

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintln(w, "Usage of TELEGRAM-IPFSACER")
		flag.PrintDefaults()
	}
	flag.Parse()

	// some checks before start
	if *tgBotApiToken == "" || *tgChannel == "" {
		fmt.Println(tgBotApiToken, *tgChannel)
		log.Fatalln("Please provide Telegram token, @channelname")
	}

	tg = telegram.NewClient(tgBotApiUrl, *tgBotApiToken)

	c, err = ipfs.NewIPFSClient(*ipfsNodeUrl, *mfsDirName, *keyPath)
	if err != nil {
		log.Fatalln(fmt.Errorf("cannot create IPFS client: %w", err))
	}

	// IPNS update every 24 hours

	go func() {
		for {
			Publish()
			time.Sleep(time.Hour * time.Duration(*ipnsUpdate))
		}
	}()
	
	// event loop

	for ;; {
		updates, err := tg.Updates()
		time.Sleep(time.Minute)

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
	log.Printf("got youtube link %s", ytLink)

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

	lines[len(lines)-1] = fmt.Sprintf("[youtube](https://www.youtube.com/watch?v=%s) | [ipfs](https://%s/ipfs/%s)", videoID, *ipfsGateway, cid)

	message := telegram.EditedPost{
		Id: result.Post.Id,
		ChatId: *tgChannel,
		Text: strings.Join(lines, "\n"),
		ParseMode: "Markdown",
	}

	tg.EditMessage(message)

	go Publish()
}


func Publish() {
	r, _ := c.PublishCurrent(context.Background())
	log.Printf("IPNS published %s -> %s", r.Value, r.Name)
}