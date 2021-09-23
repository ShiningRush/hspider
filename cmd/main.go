package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/shiningrush/goreq"
	"github.com/shiningrush/hspider/pkg/api"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	Date     []string          `json:"date"`
	Doctor   []*Doctor         `json:"doc"`
	Schedule map[int]*Schedule `json:"sch, omitempty"`
	Weeks    []string          `json:"week"`
}

type Doctor struct {
	Id   int    `json:"doctor_id"`
	Name string `json:"doctor_name"`
}

type Schedule struct {
	AM MapSchedule `json:"am"`
	PM MapSchedule `json:"pm"`
}

type MapSchedule map[string]*ScheduleDetail

func (s *MapSchedule) UnmarshalJSON(b []byte) error {
	rs := make(map[string]*ScheduleDetail)
	strb := string(b)
	if strb != "[]" {
		if strings.HasPrefix(strb, "[") {
			var arrRs []*ScheduleDetail
			if err := json.Unmarshal(b, &arrRs); err != nil {
				return err
			}
			rs["0"] = arrRs[0]
		} else {
			if err := json.Unmarshal(b, &rs); err != nil {
				return err
			}
		}
	}

	*s = rs
	return nil
}

type ScheduleDetail struct {
	DoctorName string `json:"doctor_name"`
	DoctorId   string `json:"doctor_id"`
	Date       string `json:"to_date"`
	// 1:ok -1:no
	State     string `json:"y_state"`
	StateDesc string `json:"y_state_desc"`
}

var (
	baseUrl    = "https://www.91160.com/dep/getschmast/uid-%d/depid-%d/date-%s/p-0.html"
	baseDocUrl = "https://www.91160.com/doctors/index/unit_id-%d/dep_id-%d/docid-%s.html"
	uid        = flag.Int("uid", 22, "医院ID")
	depId      = flag.Int("dep_id", 126, "科室ID")
	epWeeks    = flag.String("ep_weeks", "", "期望的工作日, 星期1-7，逗号分割")
	// 星期一: 1 ~ 7
	epDocIds = flag.String("ep_doc", "541,551,200224658", "期望的医生ID, 逗号分隔")
)

func main() {
	api.NewLogin().ConfigRoutes(restful.DefaultContainer)
	log.Fatal(http.ListenAndServe(":9080", nil))

	//flag.Parse()
	//expectedDocId := strings.Split(*epDocIds, ",")
	//expectedWeek := strings.Split(*epWeeks, ",")
	//fmt.Printf("ready for seach\n uid: %d\n depId: %d\n epWeeks:%s\n epDocIds:%s\n", *uid, *depId, *epWeeks, *epDocIds)
	//
	//isFinded := false
	//for !isFinded {
	//	url := fmt.Sprintf(baseUrl, *uid, *depId, time.Now().Format("2006-01-02"))
	//	if searchResult(url, expectedWeek, expectedDocId) {
	//		isFinded = true
	//	}
	//
	//	url = fmt.Sprintf(baseUrl, *uid, *depId, time.Now().Add(time.Hour * 24 * 7).Format("2006-01-02"))
	//	if searchResult(url, expectedWeek, expectedDocId) {
	//		isFinded = true
	//	}
	//
	//	if !isFinded {
	//		log.Println("no info wait a seconds and restart...")
	//		time.Sleep(20 * time.Second)
	//	}
	//}
}

func searchResult(url string, epWeek, epDocId []string) bool {
	var rs Response
	err := goreq.Get(url, goreq.SetHeader(http.Header{
		"User-Agent": []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36"},
	}), goreq.JsonResp(rs)).Do()
	if err != nil {
		log.Fatal(err)
	}
	epWeekIdx := getExpectedWeekIdx(epWeek, rs.Weeks)
	for _, schs := range rs.Schedule {

		am, pm := checkHasSchedule(epWeekIdx, epDocId, schs.AM), checkHasSchedule(epWeekIdx, epDocId, schs.PM)
		if am != nil {
			notify(fmt.Sprintf("find docname: %s, date: %s, %s \n", am.DoctorName, am.Date, "上午"), am.DoctorId)
			return true
		}

		if pm != nil {
			notify(fmt.Sprintf("find docname: %s, date: %s, %s \n", pm.DoctorName, pm.Date, "下午"), pm.DoctorId)
			return true
		}
	}

	return false
}

func getExpectedWeekIdx(expectedWeek, weeks []string) string {
	var epIdx []string
	for i := range weeks {
		for j := range expectedWeek {
			if weeks[i] == expectedWeek[j] {
				epIdx = append(epIdx, strconv.Itoa(i))
			}
		}
	}

	return strings.Join(epIdx, ",")
}

func findDoctor(docId int, docs []*Doctor) *Doctor {
	for _, v := range docs {
		if docId == v.Id {
			return v
		}
	}

	return &Doctor{Name: "未找到"}
}

func checkHasSchedule(dateIdx string, epDocId []string, schs MapSchedule) *ScheduleDetail {
	for k, v := range schs {
		if len(epDocId) > 0 {
			isInRange := false
			for i := range epDocId {
				if epDocId[i] == v.DoctorId {
					isInRange = true
				}
			}

			if !isInRange {
				continue
			}
		}

		if dateIdx != "" && !strings.Contains(dateIdx, k) {
			continue
		}

		if v.State == "1" {
			return v
		}
	}

	return nil
}

func notify(text, docId string) {
}
