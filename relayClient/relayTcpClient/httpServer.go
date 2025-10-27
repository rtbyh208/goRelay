package relaytcpclient

import (
	"embed"
	"encoding/json"
	"fmt"
	dbcache "goRelay/dbCache"
	"goRelay/pkg"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//go:embed static/*
var content embed.FS

type requestJsonDelClient struct {
	ID string `json:"id"`
}

type requestJsonSwitchService struct {
	ID string `json:"id"`
}

type requestJsonAdd_UpdateClient struct {
	ID     string `json:"id"`
	Value  string `json:"value"`
	Remask string `json:"remask"`
}

type responseMsgs struct {
	Status  bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type responseLstClents struct {
	Value  string `json:"value"`
	Remask string `json:"remask"`
	TX     string `json:"TX"`
	RX     string `json:"RX"`
	Status bool   `json:"status"`
}

type relayServerInfos struct {
	Key                   string `json:"key" yaml:"key"`
	Id                    string `json:"id" yaml:"id"`
	PipeServerAddr        string `json:"pipe_server_addr" yaml:"pipe_server_addr"`
	ListenRelayServerAddr string `json:"listen_relay_server_addr" yaml:"listen_relay_server_addr"`
	DebugLog              bool   `json:"debug_log" yaml:"debug_log"`
	RegisterID            string `json:"register_id" yaml:"register_id"`
	DirectConn            string `json:"direct_conn" yaml:"direct_conn"`
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func StartHttpServer(key string, httpserver string, pipeServerAddr string, registerID string, cache *dbcache.Caches) {

	http.HandleFunc("/modifyClient", func(w http.ResponseWriter, r *http.Request) {

		var client requestJsonAdd_UpdateClient
		decoder := json.NewDecoder(r.Body)

		decoder.Decode(&client)

		var clientInfos responseMsgs

		if client.Value == "" {
			clientInfos.Status = false
			clientInfos.Message = "value为空,错误"
			goto returnClientInfos
		}

		if client.ID == "" {
			new_id := pkg.IDHash(fmt.Sprintf("%s%d", randStrings(6), time.Now().Unix()))
			cache.Set(pkg.IDHash(new_id), client.Value)
			clientInfos.Status = true

			remask_id := fmt.Sprintf("Remask_%s", pkg.IDHash(new_id))
			cache.Set(remask_id, client.Remask)

			clientInfos.Message = fmt.Sprintf("<b style='color: red;'>以下信息仅出现一次，请妥善保管！如有必要，请修改pipeserver的连接地址</b>" +
				"</br></br></br>" +
				"启动的命令:</br>" +
				"</b>" +
				"</br>" +
				"</br>" +
				"如:</br>" +
				"./goRelay_linux_amd64 relayServer --conf relayServerConf.yaml" +
				"</br>" +
				"</br>" +
				"<b style='color: red;'></b>" +
				"</br>" +
				"</br>" +
				"</br>" +
				"")

			var c relayServerInfos
			c.Key = key
			c.Id = new_id
			c.PipeServerAddr = pipeServerAddr
			c.RegisterID = registerID
			c.DebugLog = true
			c.ListenRelayServerAddr = ":6000"

			clientInfos.Data = c

		} else {
			if cache.Exists(client.ID) {
				cache.Set(client.ID, client.Value)

				remask_id := fmt.Sprintf("Remask_%s", client.ID)
				cache.Set(remask_id, client.Remask)

				clientInfos.Status = true
				clientInfos.Message = "修改成功"
				goto returnClientInfos
			}
			clientInfos.Status = true
			clientInfos.Message = "未找到记录"
		}

	returnClientInfos:
		jsonData, err := json.Marshal(clientInfos)
		if err != nil {
			goLog.Error("marshal error ", err)
		}
		w.Write(jsonData)
	})

	http.HandleFunc("/delClient", func(w http.ResponseWriter, r *http.Request) {

		var client requestJsonDelClient
		decoder := json.NewDecoder(r.Body)

		decoder.Decode(&client)

		var clientInfos responseMsgs

		if cache.Exists(client.ID) {
			cache.Del(client.ID)

			clientInfos.Status = true
			clientInfos.Message = "删除成功"
		} else {
			var clientInfos responseMsgs
			clientInfos.Status = false
			clientInfos.Message = "删除失败"
		}

		jsonData, err := json.Marshal(clientInfos)
		if err != nil {
			goLog.Error("marshal error ", err)
		}
		w.Write(jsonData)
	})

	http.HandleFunc("/listClient", func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}

		clientMaps := make(map[string]responseLstClents, 0)

		for k, v := range cache.List() {
			if !strings.HasPrefix(k, "Stop_") && !strings.HasPrefix(k, "Usage_") && !strings.HasPrefix(k, "Remask_") {
				var rc responseLstClents
				rc.Value = v
				var newRX string = fmt.Sprintf("Usage_RX_%s", k)
				var newTx string = fmt.Sprintf("Usage_TX_%s", k)
				rc.TX = numberConver(cache.Cache[newTx])
				rc.RX = numberConver(cache.Cache[newRX])

				remask_id := fmt.Sprintf("Remask_%s", k)
				rc.Remask, _ = cache.Get(remask_id)
				rc.Status = !cache.Exists(fmt.Sprintf("Stop_%s", k))

				clientMaps[k] = rc
			}
		}

		jsonData, err := json.Marshal(clientMaps)
		if err != nil {
			goLog.Error("marshal error ", err)
		}
		w.Write(jsonData)
	})

	http.HandleFunc("/switchServicee", func(w http.ResponseWriter, r *http.Request) {

		var client requestJsonDelClient
		decoder := json.NewDecoder(r.Body)

		decoder.Decode(&client)

		var clientInfos responseMsgs

		stopClientID := fmt.Sprintf("Stop_%s", client.ID)

		if client.ID == "" {
			clientInfos.Status = false
			clientInfos.Message = "ID为,切换失败"
			goto returnClientInfos
		}

		if cache.Exists(stopClientID) {
			cache.Del(stopClientID)
			clientInfos.Status = true
			clientInfos.Message = "启用成功"
			goto returnClientInfos
		} else {
			cache.Set(stopClientID, "")
			clientInfos.Status = true
			clientInfos.Message = "禁用成功"
			goto returnClientInfos
		}

	returnClientInfos:
		jsonData, err := json.Marshal(clientInfos)
		if err != nil {
			goLog.Error("marshal error ", err)
		}
		w.Write(jsonData)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := content.ReadFile("static" + r.URL.Path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Write(data)
	})

	goLog.Info("listen http server: ", httpserver)
	err := http.ListenAndServe(httpserver, nil)
	if err != nil {
		goLog.Error("listen ", httpserver, " error: ", err)
	}
}

func randStrings(length int) []byte {
	source := rand.NewSource(time.Now().UnixMicro())
	r := rand.New(source)

	var result []byte
	for i := 0; i < length; i++ {
		result = append(result, charset[r.Intn(len(charset))])
	}
	return result
}

func numberConver(x string) string {

	if x == "" {
		return "0 Bytes"
	}

	converNumer, err := strconv.ParseUint(x, 10, 64)
	if err != nil {
		return "full max"
	}

	var capactiyLevels = []string{"Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	var level int = 0

	for converNumer > (1024 * 1024) {
		level += 1
		converNumer = converNumer / 1024
	}

	var converNumerFloat float64

	converNumerFloat = float64(converNumer) / 1024

	return fmt.Sprintf("%.6v%s", converNumerFloat, capactiyLevels[level+1])
}
