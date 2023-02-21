package menu

import (
	"cmd-v2ex/constant"
	"cmd-v2ex/handler"
	"cmd-v2ex/handler/initEnv"

	"github.com/AlecAivazis/survey/v2"
)

func Process() {
	action := ""
	prompt := &survey.Select{
		Message: `
大家一起摸鱼吧!!!`,
		Options: actions(),
	}
	survey.AskOne(prompt, &action, handler.GetSurveyIcon())

	// fmt.Printf(action)

	switch action {
	case constant.ActionEnv:
		initEnv.NewInfo().Process()
	}

}

func actions() []string {
	return []string{
		constant.ActionEnv,
	}
}
