package main

type ListReply struct {
	Code     int        `json:"code"`
	Message  string     `json:"message"`
	Value    *ListValue `json:"value"`
	Redirect string     `json:"redirect"`
}

type ListValue struct {
	Count int     `json:"count"`
	File  []*File `json:"file"`
	End   int     `json:"end"`
}

type File struct {
	Id              int    `json:"id"`
	FileName        string `json:"fileName"`
	DirName         string `json:"dirName"`
	DirId           int    `json:"dirId"`
	UserId          int    `json:"userId"`
	Size            int    `json:"size"`
	Url             string `json:"url"`
	Thumb256        string `json:"thumb256"`
	Thumb1024       string `json:"thumb1024"`
	ShootTime       int64  `json:"shootTime"`
	Md5             string `json:"md5"`
	Status          int    `json:"status"`
	IsVideo         bool   `json:"isVideo"`
	Height          int    `json:"height"`
	Width           int    `json:"width"`
	Tags            string `json:"tags"`
	TrashTime       int    `json:"trashTime"`
	CreateTime      int64  `json:"createTime"`
	ModifyTime      int64  `json:"modifyTime"`
	SqlNow          int64  `json:"sqlNow"`
	RemainTrashTime int    `json:"remainTrashTime"`
	DoubleCamera    int    `json:"doubleCamera"`
}
