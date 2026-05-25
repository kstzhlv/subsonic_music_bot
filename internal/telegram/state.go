package telegram

type WishlistState string

const (
	StateIdle WishlistState = ""
	StateWaitingWishlistAdd WishlistState = "waiting wishlist add"
)

type SessionState struct {
	WishlistState WishlistState
}
