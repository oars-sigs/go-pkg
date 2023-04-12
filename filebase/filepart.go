package filebase

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/http"
	"os"

	"pkg.oars.vip/go-pkg/req"
	"pkg.oars.vip/go-pkg/server/base"
)

type CreateMultipartResp struct {
	base.DataResponse
	Data *FileUpload `json:"data"`
}

type FileUpload struct {
	ID       string `gorm:"column:upload_id" json:"uploadId"`
	Creator  string `gorm:"column:upload_creator" json:"creator"`
	PartNum  int    `gorm:"column:upload_part_num" json:"partNum"`
	Digest   string `gorm:"column:upload_digest" json:"digest"`
	Metadata string `gorm:"column:upload_metadata" json:"metadata"`
	Parent   string `gorm:"-" json:"parent"`
	Created  int64  `gorm:"column:upload_created;autoCreateTime:milli" json:"created"`
	Updated  int64  `gorm:"column:upload_updated;autoUpdateTime:milli" json:"updated"`
	Status   string `gorm:"-" json:"status"`
}

type createMultipartReq struct {
	Parent string `json:"parent"`
	Digest string `json:"digest"`
}

func (c *Client) CreateMultipart(namespace, parentId, digest string) (*CreateMultipartResp, error) {
	ustr := fmt.Sprintf("%s/filebase/api/v1/%s/multiparts", c.cfg.Address, namespace)
	in := createMultipartReq{
		Parent: parentId,
		Digest: digest,
	}
	var out CreateMultipartResp
	err := req.ReqJSON("POST", ustr, in, &out, c.setAuthHeader(nil))
	if err != nil {
		return nil, err
	}
	if out.Error.Error() != nil {
		return nil, out.Error.Error()
	}
	return &out, nil
}

func (c *Client) PutMultipart(path, namespace, uploadId string, num int) (int, error) {
	fs, err := os.Open(path)
	if err != nil {
		return num, err
	}
	info, err := fs.Stat()
	if err != nil {
		return num, err
	}
	filesize := info.Size()
	const filechunk = 4 * 1 << 20
	blocks := int(math.Ceil(float64(filesize) / float64(filechunk)))

	for i := num; i < blocks; i++ {
		hash := md5.New()
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)
		fs.Read(buf)
		io.WriteString(hash, string(buf))
		digest := hex.EncodeToString(hash.Sum(nil))
		err := c.putPart(buf, namespace, uploadId, digest, i)
		if err != nil {
			return i, err
		}
	}
	return num, nil

}

func (c *Client) MergeMultiPart(namespace, uploadId string, m *FileMetadata) error {
	ustr := fmt.Sprintf("%s/filebase/api/v1/%s/multiparts/%s", c.cfg.Address, namespace, uploadId)
	var out base.DataResponse
	err := req.ReqJSON("PUT", ustr, m, &out, c.setAuthHeader(nil))
	if err != nil {
		return err
	}
	return out.Error.Error()
}

func (c *Client) AbortMultiPart(namespace, uploadId string) error {
	ustr := fmt.Sprintf("%s/filebase/api/v1/%s/multiparts/%s", c.cfg.Address, namespace, uploadId)
	var out base.DataResponse
	err := req.ReqJSON("DELETE", ustr, nil, &out, c.setAuthHeader(nil))
	if err != nil {
		return err
	}
	return out.Error.Error()
}

func (c *Client) putPart(body []byte, namespace, uploadId, digest string, num int) error {
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	part, err := writer.CreateFormFile("file", "file.txt")
	if err != nil {
		return err
	}
	part.Write(body)
	writer.Close()
	ustr := fmt.Sprintf("%s/filebase/api/v1/%s/multiparts/%s/%d?md5=%s", c.cfg.Address, namespace, uploadId, num, digest)
	fmt.Println(ustr)
	resp, err := req.Req("POST", ustr, requestBody, c.setAuthHeader(map[string]string{"Content-Type": writer.FormDataContentType()}))
	if err != nil {
		return err
	}
	defer resp.Close()
	data, err := ioutil.ReadAll(resp)
	if err != nil {
		return err
	}
	var res base.DataResponse
	err = json.Unmarshal(data, &res)
	if err != nil {
		return err
	}
	if res.Error.Error() != nil {
		return res.Error.Error()
	}
	return nil
}

func (c *Client) Download(namespace, id, localFilePath string) error {
	bufferSize := 1024 * 1024
	fileUrl := fmt.Sprintf("%s/filebase/api/v1/%s/files/%s", c.cfg.Address, namespace, id)
	fileInfo, err := os.Stat(localFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(localFilePath)
			if err != nil {
				return err
			}
			file.Close()
			fileInfo, err = os.Stat(localFilePath)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	file, err := os.OpenFile(localFilePath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	request, err := http.NewRequest("GET", fileUrl, nil)
	if err != nil {
		return err
	}

	request.Header.Set("Range", fmt.Sprintf("bytes=%d-", fileInfo.Size()))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// 如果服务器返回206 Partial Content，则支持断点续传
	if response.StatusCode == http.StatusPartialContent {
		_, err = file.Seek(fileInfo.Size(), 0)
		if err != nil {
			return err
		}
	}

	buffer := make([]byte, bufferSize)
	for {
		n, err := response.Body.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		_, err = file.Write(buffer[:n])
		if err != nil {
			return err
		}
	}
	return err
}
