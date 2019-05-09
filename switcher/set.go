package switcher

import (
	"context"
	"fmt"
	"strconv"

	"github.com/byuoitav/common/nerr"
)

// SwitchInput changes the input on the given output to input
func SwitchInput(ctx context.Context, address, output, input string) *nerr.E {
	// atlona switchers are 1-based
	out, gerr := strconv.Atoi(output)
	if gerr != nil {
		return nerr.Translate(gerr).Addf("unable to switch input on %s", address)
	}

	in, gerr := strconv.Atoi(input)
	if gerr != nil {
		return nerr.Translate(gerr).Addf("unable to switch input on %s", address)
	}

	out++
	in++

	// validate that input/output are valid numbers
	var settings AVSettings
	err := getPage(ctx, address, avSettingsPage, &settings)
	if err != nil {
		return err.Addf("unable to switch input")
	}

	if in > len(settings.HDCPSettings) || in <= 0 {
		return nerr.Createf("error", "unable to switch input on %s - input %s is out of range", address, input)
	}

	if out != 1 {
		return nerr.Createf("error", "unable to switch input on %s - output %s is invalid", address, output)
	}

	err = sendCommand(ctx, address, fmt.Sprintf("x%vAVx%v", in, out))
	if err != nil {
		return err.Addf("unable to switch input")
	}

	return nil
}
