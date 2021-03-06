package models

import (
	"regexp"
	. "security-cop/app/util"
	"strconv"
	"time"
)

type Issue struct {
	GorpController
}

func (issue *Issue) Create(issue_data *IssueData) error {
	if err := issue_data.Validate(); err != nil {
		return err
	}

	issue_data.Limit = DayStringToUnixTime(issue_data.LimitStr)
	//gorp doesn't support time type. we use unix time on DB.
	issue_data.Created = time.Now().Unix()
	issue_data.Updated = time.Now().Unix()

	err := Txn.Insert(issue_data)
	if err != nil {
		return err
	}

	return nil
}

func (issue *Issue) Update(issue_data *IssueData) error {
	if err := issue_data.Validate(); err != nil {
		return err
	}

	issue_data.Updated = time.Now().Unix()
	_, err := Txn.Update(issue_data)
	if err != nil {
		return err
	}

	return nil
}

func (issue *Issue) GetListAll() []IssueData {
	return issue.getList("")
}

func (issue *Issue) GetList(status, priority string) []IssueData {
	conditionList := []string{}
	if status != "" {
		if _, err := strconv.Atoi(status); err == nil {
			conditionList = append(conditionList, "status="+status)
		}
	}
	if priority != "" {
		if _, err := strconv.Atoi(priority); err == nil {
			conditionList = append(conditionList, "priority="+priority)
		}
	}

	if len(conditionList) > 0 {
		condition := "where "
		for i := 0; i < len(conditionList); i++ {
			if i == 0 {
				condition += conditionList[i]
			} else {
				condition += " and " + conditionList[i]
			}
		}
		return issue.getList(condition)
	} else {
		return issue.getList("")
	}
}

func (issue *Issue) GetByID(id int) []IssueData {
	return issue.getList("where id=" + strconv.Itoa(id))
}

//common function for get issues
func (issue *Issue) getList(condition string) []IssueData {
	sql := "select * from issue " + condition
	rows, _ := Dbm.Select(IssueData{}, sql)
	issue_list := make([]IssueData, len(rows))
	cnt := 0
	for _, row := range rows {
		issuedata := row.(*IssueData)
		issue_list[cnt].ID = issuedata.ID
		issue_list[cnt].Title = issuedata.Title
		issue_list[cnt].Source = issuedata.Source
		issue_list[cnt].Detail = issuedata.Detail
		issue_list[cnt].Priority = issuedata.Priority
		issue_list[cnt].Status = issuedata.Status
		issue_list[cnt].LimitStr = UnixTimeToDayString(issuedata.Limit)
		issue_list[cnt].CreatedStr = UnixTimeToDateString(issuedata.Created)
		issue_list[cnt].UpdatedStr = UnixTimeToDateString(issuedata.Updated)
		cnt++
	}
	return issue_list
}

//list issues of service
func (issue *Issue) GetServiceIssueList(serviceid int, status string) []ServiceIssueView {
	sql_fmt := "SELECT" +
		" s.serviceid ServiceID, i.id IssueId, i.title IssueTitle, i.priority IssuePriority, s.status StatusCode, s.reflectdate ReflectDate" +
		" FROM service_issue s INNER JOIN issue i ON s.issueid = i.id"
	condition := " where s.serviceid=" + strconv.Itoa(serviceid)
	if status != "" {
		r := regexp.MustCompile("^[0-9]$")
		if !r.MatchString(status) {
			return nil
		}
		condition += " and s.status=" + status
	}

	sql := sql_fmt + condition
	rows, err := Dbm.Select(ServiceIssueView{}, sql)
	if err != nil {
		panic(err)
		return nil
	}

	issue_list := make([]ServiceIssueView, len(rows))
	cnt := 0
	for _, row := range rows {
		issue_data := row.(*ServiceIssueView)
		issue_list[cnt].IssueID = issue_data.IssueID
		issue_list[cnt].IssueTitle = issue_data.IssueTitle
		issue_list[cnt].IssuePriorityStr = GetPriority(issue_data.IssuePriority)
		issue_list[cnt].StatusCode = issue_data.StatusCode
		issue_list[cnt].Status = GetStatus(issue_data.StatusCode)
		if issue_data.ReflectDate > 0 {
			issue_list[cnt].ReflectDateStr = UnixTimeToDayString(issue_data.ReflectDate)
		}

		cnt++
	}

	return issue_list
}
