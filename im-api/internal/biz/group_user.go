package biz

type GroupUserUseCase struct {
	groupUserRepo GroupUserRepo
}

func NewGroupUseCase(groupUserRepo GroupUserRepo) *GroupUserUseCase {
	return &GroupUserUseCase{
		groupUserRepo: groupUserRepo,
	}
}

type GroupUserRepo interface {
}
