package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Instance holds data to send request to designated url.
type Instance struct {
	Cl  *http.Client    // Http client to send request.
	Url string          // Should be filled with CallbackUrl.
	Ctx context.Context // Context to run this http client.
}

// DispatchPOST send POST request to designated url using given http client.
func (i *Instance) DispatchPOST(r interface{}) error {
	js, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("failed to marshaling docker.StdResponse to bytes: %s", err)
	}

	req, err := http.NewRequestWithContext(i.Ctx, fiber.MethodPost, i.Url, bytes.NewBuffer(js))
	if err != nil {
		return fmt.Errorf("failed creating request instance: %s", err)
	}
	req.Header.Set("content-type", fiber.MIMEApplicationJSON)

	_, err = i.Cl.Do(req)
	if err != nil {
		return fmt.Errorf("failed sending POST request: %s", err)
	}

	return nil
}
