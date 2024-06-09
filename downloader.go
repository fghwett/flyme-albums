package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type Downloader struct {
	savePath string
	dataPath string
	token    string

	c []*Album

	bucket *oss.Bucket
}

func NewDownloader(savePath string, token string) *Downloader {
	return &Downloader{
		savePath: savePath,
		dataPath: "./albums.json",
		token:    token,
	}
}

func (d *Downloader) Run() error {
	if _, err := os.Stat(d.dataPath); os.IsNotExist(err) {

		if err = d.GetAlbums(); err != nil {
			return fmt.Errorf("获取相册列表失败 %v", err)
		}

		if err = d.GetPhotos(); err != nil {
			return fmt.Errorf("获取相片失败 %v", err)
		}

		if err := d.Save(); err != nil {
			return fmt.Errorf("缓存数据失败 err: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("检查缓存数据是否存在失败 err: %s\n", err)
	} else {
		body, err := ioutil.ReadFile(d.dataPath)
		if err != nil {
			return fmt.Errorf("read data file error, err: %s\n", err)
		}
		a := &[]*Album{}
		if err := json.Unmarshal(body, a); err != nil {
			return fmt.Errorf("json unmarshal erorr, err: %s", err)
		}
		d.c = *a
	}
	log.Printf("一共%d张照片", d.Count())

	// 3. 获取oss token
	if err := d.InitOss(); err != nil {
		return fmt.Errorf("获取oss token失败 err: %s", err)
	}

	// 4. 下载图片
	if err := d.Download(); err != nil {
		return fmt.Errorf("下载图片失败 err: %s", err)
	}

	return nil
}

func (d *Downloader) InitOss() error {
	q := &url.Values{}
	q.Set("type", "2")
	q.Set("token", d.token)

	req, err := http.NewRequest(http.MethodPost, "https://mzstorage.meizu.com/file/get_sig", bytes.NewBufferString(q.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("file/get_sig body: %s", string(body))

	reply := &SigReply{}
	if err := json.Unmarshal(body, reply); err != nil {
		return err
	}
	client, err := oss.New(reply.Value.Endpoint, reply.Value.AccessKeyId, reply.Value.AccessKeySecret, oss.SecurityToken(reply.Value.SecurityToken))
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(reply.Value.Bucket)
	if err != nil {
		return err
	}
	d.bucket = bucket

	return nil
}

func (d *Downloader) Download() error {
	for _, album := range d.c {
		if err := d.downAlbum(album); err != nil {
			return fmt.Errorf("下载相册 %d-%s 失败 err: %v", album.Id, album.DirName, err)
		}
	}
	return nil
}

func (d *Downloader) downAlbum(album *Album) error {
	folderPath := filepath.Join(d.savePath, album.DirName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	for _, f := range album.Files {
		if err := d.downPhoto(album.DirName, f); err != nil {
			e := fmt.Errorf("down %s err: %v", f.Url, err)
			log.Printf("downAlbum：%s \n", e)
		}
		log.Printf("下载%s成功", folderPath)
	}

	return nil
}

func (d *Downloader) downPhoto(folderName string, f *File) error {
	filePath := filepath.Join(d.savePath, folderName, f.FileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := d.bucket.GetObjectToFile(f.Url, filePath); err != nil {
			return fmt.Errorf("下载照片失败 url: %s, filePath: %s, err: %s", f.Url, filePath, err)
		}
	} else if err != nil {
		return fmt.Errorf("检查文件是否存在失败 filePath: %s, err: %s", filePath, err)
	}

	return nil
}

func (d *Downloader) GetAlbums() error {
	q := &url.Values{}
	q.Set("limit", "100")
	q.Set("offset", "0")
	q.Set("order", "1")
	q.Set("token", d.token)

	req, err := http.NewRequest(http.MethodPost, "https://mzstorage.meizu.com/album/dir/list", bytes.NewBufferString(q.Encode()))
	if err != nil {
		return fmt.Errorf("new request album/dir/list error: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request album/dir/list error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("album/dir/list read body error: %s", err)
	}
	log.Printf("album/dir/list body: %s\n", string(body))

	var dirReply DirReply
	if err := json.Unmarshal(body, &dirReply); err != nil {
		return fmt.Errorf("album/dir/list json unmarshal error: %v", err)
	}

	if dirReply.Code != 200 {
		return fmt.Errorf("album/dir/list code error: %v", dirReply.Message)
	}
	d.c = dirReply.Value.Dir

	return nil
}

func (d *Downloader) Count() (count int) {
	for _, album := range d.c {
		count += len(album.Files)
	}
	return
}

func (d *Downloader) GetPhotos() error {
	for _, album := range d.c {
		if err := album.GetPhotos(d.token); err != nil {
			return fmt.Errorf("%d-%s 获取相片失败 err: %s", album.Id, album.DirName, err)
		}
	}
	return nil
}

func (d *Downloader) Save() error {
	body, err := json.Marshal(d.c)
	if err != nil {
		return fmt.Errorf("保存相册信息失败 err: %v", err)
	}

	if err = ioutil.WriteFile(d.dataPath, body, 0666); err != nil {
		return fmt.Errorf("写入文件失败 err: %v", err)
	}

	return nil
}
