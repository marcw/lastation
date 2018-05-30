package main

import "fmt"
import "net/http"
import "log"
import "time"
import "io/ioutil"
import "encoding/json"
import "github.com/faiface/beep"
import "github.com/faiface/beep/mp3"
import "github.com/faiface/beep/speaker"

type trackinfo struct {
    Trackname string
}

func main() {
	fmt.Println(">>> Playing lastation.fm");
    resp, err := http.Get("https://radio.lastation.fm/listen.mp3")
    if (err != nil) {
        log.Fatalln(err);
	}

    // Decoding mp3 from raw response
	s, format, err := mp3.Decode(resp.Body)
    if err != nil {
        log.Fatalln(err);
	}
    defer s.Close()

    // Open speakers
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan struct{})
    speaker.Play(beep.Seq(s, beep.Callback(func() {
        close(done)
    })))

    // Read track info from JSON
    resp2, err := http.Get("https://lastation.fm/track.json")
    if (err != nil) {
        log.Fatalln(err);
	}
	rawTrackInfo, _ := ioutil.ReadAll(resp2.Body);
    defer resp2.Body.Close();

	t := trackinfo{}
	if err := json.Unmarshal(rawTrackInfo, &t); err != nil {
        log.Println("Error while decoding JSON", err);
	}
	fmt.Println(">>> Current track: ", t.Trackname);

    <-done
}
