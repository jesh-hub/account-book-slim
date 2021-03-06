package slack

import (
	"strconv"
	"strings"
	"time"

	"github.com/cloudflare/ahocorasick"
)

const (
	DEFAULT_HISTORY_LATEST = ""
	DEFAULT_HISTORY_OLDEST = "0"
	DEFAULT_HISTORY_Limit  = 100
	CLIENT_TIMEZONE        = "Asia/Seoul"
	TIME_FORMAT_YYYY_MM    = "2006-01"
)

type History struct {
	Messages []*Message `json:"messages"`
}

type HistoryParameters struct {
	Latest string `json:"latest"`
	Oldest string `json:"oldest"`
	Limit  int    `json:"limit"`
}

type Message struct {
	ClientMsgID string `json:"client_msg_id,omitempty"`
	Text        string `json:"text"`
	Ts          string `json:"ts"`
}

type MessageParameters struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func (mp *MessageParameters) StartAsTime(location string) time.Time {
	// Slack 앱 사용자의 Local 타임존을 파라미터로 받기 때문에 UTC 타임존으로 변환
	// 예시: 5월 1일 00시 ~ 09시에 입력한 채팅은 UTC 기준으로 4월 30일이기 때문에 GetMessages() 험수에서 누락됨
	var startTime time.Time
	var err error

	if len(location) > 0 {
		loc, _ := time.LoadLocation(location)
		startTime, err = time.ParseInLocation(TIME_FORMAT_YYYY_MM, mp.Start, loc)
	} else {
		startTime, err = time.Parse(TIME_FORMAT_YYYY_MM, mp.Start)
	}

	errorHandler(err)

	return startTime.UTC()
}

func (mp *MessageParameters) EndAsTime(location string) time.Time {
	var endTime time.Time
	var err error

	if len(location) > 0 {
		loc, _ := time.LoadLocation(location)
		endTime, err = time.ParseInLocation(TIME_FORMAT_YYYY_MM, mp.End, loc)
	} else {
		endTime, err = time.Parse(TIME_FORMAT_YYYY_MM, mp.End)
	}
	errorHandler(err)

	return endTime.UTC()
}

func NewHistoryParameters() HistoryParameters {
	return HistoryParameters{
		Latest: DEFAULT_HISTORY_LATEST,
		Oldest: DEFAULT_HISTORY_OLDEST,
		Limit:  DEFAULT_HISTORY_Limit,
	}
}

func (mp *MessageParameters) dateFilter(date string) bool {
	if len(mp.Start) > 0 && len(mp.End) > 0 {
		// 날짜 범위 파라미터 있을 경우 해당 범위만 조회
		dateTime, err := time.Parse("2006-01-02", date)
		errorHandler(err)

		if dateTime.Unix() >= mp.StartAsTime("").Unix() &&
			dateTime.Unix() < mp.EndAsTime("").AddDate(0, 1, 0).Unix() {
			return true
		} else {
			return false
		}
	} else {
		// 날짜 범위 파라미터 없을 경우 전체 조회
		return true
	}
}

func (s *SlackClient) GetMessages(messageParameters *MessageParameters) []*Message {
	// url 파라미터 설정
	historyParameters := NewHistoryParameters()
	if len(messageParameters.Start) > 0 && len(messageParameters.End) > 0 {
		historyParameters.Oldest = strconv.FormatInt(messageParameters.StartAsTime(CLIENT_TIMEZONE).Unix(), 10)
		// 늦게 입력된 채팅 크롤링을 위해서 end + 1달 처리
		historyParameters.Latest = strconv.FormatInt(messageParameters.EndAsTime(CLIENT_TIMEZONE).AddDate(0, 1, 0).Unix(), 10)
	}

	// Slack API 통신
	var history *History
	s.NewAPI("https://slack.com/api/conversations.history", historyParameters, &history)

	return history.Messages
}

func (s *SlackClient) FilterMessages(messages []*Message) []*Message {
	mentionFilter := "<@" + s.botId + "> "
	m := ahocorasick.NewStringMatcher([]string{mentionFilter})

	var messagesFiltered []*Message

	for _, message := range messages {
		// bot이 멘션된 채팅만 필터링
		if len(m.Match([]byte(message.Text))) == 1 {
			// 멘션 텍스트 제거
			message.Text = strings.Replace(message.Text, mentionFilter, "", 1)
			messagesFiltered = append(messagesFiltered, message)
		}
	}

	return messagesFiltered
}

func (s *SlackClient) ConvertToPayment(messagesFiltered []*Message, messageParameters *MessageParameters) []*Payment {
	trim := func(s string) string {
		return strings.Trim(s, " ")
	}

	var payments []*Payment

	for _, message := range messagesFiltered {
		txtSlice := strings.Split(message.Text, ";")
		if len(txtSlice) >= 6 {
			date := parseDate(trim(txtSlice[0]))
			if messageParameters.dateFilter(date) {
				price := parsePrice(trim(txtSlice[4]))
				monthlyInstallment, _ := strconv.Atoi(trim(txtSlice[5]))
				payments = append(payments, &Payment{
					Date:               date,
					Method:             trim(txtSlice[1]),
					Category:           trim(txtSlice[2]),
					Name:               trim(txtSlice[3]),
					Price:              price,
					MonthlyInstallment: monthlyInstallment,
				})
			}
		}
	}

	return payments
}

func parseDate(date string) string {
	if strings.Index(date, `.`) > 0 {
		date = strings.Replace(date, `.`, `-`, -1)
	}

	return date
}

func parsePrice(price string) int {
	// 쉼표가 있으면 제거
	if strings.Index(price, ",") > 0 {
		price = strings.Replace(price, `,`, ``, -1)
	}

	result, err := strconv.Atoi(price)
	errorHandler(err)

	return result
}
