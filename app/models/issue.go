package models

import (
	"errors"
	"modeltest/app/util"
	"regexp"
	"strconv"
	"time"

	"github.com/revel/revel"
)

type Issue struct {
	GorpController
}

func (issue *Issue) Create(issue_data *IssueData) error {
	var v revel.Validation
	issue_data.Validate(&v)
	if v.HasErrors() {
		return errors.New("Validate Error")
	}

	//gorp doesn't support time type. we use unix time on DB.
	issue_data.Created = time.Now().Unix()
	issue_data.Updated = time.Now().Unix()

	err := Txn.Insert(issue_data)
	if err != nil {
		return errors.New("System Error")
	}

	return nil
}

//common function for get issues
func (issue *Issue) GetIssueList(condition string) []IssueData {
	sql := "select * from issue " + condition
	rows, _ := Dbm.Select(IssueData{}, sql)
	issue_list := make([]IssueData, len(rows))
	cnt := 0
	for _, row := range rows {
		issuedata := row.(*IssueData)
		issue_list[cnt].Id = issuedata.Id
		issue_list[cnt].Title = issuedata.Title
		issue_list[cnt].Source = issuedata.Source
		issue_list[cnt].Detail = issuedata.Detail
		issue_list[cnt].Priority = issuedata.Priority
		issue_list[cnt].Status = issuedata.Status
		issue_list[cnt].LimitStr = util.UnitTimeToDayString(issuedata.Limit)
		issue_list[cnt].CreatedStr = util.UnitTimeToDateString(issuedata.Created)
		issue_list[cnt].UpdatedStr = util.UnitTimeToDateString(issuedata.Updated)
		cnt++
	}
	return issue_list
}

//sql template
const SERVICE_ISSUE_SQL string = "SELECT" +
	" s.serviceid ServiceID, i.id IssueId, i.title IssueTitle, i.priority IssuePriority, s.status StatusCode, s.reflectdate ReflectDate" +
	" FROM service_issue s INNER JOIN issue i ON s.issueid = i.id"

//list issues of service
func (issue *Issue) GetServiceIssueList(serviceid int, status string) []ServiceIssueView {
	condition := " where s.serviceid=" + strconv.Itoa(serviceid)
	if status != "" {
		r := regexp.MustCompile("^[0-9]$")
		if !r.MatchString(status) {
			return nil
		}
		condition += " and s.status=" + status
	}

	sql := SERVICE_ISSUE_SQL + condition
	rows, err := Dbm.Select(ServiceIssueView{}, sql)
	if err != nil {
		panic(err)
		return nil
	}

	issue_list := make([]ServiceIssueView, len(rows))
	cnt := 0
	for _, row := range rows {
		issue_data := row.(*ServiceIssueView)
		issue_list[cnt].IssueId = issue_data.IssueId
		issue_list[cnt].IssueTitle = issue_data.IssueTitle
		issue_list[cnt].IssuePriorityStr = util.GetPriority(issue_data.IssuePriority)
		issue_list[cnt].StatusCode = issue_data.StatusCode
		issue_list[cnt].Status = util.GetStatus(issue_data.StatusCode)
		if issue_data.ReflectDate > 0 {
			issue_list[cnt].ReflectDateStr = util.UnitTimeToDayString(issue_data.ReflectDate)
		}

		cnt++
	}

	return issue_list
}
