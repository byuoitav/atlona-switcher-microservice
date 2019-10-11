package switcher2x1

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/byuoitav/atlona-switcher-microservice/structs"
	"github.com/byuoitav/common/nerr"
)

// GetInput returns the current input for the requested output
func GetInput(ctx context.Context, address string) (string, *nerr.E) {
	var resp structs.WallPlateStruct
	url := fmt.Sprintf("http://%s/aj.html?a=avs", address)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nerr.Translate(err).Addf("error when making request: %s", err)
	}
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nerr.Translate(err).Addf("error when making call: %s", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal([]byte(body), &resp) // here!
	if err != nil {
		return "", nerr.Translate(err).Addf("error when unmarshalling the response: %s", err)
	}
	defer res.Body.Close()

	return strconv.Itoa(resp.Inp), nil
}
