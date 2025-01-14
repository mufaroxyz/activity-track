package lib

import (
	"regexp"
)

/*
activity -> { Gaming, Coding, Youtube, Social, Discord, Other}
high_level_idenitifer -> by exe
if high_level_identifier has activity mapping -> {
	check if there are constraints for the activity (regex matching of window title) -> {
		if constraints are met -> activity
		else -> Other
	} else {
		activity determined by high_level_identifier
	}
} else {
	activity -> Other
}

"Brave.exe": {
	"activity": "Social",
	"constraints": [
		{"* - Youtube - Brave": "Youtube"},
	]
}
*/

type Constraint struct {
	Regex    string
	Activity string
}

type Association struct {
	Activity     string
	Constraints  []Constraint
	HighLevelExe interface{} // can be string or []string
}

var PredefinedActivities = map[string]string{
	"gaming":  "Gaming",
	"coding":  "Coding",
	"youtube": "Youtube",
	"social":  "Social",
	"discord": "Discord",
	"other":   "Other",
}

var PredefinedAssociations = []Association{
	{
		Activity: "social",
		Constraints: []Constraint{
			{
				Regex:    `\s*- YouTube - Brave$`,
				Activity: "youtube",
			},
		},
		HighLevelExe: "brave.exe",
	},
	{
		Activity:     "discord",
		HighLevelExe: []string{"discord.exe", "discordcanary.exe"},
	},
	{
		Activity: "coding",
		HighLevelExe: []string{
			"code.exe",
			"code - insiders.exe",
			"windowsterminal.exe",
			"powershell.exe",
			"goland64.exe",
			"idea64.exe",
			"zed.exe",
		},
	},
}

var IgnoreTitleList = []string{
	"Task Switching",
}

func IsTitleIgnored(title string) bool {
	for _, ignoredTitle := range IgnoreTitleList {
		if title == ignoredTitle {
			return true
		}
	}

	return false
}

func GetAssociation(fileExe string, windowTitle string) string {
	for _, association := range PredefinedAssociations {
		if association.HighLevelExe == fileExe {
			if len(association.Constraints) > 0 {
				for _, constraint := range association.Constraints {
					matched, _ := regexp.MatchString(constraint.Regex, windowTitle)
					if matched {
						return PredefinedActivities[constraint.Activity]
					}
				}
			}

			return PredefinedActivities[association.Activity]
		}

		if arr, ok := association.HighLevelExe.([]string); ok {
			for _, exe := range arr {
				if exe == fileExe {
					if len(association.Constraints) > 0 {
						for _, constraint := range association.Constraints {
							matched, _ := regexp.MatchString(constraint.Regex, windowTitle)
							if matched {
								return PredefinedActivities[constraint.Activity]
							}
						}
					}

					return PredefinedActivities[association.Activity]
				}
			}
		}
	}

	return PredefinedActivities["other"]
}
