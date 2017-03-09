// Copyright 2016 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package controllers

import (
	"encoding/base64"
	"errors"
	"net/http"
	"path"
	"strconv"

	"github.com/npiganeau/yep/yep/server"
	"github.com/npiganeau/yep/yep/tools/generate"
)

func CompanyLogo(c *server.Context) {
	c.File(path.Join(generate.YEPDir, "yep", "server", "static", "web", "src", "img", "logo.png"))
}

func Image(c *server.Context) {
	model := c.Query("model")
	field := c.Query("field")
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	uid := c.Session().Get("uid").(int64)
	img, gErr := getFieldValue(uid, id, model, field)
	res, err := base64.StdEncoding.DecodeString(img.(string))
	if err != nil || gErr != nil {
		c.Error(errors.New("Unable to fetch image"))
		return
	}
	c.Data(http.StatusOK, "image/png", res)
}