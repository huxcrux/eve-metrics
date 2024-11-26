package helpers

func GetIndustryActivityName(id int32) string {
	switch id {
	case 1:
		return "Manufacturing"
	case 3:
		return "Researching Time Efficiency"
	case 4:
		return "Researching Material Efficiency"
	case 5:
		return "Copying"
	case 7:
		return "Reverse Engineering"
	case 8:
		return "Invention"
	case 11:
		return "Reactions"
	default:
		return "Unknown"
	}
}
