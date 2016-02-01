package controllers

import (
	"security-cop/app/controllers"
	"security-cop/app/models"
	"security-cop/app/util"
	"strconv"
	"time"

	"github.com/revel/revel"
)

type ApiIssues struct {
	ApiController
}

func (c *ApiIssues) Create() revel.Result {
	issue := &models.Issue{}
	c.BindParams(issue)

	issue.Validate(c.App.Validation)
	if c.App.Validation.HasErrors() {
		return c.App.RenderJson(&ErrorResponse{ERR_VALIDATE, ErrorMessage(ERR_VALIDATE)})
	}

	issue.Created = time.Now().Unix()
	issue.Updated = time.Now().Unix()
	err := c.Txn.Insert(issue)
	if err != nil {
		panic(err)
	}

	return c.App.RenderJson(&Response{OK, issue})
}

func (c *ApiIssues) Update(id int) revel.Result {
	issue := &models.Issue{}
	c.BindParams(issue)
	issue.Id = id

	issue.Validate(c.App.Validation)
	if c.App.Validation.HasErrors() {
		return c.App.RenderJson(&ErrorResponse{ERR_VALIDATE, ErrorMessage(ERR_VALIDATE)})
	}

	_, err := c.Txn.Update(issue)
	if err != nil {
		panic(err)
	}

	return c.App.RenderJson(&Response{OK, issue})
}

func (c *ApiIssues) Show(id int) revel.Result {
	issues := getIssues("where id=" + strconv.Itoa(id))
	return c.App.RenderJson(&Response{OK, issues})
}

func (c *ApiIssues) List(q string) revel.Result {
	issues := getIssues("")
	return c.App.RenderJson(&Response{OK, issues})
}

func getIssues(condition string) []models.Issue {
	sql := "select * from issue " + condition
	rows, _ := controllers.Dbm.Select(models.Issue{}, sql)
	issues := make([]models.Issue, len(rows))
	cnt := 0
	for _, row := range rows {
		issue := row.(*models.Issue)
		issues[cnt].Id = issue.Id
		issues[cnt].Title = issue.Title
		issues[cnt].Source = issue.Source
		issues[cnt].Detail = issue.Detail
		issues[cnt].Priority = issue.Priority
		issues[cnt].Status = issue.Status
		issues[cnt].Limit = issue.Limit
		issues[cnt].CreatedStr = util.UnitTimeToString(issue.Created)
		issues[cnt].UpdatedStr = util.UnitTimeToString(issue.Updated)
		cnt++
	}
	return issues
}
