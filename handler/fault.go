package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/Jarvis-Sui/chaos-os/binding"
	"github.com/Jarvis-Sui/chaos-os/manager"
	"github.com/Jarvis-Sui/chaos-os/util"
	"github.com/labstack/echo/v4"
	"github.com/spf13/pflag"
)

func GetFaults(c echo.Context) error {
	if err := checkQueryParams(c.QueryParams(), []string{"id", "status", "limit"}); err != nil {
		return c.JSON(http.StatusBadRequest, errorMsg(err))
	}
	var id, state string
	var limit int

	flags := pflag.FlagSet{}
	flags.StringVar(&id, "id", "", "fault id")
	flags.StringVar(&state, "status", "", "status of faults to return. Ready | Running | Error | Destroyed")
	flags.IntVar(&limit, "limit", 100, "maximum number of faults returned")

	flags.Set("id", c.QueryParam("id"))
	flags.Set("status", c.QueryParam("status"))
	flags.Set("limit", c.QueryParam("limit"))
	faults := manager.FaultStatus(&flags)
	return c.JSON(http.StatusOK, faults)
}

func GetFaultById(c echo.Context) error {
	var id string
	flags := pflag.FlagSet{}
	flags.StringVar(&id, "id", "", "fault id")

	flags.Set("id", c.Param("id"))
	faults := manager.FaultStatus(&flags)
	return c.JSON(http.StatusOK, faults)
}

func AddFault(c echo.Context) error {
	if err := checkRequiredQueryParams(c.QueryParams(), []string{"cmd"}); err != nil {
		return c.JSON(http.StatusBadRequest, errorMsg(err))
	}
	cmd := exec.Command("bash", "-c", fmt.Sprintf("%s fault create %s", util.GetBinPath(), c.QueryParam("cmd")))

	if out, err := cmd.CombinedOutput(); err != nil {
		return c.JSON(http.StatusBadRequest, errorMsg(err))
	} else {
		var fault = binding.Fault{}
		if err := json.Unmarshal(out, &fault); err != nil {
			return c.JSON(http.StatusBadRequest, errorMsg(fmt.Errorf("%s", string(out))))
		} else {
			return c.JSON(http.StatusCreated, fault)
		}
	}
}

func DestroyFault(c echo.Context) error {
	if err := checkRequiredQueryParams(c.QueryParams(), []string{"id"}); err != nil {
		return c.JSON(http.StatusBadRequest, errorMsg(err))
	}

	var id string
	flags := pflag.FlagSet{}
	flags.StringVar(&id, "id", "", "fault id")
	flags.Set("id", c.QueryParam("id"))

	if err := manager.DestroyFault(&flags); err != nil {
		return c.JSON(http.StatusInternalServerError, errorMsg(err))
	} else {
		faults := manager.FaultStatus(&flags)
		return c.JSON(http.StatusOK, faults)
	}
}
