package controllers

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func controllerLog(c *gin.Context) *log.Entry {
	entryRaw, ok := c.Get("logger")
	if !ok {
		return log.NewEntry(log.StandardLogger())
	}

	entry, ok := entryRaw.(*log.Entry)
	if !ok {
		return log.NewEntry(log.StandardLogger())
	}

	return entry
}
