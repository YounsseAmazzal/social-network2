package requests

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"social/internal/models"
	"social/internal/utils"
	"strings"
)

func PostPyloadParse(req *http.Request) (models.PostCreation, error) {
	var postPyyload models.PostCreation
	contentType := req.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := req.ParseMultipartForm(32 << 20); err != nil {
			log.Printf("Error parsing multipart form on poooost ")
			return postPyyload, fmt.Errorf("Error parsing multipart form: %v", err)
		}

		postPyyload.Title = req.FormValue("PostTitle")
		postPyyload.Description = req.FormValue("PostContent")
		postPyyload.Postprivecy = req.FormValue("privacy")
		//fmt.Println("!!!!!!!!!!THee Post Privercy!!!!!!!!!!!!1", postPyyload.Postprivecy)
		if postPyyload.Postprivecy == "followers_chosen" {
			//fmt.Println("EEEEEEEEEE@@@@#####")
			shared_ids := req.Form["followers_chosen"]
			for _, id := range shared_ids {
				//fmt.Println("THe shareed post", id)
				postPyyload.SharedWith = append(postPyyload.SharedWith, id)
			}
		}
		if file_image, header, err := req.FormFile("shared_image"); err == nil {
			defer file_image.Close()
			if header != nil || header.Size > 0 {
				Imagename, err := utils.ProcessImageUpload(file_image, header)
				if err != nil {
					return postPyyload, fmt.Errorf("failed to save post image file: %w", err)
				}
				postPyyload.Image = Imagename
			}
		}

	}
	// //fmt.Println("@@@@@@@@thee user retuerner ", postPyyload.DateOfBirth)
	return postPyyload, nil
}
func ValidPostPyload(paylod *models.PostCreation) error {
	if paylod.Title == "" {
		return errors.New("Title cannot be eempty")
	}
	if len(paylod.Title) > 50 {
		return errors.New("the title is longer then (50) char ")
	}
	if paylod.Description == "" {
		return errors.New("Description cannot be eempty")
	}
	if len(paylod.Description) > 300 {
		return errors.New("the title is longer then (300) char ")
	}
	if paylod.Postprivecy != "publicpost" && paylod.Postprivecy != "only_followers" && paylod.Postprivecy != "followers_chosen" {
		return errors.New("check your post shaaaring ")
	}
	if paylod.Postprivecy == "followers_chosen" && len(paylod.SharedWith) == 0 {
		return errors.New("You should choice people to share with them this possst hh")
	}
	if paylod.Image == "" {
		//fmt.Println("PPPPPPPPPPPPPPPPPPOSST PIIC", paylod.Image)
		return errors.New("Enter an image it's nooot option enter valid image !!! ")
	}
	return nil
}
func ExtractPosstid(request *http.Request) (string, error) {
	Postt_id := request.URL.Query().Get("Possst_id")
	if Postt_id == "" {
		return "", fmt.Errorf("Post Id not Foundd")
	}
	return Postt_id, nil
}
func Parsecommmenpyload(req *http.Request) (models.Comment, error) {
	var commentpyload models.Comment
	contentType := req.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := req.ParseMultipartForm(32 << 20); err != nil {
			log.Printf("Error parsing multipart form on comment ")
			return commentpyload, fmt.Errorf("Error parsing multipart form: %v", err)
		}

		commentpyload.Content = req.FormValue("Content")
		if file_image, header, err := req.FormFile("comment_immmage"); err == nil {
			defer file_image.Close()
			if header != nil || header.Size > 0 {
				Imagename, err := utils.ProcessImageUpload(file_image, header)
				if err != nil {
					return commentpyload, fmt.Errorf("failed to save post image file: %w", err)
				}
				commentpyload.Image = Imagename
			}
		}

	}
	//fmt.Println("Returneer payload",commentpyload)
	return commentpyload, nil
}
func ValidCommentpy(commentpayload *models.Comment) error {
	if commentpayload.Content == "" {
		return errors.New("contennt cannot be eempty")
	}
	if len(commentpayload.Content) > 1000 {
		return errors.New("the title is longer then (50) char ")
	}
	// if commentpayload.Image==""{
	// 	return errors.New("Select an Image and it should be valid!!!")
	// }
	return nil
}
