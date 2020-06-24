package main

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/client/flags"
	keyring "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

func runCmd(c *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{"r"},
		Short:   "Run the oracle process",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.ChainID == "" {
				return errors.New("Chain ID must not be empty")
			}
			keys, err := keybase.List()
			if err != nil {
				return err
			}
			if len(keys) == 0 {
				return errors.New("No key available")
			}
			c.keys = make(chan keyring.Info, len(keys))
			for _, key := range keys {
				c.keys <- key
			}
			c.gasPrices, err = sdk.ParseDecCoins(cfg.GasPrices)
			if err != nil {
				return err
			}
			c.client, err = rpcclient.NewHTTP(cfg.NodeURI, "/websocket")
			if err != nil {
				return err
			}
			c.amount = sdk.NewCoins(sdk.NewCoin("uband", sdk.NewInt(cfg.Amount)))
			r := gin.Default()
			r.Use(func(c *gin.Context) {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
				c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
				c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")

				if c.Request.Method == "OPTIONS" {
					c.AbortWithStatus(204)
					return
				}
			})

			r.POST("/request", func(gc *gin.Context) {
				handleRequest(gc, c)
			})

			return r.Run("0.0.0.0:" + cfg.Port)
		},
	}
	cmd.Flags().String(flags.FlagChainID, "", "chain ID of BandChain network")
	cmd.Flags().String(flags.FlagNode, "tcp://localhost:26657", "RPC url to BandChain node")
	cmd.Flags().String(flags.FlagGasPrices, "", "gas prices for report transaction")
	cmd.Flags().String(flagPort, "5005", "port of faucet service")
	cmd.Flags().Int64(flagAmount, 10000000, "amount in uband for each request")
	viper.BindPFlag(flags.FlagChainID, cmd.Flags().Lookup(flags.FlagChainID))
	viper.BindPFlag(flags.FlagNode, cmd.Flags().Lookup(flags.FlagNode))
	viper.BindPFlag(flags.FlagGasPrices, cmd.Flags().Lookup(flags.FlagGasPrices))
	viper.BindPFlag(flagPort, cmd.Flags().Lookup(flagPort))
	viper.BindPFlag(flagAmount, cmd.Flags().Lookup(flagAmount))
	return cmd
}
