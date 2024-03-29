package testdata

import (
	"umigame-api/models"
)

var SelectActivityList_Basic = []models.Activity{
	{
		UserID:    1,
		ProblemID: 1,
		IsSolved:  true,
		IsLiked:   true,
	},
	{
		UserID:    1,
		ProblemID: 2,
		IsSolved:  true,
		IsLiked:   false,
	},
	{
		UserID:    1,
		ProblemID: 3,
		IsSolved:  false,
		IsLiked:   true,
	},
}

// var AddSolved_Basic = []models.Activity{
// 	{
// 		UserID:    1,
// 		ProblemID: 1,
// 		IsSolved:  true,
// 	},
// }

var ChangeLiked_Basic = []models.Activity{
	{
		UserID:    1,
		ProblemID: 1,
		IsSolved:  true,
		IsLiked:   false,
	},
	{
		UserID:    1,
		ProblemID: 1,
		IsSolved:  true,
		IsLiked:   true,
	},
}
