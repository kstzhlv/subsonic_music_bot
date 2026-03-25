type User struct {
	ID int64
	TelegramID int64
}

type UserRepository interface {
	GetAll(ctx context.Context) ([]User, error)
}
