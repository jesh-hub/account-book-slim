package controller

import (
	"abs/model"
	"abs/service"
	"abs/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AddPaymentMethod
// @Summary Add paymentMethod
// @Description 결제수단 등록
// @Tags payment_method
// @Accept json
// @Produce json
// @Param groupId path string true "Group ID"
// @Param paymentMethodAdd body model.PaymentMethodAdd true "PaymentMethodAdd"
// @Success 200 {object} model.PaymentMethod
// @Router /v1/group/{groupId}/paymentMethod [post]
func AddPaymentMethod(c *gin.Context) {
	groupId := c.Param("groupId")
	paymentMethodAdd := &model.PaymentMethodAdd{}
	if err := c.ShouldBindJSON(paymentMethodAdd); err != nil {
		util.ErrorHandler(c, 400, err)
		return
	}

	paymentMethod, err := service.AddPaymentMethod(groupId, paymentMethodAdd)
	if err != nil {
		util.ErrorHandler(c, 500, err)
		return
	}
	c.JSON(http.StatusOK, paymentMethod)
}

// FindPaymentMethodByGroupId
// @Summary Find paymentMethod by group id
// @Description group id로 결제수단 조회
// @Tags payment_method
// @Accept json
// @Produce json
// @Param groupId path string true "Group ID"
// @Success 200 {array} model.PaymentMethod
// @Router /v1/group/{groupId}/paymentMethod [get]
func FindPaymentMethodByGroupId(c *gin.Context) {
	groupId := c.Param("groupId")

	paymentMethods, err := service.FindPaymentMethodByGroupId(groupId)
	if err != nil {
		util.ErrorHandler(c, 500, err)
		return
	}
	c.JSON(http.StatusOK, paymentMethods)
}

// UpdatePaymentMethod
// @Summary Update paymentMethod
// @Description 결제수단 수정
// @Tags payment_method
// @Accept json
// @Produce json
// @Param paymentMethodId path string true "PaymentMethod ID"
// @Param paymentMethodUpdate body model.PaymentMethodUpdate true "PaymentMethodUpdate"
// @Success 200 {array} model.PaymentMethod
// @Router /v1/paymentMethod/{paymentMethodId} [put]
func UpdatePaymentMethod(c *gin.Context) {
	paymentMethodId := c.Param("paymentMethodId")
	paymentMethodUpdate := &model.PaymentMethodUpdate{}

	if err := c.ShouldBindJSON(paymentMethodUpdate); err != nil {
		util.ErrorHandler(c, 400, err)
		return
	}

	updated, err := service.UpdatePaymentMethod(paymentMethodId, paymentMethodUpdate)
	if err != nil {
		util.ErrorHandler(c, 500, err)
		return
	}
	c.JSON(http.StatusOK, updated)
}
