package telegram

type WishlistState string

const (
	StateIdle WishlistState = ""
	StateWaitingWishlistAdd WishlistState = "waiting wishlist add"
	StateWaitingWishlistRemove WishlistState = "waiting wishlist remove"
)

type SessionState struct {
	WishlistState WishlistState
}
