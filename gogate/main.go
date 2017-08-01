package gogate

import (
	"net/http"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"encoding/base64"
	"encoding/json"
	"bytes"
	"errors"
	"strings"
)

func init() {
	http.HandleFunc("/", handler)
}

var speechURL = "https://speech.googleapis.com/v1beta1/speech:syncrecognize?key=" + os.Getenv("SPEECH_API_KEY")

const (
	minAccuracy = 0.7
	passwd = "who is john galt?"

	welcomeMsg = `<?xml version="1.0" encoding="UTF-8"?>
	<Response>
		<Say>Hello there, what's the password?'</Say>
		<Record timeout="3" />
	</Response>`

	okPasswd = `<?xml version="1.0" encoding="UTF-8"?>
	<Response>
		<Say>That is the correct password!</Say>
	</Response>`

	badPasswd = `<?xml version="1.0" encoding="UTF-8"?>
	<Response>
		<Say>Go away you evil person!</Say>
	</Response>`

	noTranscription = `<?xml version="1.0" encoding="UTF-8"?>
	<Response>
		<Say>Sorry, I didn't quite understand that. Please try again.  What's the password'?'</Say>
		<Record timeout="3" />
	</Response>`
)


func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	/* err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for k, v := range r.PostForm {
		log.Infof(c,"%s: %v", k, v)
	} */

	w.Header().Set("Concent-Type", "text/xml")
	rec := r.FormValue("RecordingUrl")
	if rec == "" {
		fmt.Fprintf(w, welcomeMsg)
		return
	}

	text, err  := transcribe(c, rec)
	if err != nil {
		if err == ErrNoResults || err == ErrInaccurate {
			fmt.Fprintf(w, noTranscription)
			return
		}
		http.Error(w, "We couldn't transcribe", http.StatusInternalServerError)
		log.Errorf(c, "Could not transcribe: %v", err)
		return
	}

	i := strings.ToLower(text)

	if  i == passwd {
		fmt.Fprintf(w, okPasswd)
		return
	} else {
		fmt.Fprintf(w, badPasswd)
		return
	}

}

func transcribe(c context.Context, url string) (string, error) {
	b, err := fetchAudio(c, url)
	if err != nil {
		return "", err
	}

	return fetchTranscription(c, b, minAccuracy)
}

func fetchAudio(c context.Context, url string) ([]byte, error) {
	client := urlfetch.Client(c)

	res, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Could not fetch %v: %v", url, err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Fetched with status: %s", res.Status)
	}

	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read response: %v", err)
	}

	return b, nil
}

type speechReq struct {
	Config struct {
		Encoding string `json:"encoding"`
		SampleRate int `json:"sampleRate"`
	} `json:"config"`
	Audio struct {
		Content string `json:"content"`
	} `json:"audio"`
}

var ErrNoResults = errors.New("No transcription results found")
var ErrInaccurate = errors.New("Transcription is too inaccurate")

func fetchTranscription(c context.Context, b []byte, a float64) (string, error) {
	var req speechReq

	req.Config.Encoding = "LINEAR16"
	req.Config.SampleRate = 8000
	req.Audio.Content = base64.StdEncoding.EncodeToString(b)

	j, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("Could not encode speech request: %v", err)
	}

	res, err := urlfetch.Client(c).Post(speechURL, "application/json", bytes.NewReader(j))
	if err != nil {
		return "", fmt.Errorf("Could not transcribe: %v", err)
	}

	var data struct {
		Error struct {
			Code int
			Message string
			Status string
		}
		Results []struct {
			Alternatives []struct {
				Transcript string
				Confidence float64
			}
		}
	}

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", fmt.Errorf("Could not decode speech response: %v", err)
	}
	if data.Error.Code != 0 {
		return "", fmt.Errorf("Speech API error: %d %s %s", data.Error.Code, data.Error.Status, data.Error.Message)
	}
	log.Infof(c, "Data received: %+v", data)
	if len(data.Results) == 0 || len(data.Results[0].Alternatives) == 0 {
		return "", ErrNoResults
	}
	if data.Results[0].Alternatives[0].Confidence < a {
		return "", ErrInaccurate
	}

	return data.Results[0].Alternatives[0].Transcript, nil
}