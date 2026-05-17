package posts

func CanDeletePost(currentRole string) bool {
	if currentRole == "admin" || currentRole == "moderator" {
		return true
	}
	return false
}

func CanApprovePost(currentRole string) bool {
	if currentRole == "admin" || currentRole == "moderator" {
		return true
	}
	return false
}

func CanEditPost(currentRole string, currentUserId, authorId int) bool {
	if currentRole == "admin" || currentRole == "moderator" {
		return true
	}
	return currentUserId == authorId
}