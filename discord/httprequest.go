package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

var MAX_REQUEST_TIMES int = 30

func (myContext *MyContext) docomoAPIrequest(m *discordgo.MessageCreate) ([]byte, error) {
	tUrl := myContext.DocomoAPI.Url + "?APIKEY=" + myContext.DocomoAPI.Token

	var apiSetting_json []byte
	var err error
	apiSetting, ok := myContext.VoiceSetting[m.Author.ID]
	if !ok {
		apiSetting_json, err = ioutil.ReadFile("defaultVoice.json")
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(apiSetting_json, &apiSetting); err != nil {
			log.Fatal(err)
		}
	}

	apiSetting.TextData = m.Author.Username + "::" + m.Content

	apiSetting_json, err = json.Marshal(apiSetting)
	if err != nil {
		log.Fatal(err)
	}
	count := 0
REQUEST:
	resp, err := http.Post(tUrl, "application/json", bytes.NewBuffer(apiSetting_json))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	af := getAudio(body)
	if len(af.data) == 0 && count < MAX_REQUEST_TIMES {
		time.Sleep(300 * time.Millisecond)
		count++
		goto REQUEST
	}
	if count == MAX_REQUEST_TIMES {
		return nil, errors.New("http request failed")
	}
	return af.data, nil
}
func byteslice2int16slice(bs []byte) []int16 {
	length := len(bs)
	i16 := make([]int16, length)
	for i := 0; i < length; i += 2 {
		i16[i] = int16(bs[i]) + int16(bs[i+1])*256
	}
	return i16
}
func littleendian4(d []byte) int64 {
	a := int64(d[0])
	b := int64(d[1])
	c := int64(d[2])
	e := int64(d[3])
	return a + b*256 + c*65536 + e*16777216
}
func littleendian2(d []byte) int {
	a := int(d[0])
	b := int(d[1])
	return a + b*256
}

func getAudio(raw []byte) AudioFormat {
	af := AudioFormat{}
	af.headchunkId = string(raw[0:4])
	if af.headchunkId != "RIFF" {
		return AudioFormat{}
	}
	af.headchunkSize = littleendian4(raw[4:8]) - 36
	af.formType = string(raw[8:12])

	af.fmtchunkID = string(raw[12:16])
	af.fmtchunkSize = littleendian4(raw[16:20])
	af.waveFormatType = littleendian2(raw[20:22])
	af.channel = littleendian2(raw[22:24])
	af.samplePerSec = littleendian4(raw[24:28])
	af.bytePerSec = littleendian2(raw[28:30])
	af.blockSize = littleendian2(raw[32:34])
	af.bitsPerSample = littleendian2(raw[34:36])

	af.datachunkID = string(raw[36:40])
	af.datachunkSize = littleendian4(raw[40:44])
	af.data = raw[44:]
	return af
}
