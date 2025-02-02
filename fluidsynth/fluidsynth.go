package fluidsynth

import (
        "encoding/json"
        "html/template"
        "log"
	"net"
	"net/http"
	"sync"
	"time"

)

var templates *template.Template
var Mutex sync.Mutex

func Include(t *template.Template) {
	templates = t
}

func init() {
	http.HandleFunc("POST /fluidsynth/{host}/{port}/setting/{setting}", postSettingHandler)
	http.HandleFunc("POST /fluidsynth/{host}/{port}/midicc/{channel}/{control}", postMidiccHandler)
	http.HandleFunc("POST /fluidsynth/{host}/{port}/selectfont/{channel}", postSelectFontHandler)
	http.HandleFunc("POST /fluidsynth/{host}/{port}/setbankpreset/{channel}", postSetBankPresetHandler)
	http.HandleFunc("POST /fluidsynth/{host}/{port}/setvolume/{channel}", postSetVolumeHandler)
}

func postSettingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	host := r.PathValue("host")
	port := r.PathValue("port")
	setting := r.PathValue("setting")
        err := r.ParseForm()
        if err != nil {
                log.Printf("postSettingHandler service: parse form error %v", err)
                return
        }
        settingvalue := r.Form.Get("settingvalue")
	FluidsynthCommand(host, port, "set "+setting+" "+settingvalue)
}

func postMidiccHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	host := r.PathValue("host")
	port := r.PathValue("port")
	channel := r.PathValue("channel")
	control := r.PathValue("control")
        err := r.ParseForm()
        if err != nil {
                log.Printf("postMidiccHandler service: parse form error %v", err)
                return
        }
        ccvalue := r.Form.Get("ccvalue")
	FluidsynthCommand(host, port, "cc "+channel+" "+control+" "+ccvalue)
}

func postSetVolumeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	host := r.PathValue("host")
	port := r.PathValue("port")
	channel := r.PathValue("channel")
        err := r.ParseForm()
        if err != nil {
                log.Printf("postSetVolumeHandler service: parse form error %v", err)
                return
        }
        volume := r.Form.Get("volume")
	FluidsynthCommand(host, port, "cc "+channel+" "+"7"+" "+volume)
	//engineKey := EngineKey{}
	//engineKey.Host = host
	//engineKey.Port = port
	//engine := engines[engineKey]
	//channelRecord := engine.Channels[channel]
	//engine.Channels[channel] = channelRecord 
}


func postSelectFontHandler(w http.ResponseWriter, r *http.Request) {
	host := r.PathValue("host")
	port := r.PathValue("port")
	channel := r.PathValue("channel")
        err := r.ParseForm()
        if err != nil {
                log.Printf("postSelectFontHandler service: parse form error %v", err)
                return
        }
        font := r.Form.Get("font")
	//engineKey := EngineKey{}
	//engineKey.Host = host
	//engineKey.Port = port
	//engine := engines[engineKey]
	//channelRecord := engine.Channels[channel]
	//engine.Channels[channel] = channelRecord 
	log.Printf("Select Font %v",r.Form)
        FluidsynthCommand(host, port, "select "+channel+" "+font+" 0 0") 
	triggerEvent, exists := r.Header["Triggering-Event"]
	if exists {
		var evt map[string]interface{}
		if err = json.Unmarshal([]byte(triggerEvent[0]), &evt); err != nil {
			log.Printf("Triggering-event unmarshal error  %v",err)
		} else {
			if evt["type"] == "change" {
	   			err = templates.ExecuteTemplate(w, "fluidprogram", FluidChannel(host,port,channel,font))
	   			if err != nil {
	   				log.Printf("Error executing fluidprogram template %v", err)
	   			}
		   		return
			}
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func postSetBankPresetHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	host := r.PathValue("host")
	port := r.PathValue("port")
	channel := r.PathValue("channel")
        err := r.ParseForm()
        if err != nil {
                log.Printf("postSetBankPresetHandler service: parse form error %v", err)
                return
        }
        bankpreset := r.Form.Get("bankpreset")
	font := r.Form.Get("font")
	bank := bankpreset[0:3]
	preset := bankpreset[4:7]
	//engineKey := EngineKey{}
	//engineKey.Host = host
	//engineKey.Port = port
	//engine := engines[engineKey]
	//channelRecord := engine.Channels[channel]
	//engine.Channels[channel] = channelRecord 
	FluidsynthCommand(host, port, "select "+channel+" "+font+" "+bank+" "+preset) 
}

func FluidsynthCommand(host string, port string, command string) ([]byte, error) {
	Mutex.Lock()
	defer Mutex.Unlock()

	log.Printf("Fluidsynth send: host %s port %s command %s", host, port, command)
	tcpAddr, err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err != nil {
		log.Printf("Error %v",err)
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Printf("Error %v",err)
		return nil, err
	}
	conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))

	_, err = conn.Write([]byte(command + "\n"))
	if err != nil {
		log.Printf("Error %v",err)
		return nil, err
	}

	response := make([]byte, 0)
	reply := make([]byte, 8196)

	goon := true
	for goon {

		l, err := conn.Read(reply)
		response = append(response, reply[:l]...)
		if err != nil {
			//log.Printf("Error %v",err)
			goon = false
		}

	}

	conn.Close()
	return response, nil
}

