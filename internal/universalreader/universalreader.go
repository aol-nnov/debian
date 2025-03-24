package universalreader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mholt/archives"
)

func New(src string) (io.Reader, error) {

	var srcReader io.Reader
	var err error

	if strings.HasPrefix(src, "http") {
		client := http.DefaultClient
		var req *http.Request

		if req, err = http.NewRequest("GET", src, nil); err != nil {
			return nil, err
		}

		if authUser := os.Getenv("HTTP_USER"); authUser != "" {
			req.SetBasicAuth(authUser, os.Getenv("HTTP_PASSWD"))
		}

		if res, err := client.Do(req); err != nil {
			return nil, err
		} else {
			if res.StatusCode == 200 {
				srcReader = res.Body
			} else {
				return nil, fmt.Errorf("UniversalReader: %s", res.Status)
			}
		}
	} else {
		if srcReader, err = os.Open(src); err != nil {
			return nil, err
		}
	}

	if format, input, err := archives.Identify(context.TODO(), "", srcReader); err != nil {
		// format not detected (i.e. file is not compressed)
		return input, nil
	} else {
		// fmt.Printf("%s %v\n", format.Name(), input)
		if decompressor, ok := format.(archives.Decompressor); ok {
			return decompressor.OpenReader(input)
		}
	}
	return nil, fmt.Errorf("NewUniversalReader: stranger things")
}

func MaybeClose(reader io.Reader) error {
	if closer, ok := reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
