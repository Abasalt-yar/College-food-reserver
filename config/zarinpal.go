package config

import (
	"encoding/json"
)

type IZarinpal interface {
	DoRequest(jsonData string) error
}

type ZarinpalRequest struct {
	Data struct {
		Code      int
		Message   string
		Authority string
	}
}

type ZarinaplRedis struct {
	UserID uint
	Amount uint64
}

func (z *ZarinpalRequest) DoRequest(jsonData string) error {
	resp, err := DoPostRequest("https://sandbox.zarinpal.com/pg/v4/payment/request.json", jsonData)
	if err != nil {
		return err
	}
	if er := json.Unmarshal(resp, z); er != nil {
		return er
	}
	return nil

}

type ZarinpalVerify struct {
	Data struct {
		Code  int
		RefID uint64 `json:"ref_id"`
	}
}

func (z *ZarinpalVerify) DoRequest(jsonData string) error {
	res, err := DoPostRequest("https://sandbox.zarinpal.com/pg/v4/payment/verify.json", jsonData)
	if err != nil {
		return err
	}
	if er := json.Unmarshal(res, z); er != nil {
		return er
	}
	return nil
}
