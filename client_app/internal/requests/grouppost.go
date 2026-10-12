package requests

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"social/internal/models"
	"social/internal/utils"
	"strconv"
	"strings"
)

func ParseGroupPOst(req *http.Request) (models.GroupPostCreation, error) {
	var groupPostPyload models.GroupPostCreation
	contentType := req.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := req.ParseMultipartForm(32 << 20); err != nil {
			log.Printf("Error parsing multipart form on poooost ")
			return groupPostPyload, fmt.Errorf("Error parsing multipart form: %v", err)
		}
		groupPostPyload.GroupPost.Title = req.FormValue("GroupPostTitle")
		groupPostPyload.GroupPost.Description = req.FormValue("GroupPostContent")
		if file_image, header, err := req.FormFile("Group_post_image"); err == nil {
			defer file_image.Close()
			if header != nil || header.Size > 0 {
				Imagename, err := utils.ProcessImageUpload(file_image, header,"post")
				if err != nil {
					return groupPostPyload, fmt.Errorf("failed to save group post image file: %w", err)
				}
				groupPostPyload.GroupPost.Image = Imagename
			}
		}

	}
	return groupPostPyload, nil
}
func ValidGroupPost(paylod *models.GroupPostCreation) error {
	if paylod.GroupPost.Title == "" {
		return errors.New("Title cannot be eempty")
	}
	if len(paylod.GroupPost.Title) > 50 {
		return errors.New("the title is longer then (50) char ")
	}
	if paylod.GroupPost.Description == "" {
		return errors.New("Description cannot be eempty")
	}
	if len(paylod.GroupPost.Description) > 300 {
		return errors.New("the title is longer then (300) char ")
	}
	if paylod.GroupPost.Image == "" {
		//fmt.Println("PPPPPPPPPPPPPPPPPPOSST PIIC", paylod.Image)
		return errors.New("Enter an image it's nooot option enter valid image !!! ")
	}
	return nil
}
func GetTheLimitandOffset(request *http.Request) (models.GrouupPR) {
	 liimit := request.URL.Query().Get("limit")
    offsetstr := request.URL.Query().Get("offset")
    limit, err := strconv.Atoi(liimit)
    if err != nil {
        limit = 10 
    }
    offset, err := strconv.Atoi(offsetstr)
    if err != nil {
        offset = 0
    }
    grouppagination := models.GrouupPR{Limit: limit, Offset: offset}
    grouppagination.SetDefaults() 
	return grouppagination
}
func ValidGroupComment(paylod *models.GroupPostCreation) error {
	if paylod.GroupPost.Title == "" {
		return errors.New("Title cannot be eempty")
	}
	if len(paylod.GroupPost.Title) > 50 {
		return errors.New("the title is longer then (50) char ")
	}
	if paylod.GroupPost.Description == "" {
		return errors.New("Description cannot be eempty")
	}
	if len(paylod.GroupPost.Description) > 300 {
		return errors.New("the title is longer then (300) char ")
	}
	if paylod.GroupPost.Image == "" {
		//fmt.Println("PPPPPPPPPPPPPPPPPPOSST PIIC", paylod.Image)
		return errors.New("Enter an image it's nooot option enter valid image !!! ")
	}
	return nil
}
func ParseCommentpy(req *http.Request) (models.GroupPostCreation, error) {
	var groupPostPyload models.GroupPostCreation
	contentType := req.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := req.ParseMultipartForm(32 << 20); err != nil {
			log.Printf("Error parsing multipart form on poooost ")
			return groupPostPyload, fmt.Errorf("Error parsing multipart form: %v", err)
		}
		groupPostPyload.GroupPost.Title = req.FormValue("GroupCommentTitle")
		groupPostPyload.GroupPost.Description = req.FormValue("GroupPostContent")
		if file_image, header, err := req.FormFile("Group_Comment_image"); err == nil {
			defer file_image.Close()
			if header != nil || header.Size > 0 {
				Imagename, err := utils.ProcessImageUpload(file_image, header,"post")
				if err != nil {
					return groupPostPyload, fmt.Errorf("failed to save group post image file: %w", err)
				}
				groupPostPyload.GroupPost.Image = Imagename
			}
		}

	}
	return groupPostPyload, nil
}