package contextutil

import (
	"net/http"

	reportModel "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/reportModule/model"
	"github.com/gin-gonic/gin"

)

func ExtractUserContext(c *gin.Context) (*reportModel.ContextUser, bool) {
	idValue, idExists := c.Get("id")
	roleIdValue, roleIdExists := c.Get("roleId")
	branchIdValue, branchIdExists := c.Get("branchId")

	if !idExists || !roleIdExists || !branchIdExists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "User ID, RoleID, Branch ID not found in request context.",
		})
		return nil, false
	}

	return &reportModel.ContextUser{
		ID:       idValue,
		RoleID:   roleIdValue,
		BranchID: branchIdValue,
	}, true
}
