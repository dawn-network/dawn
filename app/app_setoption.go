package app

import (
	"encoding/hex"
	"log"
	"strings"
	"strconv"
)

/**
glogChainApp.SetOption("genesis.block/create.account", "pool/449F6F39391BC9E918CE51DB10F6FAADF65077263E715B17652425CD7827C814/1000000")
glogChainApp.SetOption("genesis.block/create.account", "jan/CDD6774218138DF657C7B0494894BBA596EB2ECCCC4C946D5ACEF3B5FCD2CE7B/1000")
glogChainApp.SetOption("genesis.block/create.account", "jake/488B8FF58E8E9868823C3388BAAB9C1F7CFCB3D7482376E7495639A1EC0F7407/1000")
glogChainApp.SetOption("genesis.block/create.account", "tuan/EB3B42091EF6C2F8FA951319940C003BEC7AAE2336BD2AFABD6FB59EB4A3EF6E/1000")
 */
func Exec_SetOption(app *GlogChainApp, key string, value string) (logstr string) {
	var err error

	switch key {
	case "genesis.block/create.account":
		strs := strings.Split(value, "/")
		if (len(strs) != 3) {
			return "Invalid input data"
		}

		var account Account
		account.PubKey, err = hex.DecodeString(strs[1])
		if (err != nil) {
			log.Println(err.Error())
			return err.Error()
		}
		account.Sequence = 1

		account.Balance, err = strconv.ParseInt(strs[2], 10, 64)
		if err != nil {
			return err.Error()
		}

		err = TreeSaveAccount(app.State, account)
		if (err != nil) {
			log.Println(err.Error())
			return err.Error()
		}

		break
	case "genesis.block/token.transfer":
		break
	default:
	}

	return ""
}