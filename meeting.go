package ysx

import (
	"fmt"
	"strings"
	"time"

	"github.com/imroc/req"
)

type MeetingData struct {
	Id        string `json:"Id"`
	MeetingNo uint64 `json:"MeetingNo"`
}

type meeting struct {
	host      string
	apiKey    string
	apiSecret string
	apiEcid   string
}

type meetingResult struct {
	Code    int         `json:"Code"`
	Message string      `json:"Message"`
	Data    MeetingData `json:"Data"`
}

func newMeeting(host, apiKey, apiSecret, apiEcid string) *meeting {
	return &meeting{
		host:      host,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		apiEcid:   apiEcid,
	}
}

func (m *meeting) Create(startTime int64, topic, hostId string) (data *MeetingData, err error) {
	url := fmt.Sprintf("%s/v20/meeting/createScheduledMeeting", m.host)
	params := req.Param{
		"Topic":                 topic,
		"Agenda":                "",
		"Duration":              360,
		"UTCStartTime":          time.Unix(startTime, 0).In(time.Local).Format("2006-01-02 15:04:05"),
		"LocalStartTime":        time.Unix(startTime, 0).In(time.Local).UTC().Format("2006-01-02 15:04:05"),
		"HostId":                hostId,
		"OpenHostVideo":         true,
		"OpenParticipantsVideo": true,
	}

	return m.postReq(url, params)
}

func (m *meeting) Update(startTime int64, ID, topic, hostID string, participants []string) (data *MeetingData, err error) {
	url := fmt.Sprintf("%s/v20/meeting/update", m.host)
	params := req.Param{
		"id":                    ID,
		"Topic":                 topic,
		"Agenda":                "",
		"Duration":              360,
		"UTCStartTime":          time.Unix(startTime, 0).In(time.Local).Format("2006-01-02 15:04:05"),
		"LocalStartTime":        time.Unix(startTime, 0).In(time.Local).UTC().Format("2006-01-02 15:04:05"),
		"HostId":                hostID,
		"Participants":          strings.Join(participants, ","),
		"OpenHostVideo":         true,
		"OpenParticipantsVideo": true,
	}

	return m.postReq(url, params)
}

func (m *meeting) Delete(id, hostId string) (data *MeetingData, err error) {
	url := fmt.Sprintf("%s/v20/meeting/delete", m.host)
	params := req.Param{
		"Id":     id,
		"HostId": hostId,
	}

	return m.postReq(url, params)
}

func (m *meeting) End(id, hostId string) (data *MeetingData, err error) {
	url := fmt.Sprintf("%s/v20/meeting/end", m.host)
	params := req.Param{
		"Id":     id,
		"HostId": hostId,
	}

	return m.postReq(url, params)
}

func (m *meeting) Get(id string) (data *MeetingData, err error) {
	url := fmt.Sprintf("%s/v20/meeting/get", m.host)
	params := req.Param{
		"Id": id,
	}

	return m.postReq(url, params)
}

func (m *meeting) postReq(url string, params req.Param) (data *MeetingData, err error) {
	var (
		resp   *req.Resp
		result = new(meetingResult)

		header = req.Header{
			"Authorization": fmt.Sprintf("Bearer %s", getAPIToken(m.apiKey, m.apiSecret, m.apiEcid)),
			"Content-Type":  "application/json; charset=utf-8",
		}
	)

	if resp, err = req.Post(url, header, req.BodyJSON(params)); nil != err {
		return
	}

	if err = resp.ToJSON(result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}
