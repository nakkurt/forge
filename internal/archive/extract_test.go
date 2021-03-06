package archive

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchAndExtract(t *testing.T) {
	srv := httptest.NewServer(nil)
	defer srv.Close()

	cases := []struct {
		name    string
		archive string
		files   []string
	}{
		{"tarball", "testdata/simple-app.tar", []string{"Dockerfile", "app.py"}},
		{"gzipped-tarball", "testdata/simple-app.tgz", []string{"Dockerfile", "app.py"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				bs, err := ioutil.ReadFile(tc.archive)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := w.Write(bs); err != nil {
					t.Fatal(err)
				}
			})

			ext, err := FetchAndExtract(srv.URL)
			if err != nil {
				t.Error(err)
			}

			fi, err := os.Stat(ext.Archive)
			if os.IsNotExist(err) {
				t.Errorf("archive was not created %q", ext.Archive)
			}
			if !fi.Mode().IsRegular() {
				t.Errorf("archive is not a regular file %q", ext.Archive)
			}

			var actual []string
			err = filepath.Walk(ext.ContentsDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					return nil
				}

				path, _ = filepath.Rel(ext.ContentsDir, path)
				actual = append(actual, path)

				return nil
			})
			if err != nil {
				t.Errorf("failed to list extracted archive files: %v", err)
			}

			assert.ElementsMatch(t, tc.files, actual, "expected archive contents to match")

			os.RemoveAll(ext.RootDir)
		})
	}

	t.Run("fetch-failure", func(t *testing.T) {
		t.SkipNow()
	})

	t.Run("unsupported-format", func(t *testing.T) {
		t.SkipNow()
	})
}
