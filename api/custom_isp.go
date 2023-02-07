package api

import (
	"errors"
	"strconv"
	"strings"
)

const FUNC_NAME_CUSTOM_ISP = "custom_isp"

type CustomIspData struct {
	Ipgroup string `json:"ipgroup"`
	Time    string `json:"time"`
	ID      int    `json:"id"`
	Comment string `json:"comment"`
	Name    string `json:"name"`
}

func (i *IKuai) ShowCustomIspByComment() (result []CustomIspData, err error) {
	param := struct {
		Type    string `json:"TYPE"`
		Limit   string `json:"limit"`
		OrderBy string `json:"ORDER_BY"`
		Order   string `json:"ORDER"`
	}{
		Type: "data",
	}
	req := CallReq{
		FuncName: FUNC_NAME_CUSTOM_ISP,
		Action:   "show",
		Param:    &param,
	}
	resp := CallResp{Data: &CallRespData{Data: &result}}
	err = postJson(i.client, i.baseurl+"/Action/call", &req, &resp)
	if err != nil {
		return
	}
	if resp.Result != 30000 {
		err = errors.New(resp.ErrMsg)
		return
	}
	return
}

func (i *IKuai) AddCustomIsp(name, ipgroup string) error {
	param := struct {
		Name    string `json:"name"`
		Ipgroup string `json:"ipgroup"`
		Comment string `json:"comment"`
	}{
		Name:    name,
		Ipgroup: ipgroup,
		Comment: COMMENT_IKUAI_BYPASS,
	}
	req := CallReq{
		FuncName: FUNC_NAME_CUSTOM_ISP,
		Action:   "add",
		Param:    &param,
	}
	resp := CallResp{}
	err := postJson(i.client, i.baseurl+"/Action/call", &req, &resp)
	if err != nil {
		return err
	}
	if resp.Result != 30000 {
		return errors.New(resp.ErrMsg)
	}
	return nil
}

func (i *IKuai) DelCustomIsp(id string) error {
	param := struct {
		Id string `json:"id"`
	}{
		Id: id,
	}
	req := CallReq{
		FuncName: FUNC_NAME_CUSTOM_ISP,
		Action:   "del",
		Param:    &param,
	}
	resp := CallResp{}
	err := postJson(i.client, i.baseurl+"/Action/call", &req, &resp)
	if err != nil {
		return err
	}
	if resp.Result != 30000 {
		return errors.New(resp.ErrMsg)
	}
	return nil
}

func (i *IKuai) DelIKuaiBypassCustomIsp() (err error) {
	data, err := i.ShowCustomIspByComment()
	if err != nil {
		return
	}
	var ids []string
	for _, d := range data {
		if d.Comment == COMMENT_IKUAI_BYPASS {
			ids = append(ids, strconv.Itoa(d.ID))
		}
	}
	id := strings.Join(ids, ",")
	err = i.DelCustomIsp(id)
	return
}
