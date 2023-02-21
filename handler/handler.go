package handler

import "github.com/AlecAivazis/survey/v2"

type Processor interface {
	Process()
}

func GetSurveyIcon() survey.AskOpt {
	return survey.WithIcons(func(icons *survey.IconSet) {
		icons.MarkedOption.Text = "[✔]"
		icons.UnmarkedOption.Text = "[✖]"
		icons.SelectFocus.Text = "➜"
		icons.UnmarkedOption.Format = "red"
		icons.Question.Text = `
cmd-v2ex(命令行浏览v2ex)
`
	})
}
