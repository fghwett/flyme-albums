package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type DirReply struct {
	Code     int       `json:"code"`
	Message  string    `json:"message"`
	Value    *DirValue `json:"value"`
	Redirect string    `json:"redirect"`
}

type DirValue struct {
	Time  int64    `json:"time"`
	Count int      `json:"count"`
	Dir   []*Album `json:"dir"`
	End   int      `json:"end"`
}

type Album struct {
	Id         int     `json:"id"`
	DirName    string  `json:"dirName"`
	ParentId   int     `json:"parentId"`
	UserId     int     `json:"userId"`
	Status     int     `json:"status"`
	TotalSize  int     `json:"totalSize"`
	FileNum    int     `json:"fileNum"`
	Icon       string  `json:"icon"`
	CreateTime int64   `json:"createTime"`
	ModifyTime int64   `json:"modifyTime"`
	SqlNow     int64   `json:"sqlNow"`
	Files      []*File `json:"files,omitempty"`

	token string
}

func (a *Album) GetPhotos(token string) error {
	a.token = token

	limit := 100
	offset := 0

	var fs []*File
	var hasNext bool
	var err error

	i := 1
	for {
		offset = (i - 1) * limit
		fs, hasNext, err = a.getPhoto(a.Id, limit, offset)
		if err != nil {
			return fmt.Errorf("获取相片失败 albumId: %d, limit: %d, offset: %d, err: %s", a.Id, limit, offset, err)
		}

		a.addFiles(fs)

		if !hasNext {
			break
		}
		i++
	}

	return nil
}

func (a *Album) getPhoto(albumId int, limit int, offset int) ([]*File, bool, error) {
	reply, err := a.getPhotoRequest(albumId, limit, offset)
	if err != nil {
		return nil, false, err
	}

	if limit+offset >= reply.Value.Count {
		return reply.Value.File, false, nil
	}
	return reply.Value.File, true, nil
}

func (a *Album) getPhotoRequest(albumId int, limit int, offset int) (*ListReply, error) {
	q := &url.Values{}
	q.Set("limit", strconv.Itoa(limit))
	q.Set("offset", strconv.Itoa(offset))
	q.Set("order", "1")
	q.Set("dirId", strconv.Itoa(albumId))
	q.Set("token", a.token)

	req, err := http.NewRequest(http.MethodPost, "https://mzstorage.meizu.com/album/list", bytes.NewBufferString(q.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("album/list dirId: %d, offset: %d, limit:%d, body: %s\n", albumId, offset, limit, string(body))

	reply := &ListReply{}
	err = json.Unmarshal(body, reply)

	return reply, err
}

func (a *Album) addFiles(fs []*File) {
	if a.Files == nil {
		a.Files = make([]*File, 0)
	}
	a.Files = append(a.Files, fs...)
}
