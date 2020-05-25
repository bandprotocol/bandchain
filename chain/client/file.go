package rpc

import (
	"net/http"
	"path/filepath"

	"github.com/bandprotocol/bandchain/chain/pkg/filecache"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

const (
	Filename       = "filename"
	FlagHomeDaemon = "daemon-home"
)

func GetFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filename := vars[Filename]
		dir := filepath.Join(viper.GetString(FlagHomeDaemon), "files")
		f := filecache.New(dir)
		file := f.GetFile(filename)
		w.Header().Set("Content-Disposition", "attachment;")
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

		w.Write(file)
	}
}
