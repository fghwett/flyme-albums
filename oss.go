package main

import "time"

type SigReply struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Value   struct {
		Region          string    `json:"region"`
		EndpointImage   string    `json:"endpointImage"`
		SecurityToken   string    `json:"securityToken"`
		Bucket          string    `json:"bucket"`
		ExpiredTime     time.Time `json:"expiredTime"`
		AccessKeyId     string    `json:"accessKeyId"`
		RegionImage     string    `json:"regionImage"`
		Endpoint        string    `json:"endpoint"`
		AccessKeySecret string    `json:"accessKeySecret"`
	} `json:"value"`
	Redirect string `json:"redirect"`
}
