package handler

import (
	"github.com/sicozz/project00/datatype"
	"github.com/sicozz/project00/utils"
)

func QueryGetServiceInfo() datatype.ProgramInfo {
	return datatype.ProgramInfo{Version: utils.VERSION, Banner: utils.BANNER}
}
