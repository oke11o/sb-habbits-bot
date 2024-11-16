package service

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oke11o/sb-habits-bot/internal/model"
)

type userRepo interface {
	SaveUser(ctx context.Context, income model.User) (model.User, error)
	UserCount(ctx context.Context) (int, error)
}

type incomeRepo interface {
	SaveIncome(ctx context.Context, income model.IncomeRequest) (model.IncomeRequest, error)
}

func NewIncomeServce(userrepo userRepo, incomeRepo incomeRepo) *IncomeService {
	return &IncomeService{
		userRepo:   userrepo,
		incomeRepo: incomeRepo,
	}
}

type IncomeService struct {
	userRepo   userRepo
	incomeRepo incomeRepo
}

func (s *IncomeService) Income(ctx context.Context, requestID string, update tgbotapi.Update) (model.User, error) {
	incomeRequest, err := model.NewIncomeRequestFromTgUpdate(requestID, update)
	if err != nil {
		return model.User{}, fmt.Errorf("model.NewIncomeRequestFromTgUpdate(%s) err: %w", requestID, err)
	}
	incomeRequest, err = s.incomeRepo.SaveIncome(ctx, incomeRequest)
	if err != nil {
		return model.User{}, fmt.Errorf("userRepo.SaveIncome() err: %w", err)
	}

	user, err := model.NewUserFromTgUpdate(update)
	if err != nil {
		return model.User{}, fmt.Errorf("model.NewUserFromTgUpdate() err: %w", err)
	}

	user, err = s.userRepo.SaveUser(ctx, user)
	if err != nil {
		return model.User{}, fmt.Errorf("userRepo.SaveUser() err: %w", err)
	}

	return user, nil
}
