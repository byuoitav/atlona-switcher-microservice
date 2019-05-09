package switcher

import (
	"context"
	"fmt"

	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/structs"
)

// GetInput returns the current input
func GetInput(ctx context.Context, address string) (string, *nerr.E) {
	var settings AVSettings
	err := getPage(ctx, address, avSettingsPage, &settings)
	if err != nil {
		return "", err.Addf("unable to get input")
	}

	return fmt.Sprintf("%v", settings.Input-1), nil
}

// GetHardwareInfo returns a hardware info struct
func GetHardwareInfo(ctx context.Context, address string) (structs.HardwareInfo, *nerr.E) {
	var hwinfo structs.HardwareInfo

	var info Info
	err := getPage(ctx, address, infoPage, &info)
	if err != nil {
		return hwinfo, err.Addf("unable to get hardware info")
	}

	// fill in the hwinfo
	if len(info.SystemInfo) >= 1 {
		hwinfo.ModelName = info.SystemInfo[0]
	}

	if len(info.SystemInfo) >= 2 {
		hwinfo.FirmwareVersion = info.SystemInfo[1]
	}

	return hwinfo, nil
}
