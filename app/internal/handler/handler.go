package handler

import (
	"net/http"
	"strconv"

	"comment-tree/internal/models"
	"comment-tree/internal/service"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type CommentsHandler struct {
	commService *service.CommentsService
	log         *zlog.Zerolog
}

func NewCommentsHandler(commService *service.CommentsService, log *zlog.Zerolog) *CommentsHandler {
	return &CommentsHandler{
		commService: commService,
		log:         log,
	}
}

func (h *CommentsHandler) Create(c *ginext.Context) {
	com, ok := h.getComment(c)
	if !ok {
		return
	}

	if err := h.commService.Create(c.Request.Context(), &com); err != nil {
		h.log.Error().
			Err(err).
			Msg("failed to create comment")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	h.log.Info().
		Int64("id", com.ID).
		Msg("comment created")
	c.JSON(http.StatusOK, com)
}

func (h *CommentsHandler) Update(c *ginext.Context) {
	com, ok := h.getComment(c)
	if !ok {
		return
	}

	if err := h.commService.Update(c.Request.Context(), &com); err != nil {
		h.log.Error().
			Err(err).
			Int64("id", com.ID).
			Msg("failed to update comment")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	h.log.Info().
		Int64("id", com.ID).
		Msg("comment updated")
	c.JSON(http.StatusOK, com)
}

func (h *CommentsHandler) Delete(c *ginext.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid id"})
		return
	}

	if err := h.commService.Delete(c.Request.Context(), id); err != nil {
		h.log.Error().
			Err(err).
			Int64("id", id).
			Msg("failed to delete comment")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
	}
}

func (h *CommentsHandler) GetByParent(c *ginext.Context) {
	parentStr := c.Query("parent")
	var parent *int64

	if parentStr == "null" || parentStr == "" {
		parent = nil
	} else {
		id, err := strconv.ParseInt(parentStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginext.H{
				"error": "invalid parent id",
			})
			return
		}
		parent = &id
	}

	limit, ok := h.getLimit(c)
	if !ok {
		return
	}
	offset, ok := h.getOffset(c)
	if !ok {
		return
	}

	coms, err := h.commService.GetByParent(c.Request.Context(), parent, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coms)
}

func (h *CommentsHandler) Search(c *ginext.Context) {
	query := c.Query("query")
	limit, ok := h.getLimit(c)
	if !ok {
		return
	}
	offset, ok := h.getOffset(c)
	if !ok {
		return
	}

	coms, err := h.commService.Search(c.Request.Context(), query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, coms)
}

func (h *CommentsHandler) RegisterRoutes(r *ginext.Engine) {
	g := r.Group("/comments")

	g.POST("/", h.Create)
	g.POST("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
	g.GET("/", h.GetByParent)
	g.GET("/search", h.Search)
}

func (h *CommentsHandler) getComment(c *ginext.Context) (models.Comment, bool) {
	var com models.Comment
	if err := c.ShouldBindJSON(&com); err != nil {
		h.log.Error().
			Err(err).
			Msg("failed to bind comment")
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return com, false
	}

	if err := models.Validate(&com); err != nil {
		h.log.Error().
			Err(err).
			Msg("invalid comment")
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return com, false
	}

	return com, true
}

func (h *CommentsHandler) getLimit(c *ginext.Context) (int64, bool) {
	limitStr := c.Query("limit")
	var limit int64
	var err error

	if limitStr == "" {
		limit = 10
	} else {
		limit, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginext.H{
				"error": "invalid limit",
			})
			return 0, false
		}
	}

	return limit, true
}

func (h *CommentsHandler) getOffset(c *ginext.Context) (int64, bool) {
	offsetStr := c.Query("offset")
	var offset int64
	var err error

	if offsetStr == "" {
		offset = 0
	} else {
		offset, err = strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginext.H{
				"error": "invalid offset",
			})
			return 0, false
		}
	}

	return offset, true
}
