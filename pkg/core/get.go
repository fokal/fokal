package core

import (
	"log"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/sprioc/conductor/pkg/model"
	"github.com/sprioc/conductor/pkg/rsp"
	"github.com/sprioc/conductor/pkg/store"
)

func GetUser(ref model.DBRef) (model.User, rsp.Response) {

	if strings.Compare(ref.Collection, "users") != 0 {
		return model.User{}, rsp.Response{Message: "Ref is of the wrong collection type",
			Code: http.StatusBadRequest}
	}

	var user = model.User{}

	err := store.Get(ref, &user)
	if err != nil {
		return model.User{}, rsp.Response{Message: "User not found",
			Code: http.StatusNotFound}
	}

	return user, rsp.Response{Code: http.StatusOK}
}

func GetImage(ref model.DBRef) (model.Image, rsp.Response) {
	if strings.Compare(ref.Collection, "images") != 0 {
		return model.Image{}, rsp.Response{Message: "Ref is of the wrong collection type",
			Code: http.StatusBadRequest}
	}

	var image model.Image

	err := store.Get(ref, &image)
	if err != nil {
		return model.Image{}, rsp.Response{Message: "Image not found",
			Code: http.StatusNotFound}
	}

	return image, rsp.Response{Code: http.StatusOK}
}

func GetCollection(ref model.DBRef) (model.Collection, rsp.Response) {
	if strings.Compare(ref.Collection, "collections") != 0 {
		return model.Collection{}, rsp.Response{Message: "Ref is of the wrong type",
			Code: http.StatusBadRequest}
	}

	var col model.Collection

	err := store.Get(ref, &col)
	if err != nil {
		return model.Collection{}, rsp.Response{Message: "Collection not found",
			Code: http.StatusNotFound}
	}

	return col, rsp.Response{Code: http.StatusOK}
}

func GetCollectionImages(ref model.DBRef) ([]*model.Image, rsp.Response) {
	if strings.Compare(ref.Collection, "collections") != 0 {
		return []*model.Image{}, rsp.Response{Message: "Ref is of the wrong type",
			Code: http.StatusBadRequest}
	}

	var images []*model.Image

	log.Printf("%+v", bson.M{"collections": ref})

	err := store.GetAll("images", bson.M{"collections": ref}, &images)
	if err != nil {
		return []*model.Image{}, rsp.Response{Code: http.StatusInternalServerError}
	}

	if len(images) == 0 {
		return []*model.Image{}, rsp.Response{Code: http.StatusNotFound,
			Message: "Collection does not exist or has not uploaded any images."}
	}

	return images, rsp.Response{Code: http.StatusOK}
}

func GetUserImages(ref model.DBRef) ([]*model.Image, rsp.Response) {
	if strings.Compare(ref.Collection, "users") != 0 {
		return []*model.Image{}, rsp.Response{Message: "Ref is of the wrong type",
			Code: http.StatusBadRequest}
	}

	var images []*model.Image

	log.Printf("%+v", bson.M{"owner": ref})

	err := store.GetAll("images", bson.M{"owner": ref}, &images)
	if err != nil {
		return []*model.Image{}, rsp.Response{Code: http.StatusInternalServerError}
	}

	if len(images) == 0 {
		return []*model.Image{}, rsp.Response{Code: http.StatusNotFound,
			Message: "User does not exist or has not uploaded any images."}
	}

	return images, rsp.Response{Code: http.StatusOK}
}

func GetFeaturedImages() ([]*model.Image, rsp.Response) {
	var images []*model.Image

	err := store.GetAll("images", bson.M{"featured": true}, &images)
	if err != nil {
		return []*model.Image{}, rsp.Response{Code: http.StatusInternalServerError}
	}

	if len(images) == 0 {
		return []*model.Image{}, rsp.Response{Code: http.StatusNoContent,
			Message: "No featured images exist at this time."}
	}

	return images, rsp.Response{Code: http.StatusOK}
}