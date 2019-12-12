package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"encoding/base64"
	"encoding/hex"

	"github.com/bandprotocol/d3n/chain/cmtx"
	"github.com/bandprotocol/d3n/chain/x/zoracle"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

var (
	port     string
	nodeURI  string
	queryURI string
	priv     secp256k1.PrivKeySecp256k1
)

type requestData struct {
	Code  string `json:"code"`
	Delay string `json:"delay"`
}

type Tx struct {
	Hash     string `json:"hash"`
	Height   string `json:"height"`
	Index    int    `json:"index"`
	TxResult struct {
		Code      int         `json:"code"`
		Data      interface{} `json:"data"`
		Log       string      `json:"log"`
		Info      string      `json:"info"`
		GasWanted string      `json:"gasWanted"`
		GasUsed   string      `json:"gasUsed"`
		Events    []struct {
			Type       string `json:"type"`
			Attributes []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"attributes"`
		} `json:"events"`
		Codespace string `json:"codespace"`
	} `json:"tx_result"`
	Tx string `json:"tx"`
}

func handleTestReq(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	if _, ok := params["id"]; !ok {
		w.Write([]byte("404 page not found"))
		return
	}
	fmt.Println(params)
	a := []string{"0xaa", "0xbb", "0xcc", params["id"][0]}
	type x struct {
		TTT []string `json:"ttt"`
	}
	xx := x{}
	xx.TTT = a
	json.NewEncoder(w).Encode(xx)
}

// has0xPrefix validates str begins with '0x' or '0X'.
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func GetTx(txHash string) (map[string]interface{}, error) {
	_resp, err := http.Get(fmt.Sprintf("%s/tx?hash=0x%s", strings.Replace(nodeURI, "tcp", "http", 1), txHash))
	if err != nil {
		return nil, err
	}
	defer _resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(_resp.Body)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	var sreqID string
	if val, ok := params["reqID"]; ok {
		sreqID = val[0]
	} else {
		w.Write([]byte("404 page not found"))
		return
	}
	w.Header().Set("Content-Type", "application/json")

	reqID, err := strconv.ParseUint(sreqID, 10, 64)
	if err != nil {
		http.Error(w, "reqID must be positive integer", http.StatusBadRequest)
		return
	}

	type resp struct {
		Height string `json:"height"`
		Result struct {
			CodeHash  string      `json:"codeHash"`
			ReportEnd int         `json:"reportEnd"`
			Reports   interface{} `json:"reports"`
			RequestID int         `json:"requestID"`
			Result    string      `json:"result"`
		} `json:"result"`
	}

	_resp, err := http.Get(fmt.Sprintf("%s/zoracle/request/%d", queryURI, reqID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseBytes, err := ioutil.ReadAll(_resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer _resp.Body.Close()

	response := resp{}
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	codeHash, _ := base64.StdEncoding.DecodeString(response.Result.CodeHash)
	result, _ := base64.StdEncoding.DecodeString(response.Result.Result)
	if len(codeHash) > 0 {
		response.Result.CodeHash = fmt.Sprintf("0x%x", codeHash)
	}
	if len(result) > 0 {
		response.Result.Result = fmt.Sprintf("0x%x", result)
	}

	json.NewEncoder(w).Encode(response)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(map[string]string{"status": "400", "message": "Only POST method is supported."})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	tx := cmtx.NewTxSender(priv)

	rd := requestData{}
	err := json.NewDecoder(r.Body).Decode(&rd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	delay, err := strconv.ParseUint(rd.Delay, 10, 64)
	if err != nil {
		http.Error(w, "delay should be positive integer", http.StatusBadRequest)
		return
	}
	if delay < 5 {
		http.Error(w, "delay should be at least 5 blocks", http.StatusBadRequest)
		return
	}

	var code []byte
	if has0xPrefix(rd.Code) {
		code, err = hex.DecodeString(rd.Code[2:])
	} else {
		code, err = hex.DecodeString(rd.Code)
	}
	if err != nil {
		http.Error(w, "wrong code format", http.StatusBadRequest)
		return
	}

	txHash, err := tx.SendTransaction(zoracle.NewMsgRequest(code, delay, tx.Sender()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var txResponse Tx
	for i := 0; i < 30; i++ {
		tr, err := GetTx(txHash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if txData, ok := tr["result"]; ok {
			byteData, _ := json.Marshal(txData)
			json.Unmarshal(byteData, &txResponse)
			break
		}
		time.Sleep(time.Second / 2)
	}

	if txResponse.Hash != txHash {
		http.Error(w, "tx was corrupted"+fmt.Sprintf("%v", txResponse), http.StatusInternalServerError)
		return
	}

	height, err := strconv.ParseUint(txResponse.Height, 10, 64)
	if err != nil {
		http.Error(w, "tx was corrupted"+fmt.Sprintf("%v", txResponse), http.StatusInternalServerError)
		return
	}

	m := map[string]string{}
	m["txHash"] = txResponse.Hash
	m["startHeight"] = txResponse.Height
	m["endHeight"] = fmt.Sprintf("%d", height+delay)

	events := txResponse.TxResult.Events
	for _, event := range events {
		if event.Type == "request" {
			for _, attr := range event.Attributes {
				k, _ := base64.StdEncoding.DecodeString(attr.Key)
				if string(k) == "id" {
					v, _ := base64.StdEncoding.DecodeString(attr.Value)
					m["requestId"] = string(v)
				}
			}
		}
	}

	json.NewEncoder(w).Encode(m)
}

func handleGetProof(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	if _, ok := params["reqID"]; !ok {
		w.Write([]byte("404 page not found"))
		return
	}

	u, err := strconv.ParseUint(params["reqID"][0], 10, 64)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	ap, err := GetProof(u, "a96e62ed3955e65be3aaa3f12d87b6b5cf26039ecfa948dc5107a495418e5430")

	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(ap)
}

func main() {
	var ok bool
	port, ok = os.LookupEnv("PORT")
	if !ok {
		port = "5000"
	}
	nodeURI, ok = os.LookupEnv("NODE_URI")
	if !ok {
		nodeURI = "tcp://localhost:26657"
	}
	queryURI, ok = os.LookupEnv("QUERY_URI")
	if !ok {
		queryURI = "http://localhost:1317"
	}
	privS, ok := os.LookupEnv("PRIVATE_KEY")
	if !ok {
		log.Fatal("Missing private key")
	}

	privB, _ := hex.DecodeString(privS)
	copy(priv[:], privB)

	http.HandleFunc("/request", handleRequest)
	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/proof", handleGetProof)
	http.HandleFunc("/test", handleTestReq)
	fmt.Println("live!")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
