package model

import (
	"encoding/json"
	"math"
)

type Currency struct {
	tableName struct{} `pg:"currencies,alias:currencies"`
	BaseModel
	Rate float32 `json:"rate", pg:"rate,default:0"`
	Name string  `json:"name", pg:"name,unique"`
}

type CurrencyEdit struct {
	tableName struct{}    `pg:"currencies,alias:currencies"`
	Id        string      `json:"id" form:"id"`
	Rate      json.Number `json:"rate", form:"rate"`
	Name      string      `json:"name", form:"name"`
}

type NewCurrencyBody struct {
	Rate json.Number `json:"rate", form:"rate"`
	Name string      `json:"name", form:"name"`
}

type ExchangeBody struct {
	Base  string      `json:"base"`
	Rates interface{} `json:"rates"`
}

type Rate struct {
	Code string  `json:"code"`
	Rate float64 `json:"rate"`
}

func (body *CurrencyEdit) ToCurrency() *Currency {
	rate, _ := body.Rate.Float64()
	tm := new(Currency)
	tm.Id = body.Id
	tm.Rate = float32(math.Round(rate*100) / 100)
	tm.Name = body.Name

	return tm
}

func (body *NewCurrencyBody) ToCurrency() *Currency {
	rate, _ := body.Rate.Float64()
	tm := new(Currency)
	tm.Init()
	tm.Rate = float32(math.Round(rate*100) / 100)
	tm.Name = body.Name
	return tm
}

func (body *ExchangeBody) Unmarshal(resp *[]interface{}) *ExchangeBody {
	body.Base = (*resp)[3].(string)
	// body.Rates = []Rate{}
	// for k, v := range (*resp)[5].(map[string]interface{}) {
	// 	body.Rates = append(body.Rates, Rate{Code: k, Rate: v.(float64)})
	// }
	return body
}
