package handler

import (
	"net/http"
	"strconv"

	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/internal/service"
	"github.com/jampa_trip/pkg/database"
	"github.com/jampa_trip/pkg/middleware"
	"github.com/jampa_trip/pkg/util"
	"github.com/jampa_trip/pkg/webserver"
	"github.com/labstack/echo/v4"
)

type ImageHandler struct{}

// UploadImages - upload de múltiplas imagens
func (h ImageHandler) UploadImages(ctx echo.Context) error {

	userID := middleware.GetUserID(ctx)
	userType := middleware.GetUserType(ctx)

	if userType != "company" {
		return webserver.ErrorResponse(ctx, util.WrapError("Apenas empresas podem fazer upload de imagens", nil, http.StatusForbidden))
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return webserver.ErrorResponse(ctx, util.WrapError("Erro ao processar formulário", err, http.StatusBadRequest))
	}
	defer form.RemoveAll()

	files := form.File["images[]"]
	if len(files) == 0 {
		return webserver.ErrorResponse(ctx, util.WrapError("Nenhum arquivo enviado", nil, http.StatusBadRequest))
	}

	var request contract.UploadImagesRequest

	if tourIDStr := form.Value["tour_id"]; len(tourIDStr) > 0 && tourIDStr[0] != "" {
		tourID, err := strconv.Atoi(tourIDStr[0])
		if err != nil {
			return webserver.ErrorResponse(ctx, util.WrapError("tour_id deve ser um número válido", err, http.StatusBadRequest))
		}
		request.TourID = &tourID
	}

	if description := form.Value["description"]; len(description) > 0 {
		request.Description = description[0]
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	imageService := service.ImageServiceNew(database.DB)
	response, err := imageService.UploadImages(files, &request, userID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, response)
}

// ListImages - lista imagens do usuário
func (h ImageHandler) ListImages(ctx echo.Context) error {

	userID := middleware.GetUserID(ctx)

	request := &contract.ListImagesRequest{}
	if err := request.ParseQueryParams(
		ctx.QueryParam("tour_id"),
		ctx.QueryParam("format"),
		ctx.QueryParam("page"),
		ctx.QueryParam("limit"),
	); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	imageService := service.ImageServiceNew(database.DB)
	response, err := imageService.ListImages(request, userID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// DeleteImage - deleta uma imagem específica
func (h ImageHandler) DeleteImage(ctx echo.Context) error {

	imageID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return webserver.InvalidIDResponse(ctx, err)
	}

	if imageID < 1 {
		return webserver.ErrorResponse(ctx, util.WrapError("ID deve ser um número positivo", nil, http.StatusBadRequest))
	}

	userID := middleware.GetUserID(ctx)

	imageService := service.ImageServiceNew(database.DB)
	response, err := imageService.DeleteImage(imageID, userID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// UpdateImage - atualiza metadados de uma imagem
func (h ImageHandler) UpdateImage(ctx echo.Context) error {

	imageID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return webserver.InvalidIDResponse(ctx, err)
	}

	if imageID < 1 {
		return webserver.ErrorResponse(ctx, util.WrapError("ID deve ser um número positivo", nil, http.StatusBadRequest))
	}

	request := &contract.UpdateImageRequest{}
	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	userID := middleware.GetUserID(ctx)

	imageService := service.ImageServiceNew(database.DB)
	response, err := imageService.UpdateImage(imageID, request, userID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// ReorderImages - reordena imagens de um passeio
func (h ImageHandler) ReorderImages(ctx echo.Context) error {

	request := &contract.ReorderImagesRequest{}
	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	userID := middleware.GetUserID(ctx)

	imageService := service.ImageServiceNew(database.DB)
	response, err := imageService.ReorderImages(request, userID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// GetImageInfo - obtém informações detalhadas de uma imagem
func (h ImageHandler) GetImageInfo(ctx echo.Context) error {

	imageID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return webserver.InvalidIDResponse(ctx, err)
	}

	if imageID < 1 {
		return webserver.ErrorResponse(ctx, util.WrapError("ID deve ser um número positivo", nil, http.StatusBadRequest))
	}

	userID := middleware.GetUserID(ctx)

	imageService := service.ImageServiceNew(database.DB)
	response, err := imageService.GetImageInfo(imageID, userID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// BatchDeleteImages - deleta múltiplas imagens
func (h ImageHandler) BatchDeleteImages(ctx echo.Context) error {

	request := &contract.BatchDeleteImagesRequest{}
	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	userID := middleware.GetUserID(ctx)

	imageService := service.ImageServiceNew(database.DB)
	response, err := imageService.BatchDeleteImages(request, userID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}
