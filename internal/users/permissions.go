package users

import "errors"

func CanDeleteUser(currentRole string, currentUserId, profileId int) bool {
	if currentRole == "admin" {
		return currentUserId != profileId
	}
	return currentUserId == profileId
}

func CanBanUser(currentRole string, currentUserId, profileId int) bool {
	return currentRole == "admin" && currentUserId != profileId
}

func SelfDeletionDetected(requestedId, currentId int) (int, int, error) {
	if requestedId == currentId {
		return requestedId, currentId, errors.New("SELF_DELETION")
	}
	return requestedId, currentId, nil
}
