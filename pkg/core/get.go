package core

import (
	"net/http"

	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/rsp"
	"github.com/sprioc/composer/pkg/sql"
)

func GetUser(ref model.Ref) (model.User, rsp.Response) {
	if ref.Collection != model.Users {
		return model.User{}, rsp.Response{Message: "Ref is of the wrong collection type",
			Code: http.StatusBadRequest}
	}

	user, err := sql.GetUser(ref.Shortcode)
	if err != nil {
		switch err.Error() {
		case "User not found.":
			return model.User{}, rsp.Response{Message: err.Error(), Code: http.StatusNotFound}
		default:
			return model.User{}, rsp.Response{Message: err.Error(), Code: http.StatusInternalServerError}
		}
	}
	return user, rsp.Response{Code: http.StatusOK}
}

func GetImage(ref model.Ref) (model.Image, rsp.Response) {
	if ref.Collection != model.Images {
		return model.Image{}, rsp.Response{Message: "Ref is of the wrong collection type",
			Code: http.StatusBadRequest}
	}

	image, err := sql.GetImage(ref.Shortcode)
	if err != nil {
		return model.Image{}, rsp.Response{Message: err.Error(),
			Code: http.StatusInternalServerError}
	}
	return image, rsp.Response{Code: http.StatusOK}
}
func GetUserFollowed(username string) ([]model.User, rsp.Response) {
	images, err := sql.GetUserFollowed(username)
	if err != nil {
		return []model.User{}, rsp.Response{Message: err.Error(),
			Code: http.StatusInternalServerError}
	}
	return images, rsp.Response{Code: http.StatusOK}
}
func GetUserFavorites(username string) ([]model.Image, rsp.Response) {
	images, err := sql.GetUserFavorites(username)
	if err != nil {
		return []model.Image{}, rsp.Response{Message: err.Error(),
			Code: http.StatusInternalServerError}
	}
	return images, rsp.Response{Code: http.StatusOK}
}
func GetUserImages(username string) ([]model.Image, rsp.Response) {
	images, err := sql.GetUserImages(username)
	if err != nil {
		return []model.Image{}, rsp.Response{Message: err.Error(),
			Code: http.StatusInternalServerError}
	}
	return images, rsp.Response{Code: http.StatusOK}
}
func GetRecentImages(limit int) ([]model.Image, rsp.Response) {
	images, err := sql.GetRecentImages(limit)
	if err != nil {
		return []model.Image{}, rsp.Response{Message: err.Error(),
			Code: http.StatusInternalServerError}
	}
	return images, rsp.Response{Code: http.StatusOK}
}
func GetFeaturedImages(limit int) ([]model.Image, rsp.Response) {
	images, err := sql.GetFeaturedImages(limit)
	if err != nil {
		return []model.Image{}, rsp.Response{Message: err.Error(),
			Code: http.StatusInternalServerError}
	}
	return images, rsp.Response{Code: http.StatusOK}
}
