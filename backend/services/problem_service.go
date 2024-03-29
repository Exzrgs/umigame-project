package services

import (
	"umigame-api/models"
	"umigame-api/myerrors"
	"umigame-api/repositories"

	"github.com/jinzhu/copier"
)

func (s *Service) GetProblemListService(page int) ([]models.ProblemBase, error) {
	problemList, err := repositories.SelectProblemList(s.db, page)
	if err != nil {
		return nil, err
	}
	if len(problemList) == 0 {
		err := myerrors.NoData.Wrap(ErrNoData, "0 problem found")
		return nil, err
	}

	problemIDs := make([]int, repositories.ProblemNumPerPage)
	for i, problem := range problemList {
		problemIDs[i] = problem.ID
	}
	activityList, err := repositories.SelectActivityList(s.db, s.userID, problemIDs)
	if err != nil {
		return nil, err
	}

	response := make([]models.ProblemBase, 0)
	idIndex := make(map[int]int, repositories.ProblemNumPerPage)
	for i, problem := range problemList {
		var problemBase models.ProblemBase
		if err := copier.Copy(&problemBase, &problem); err != nil {
			err = myerrors.TypeCastFailed.Wrap(err, "internal server error")
			return nil, err
		}
		problemBase.IsSolved, problemBase.IsLiked = false, false
		response = append(response, problemBase)

		idIndex[problem.ID] = i
	}
	for _, activity := range activityList {
		problemBase := &response[idIndex[activity.ProblemID]]
		problemBase.IsSolved = activity.IsSolved
		problemBase.IsLiked = activity.IsLiked
	}

	return response, nil
}

/*
TODO:IDも返すようにする
*/
func (s *Service) GetProblemDetailService(problemID int) (models.ProblemDetail, error) {
	problem, err := repositories.SelectProblem(s.db, problemID)
	if err != nil {
		return models.ProblemDetail{}, err
	}

	activity, err := repositories.SelectActivity(s.db, s.userID, problemID)
	if err != nil {
		return models.ProblemDetail{}, err
	}

	chats, err := repositories.SelectProblemChat(s.db, s.userID, problemID)
	if err != nil {
		return models.ProblemDetail{}, err
	}

	var response models.ProblemDetail
	if err := copier.Copy(&response, &problem); err != nil {
		err = myerrors.TypeCastFailed.Wrap(err, "internal server error")
		return models.ProblemDetail{}, err
	}
	if err := copier.Copy(&response, &activity); err != nil {
		err = myerrors.TypeCastFailed.Wrap(err, "internal server error")
		return models.ProblemDetail{}, err
	}
	for _, chat := range chats {
		var chatBase models.ChatBase
		if err := copier.Copy(&chatBase, &chat); err != nil {
			err = myerrors.TypeCastFailed.Wrap(err, "internal server error")
			return models.ProblemDetail{}, err
		}
		response.ChatList = append(response.ChatList, chatBase)
	}

	return response, nil
}

func (s *Service) PostProblemService(problem models.Problem) (models.Problem, error) {
	newProblem, err := repositories.InsertProblem(s.db, problem)
	if err != nil {
		return models.Problem{}, err
	}

	var response models.Problem

	if err := copier.Copy(&response, &newProblem); err != nil {
		err = myerrors.TypeCastFailed.Wrap(err, "internal server error")
		return models.Problem{}, err
	}

	return response, nil
}
