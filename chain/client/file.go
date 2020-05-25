package rpc

import (
	"net/http"
	"path/filepath"

	"github.com/bandprotocol/bandchain/chain/pkg/filecache"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

const (
	Filename       = "filename"
	FlagHomeDaemon = "daemon-home"
)

func GetFile() http.HandlerFunc {
	fileCache := filecache.New(filepath.Join(viper.GetString(FlagHomeDaemon), "files"))
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filename := vars[Filename]
		file, err := fileCache.GetFile(filename)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		w.Header().Set("Content-Disposition", "attachment;")
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.Write(file)
	}
}
