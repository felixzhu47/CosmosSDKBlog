package keeper

import (
	"encoding/binary"

	//"fmt"
	//"fmt"
	"strconv"
	"time"

	//"github.com/gen2brain/beeep"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/example/blog/x/blog/types"
)

// GetCommentCount get the total number of comment
func (k Keeper) GetCommentCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentCountKey))
	byteKey := types.KeyPrefix(types.CommentCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to iint64
		panic("cannot decode count")
	}

	return count
}

// SetCommentCount set the total number of comment
func (k Keeper) SetCommentCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentCountKey))
	byteKey := types.KeyPrefix(types.CommentCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendComment appends a comment in the store with a new id and update the count, it also appends a comment to the post list directly
func (k Keeper) AppendComment(
	ctx sdk.Context,
	id uint64,
	creator string,
	body string,
	postID uint64,
	timee string,
) uint64 {
	// Create the comment
	count := k.GetCommentCount(ctx)
	current_time := time.Now()
	var comment = types.Comment{
		Creator: creator,
		Id:      count,
		Body:    body,
		PostID:  postID,
		Time:    current_time.String(),
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	value := k.cdc.MustMarshalBinaryBare(&comment)
	store.Set(GetCommentIDBytes(comment.Id), value)

	//commentIDs := k.GetPost(ctx, comment.PostID).CommentIds
	//commentIDs = append(commentIDs, comment.Id)

	//fetch the the post with the corresponding postID
	creator_post := k.GetPost(ctx, comment.PostID).Creator
	Title_post := k.GetPost(ctx, comment.PostID).Title
	Body_post := k.GetPost(ctx, comment.PostID).Body
	commentInPost := k.GetPost(ctx, comment.PostID).Comments

	// create the comment that will show up in the post
	var new_commentInPost = types.CommentInPost{
		Creator: creator,
		Body:    body,
		Time:    current_time.String(),
	}
	commentInPost = append(commentInPost, &new_commentInPost)

	// create a new post by adding the comment to it
	var post = types.Post{
		Creator: creator_post,
		Id:      postID,
		Title:   Title_post,
		Body:    Body_post,
		//CommentIds: commentIDs,
		Comments: commentInPost,
	}

	store_post := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PostKey))
	value_post := k.cdc.MustMarshalBinaryBare(&post)
	store_post.Set(GetCommentIDBytes(comment.PostID), value_post)

	// Update comment count
	k.SetCommentCount(ctx, count+1)

	return count

}

// SetComment set a specific comment in the store
func (k Keeper) SetComment(ctx sdk.Context, comment types.Comment) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	b := k.cdc.MustMarshalBinaryBare(&comment)
	store.Set(GetCommentIDBytes(comment.Id), b)
}

// GetComment returns a comment from its id
func (k Keeper) GetComment(ctx sdk.Context, id uint64) types.Comment {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	var comment types.Comment
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetCommentIDBytes(id)), &comment)
	return comment
}

// HasComment checks if the comment exists in the store
func (k Keeper) HasComment(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	return store.Has(GetCommentIDBytes(id))
}

// GetCommentOwner returns the creator of the comment
func (k Keeper) GetCommentOwner(ctx sdk.Context, id uint64) string {
	return k.GetComment(ctx, id).Creator
}

// RemoveComment removes a comment from the store
func (k Keeper) RemoveComment(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	store.Delete(GetCommentIDBytes(id))
}

// GetAllComment returns all comment
func (k Keeper) GetAllComment(ctx sdk.Context) (list []types.Comment) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Comment
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetCommentIDBytes returns the byte representation of the ID
func GetCommentIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetCommentIDFromBytes returns ID in uint64 format from a byte array
func GetCommentIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
