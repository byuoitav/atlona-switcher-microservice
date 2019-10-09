package switcher2x1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/byuoitav/common/nerr"
)

// SetInput changes the input on the given output to input
func SetInput(ctx context.Context, address, input string) *nerr.E {

	url := fmt.Sprintf("http://%s/aj.html?a=command&cmd=x%sAVx1", address, input)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nerr.Translate(err).Addf("error when making request: %s", err)
	}
	req = req.WithContext(ctx)
	res, gerr := http.DefaultClient.Do(req)
	if gerr != nil {
		return nerr.Translate(gerr).Addf("error when making call: %s", gerr)
	}
	defer res.Body.Close()
	return nil

}
