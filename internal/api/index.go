package api

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize/english"
	"github.com/gin-gonic/gin"

	"github.com/photoprism/photoprism/internal/acl"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/internal/get"
	"github.com/photoprism/photoprism/internal/i18n"
	"github.com/photoprism/photoprism/internal/photoprism"
	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/txt"
)

// StartIndexing indexes media files in the "originals" folder.
//
// POST /api/v1/index
func StartIndexing(router *gin.RouterGroup) {
	router.POST("/index", func(c *gin.Context) {
		s := Auth(c, acl.ResourcePhotos, acl.ActionUpdate)

		if s.Abort(c) {
			return
		}

		conf := get.Config()
		settings := conf.Settings()

		if !settings.Features.Library {
			AbortFeatureDisabled(c)
			return
		}

		start := time.Now()

		var f form.IndexOptions

		if err := c.BindJSON(&f); err != nil {
			AbortBadRequest(c)
			return
		}

		// Configure index options.
		path := conf.OriginalsPath()
		convert := settings.Index.Convert && conf.SidecarWritable()
		skipArchived := settings.Index.SkipArchived

		indOpt := photoprism.NewIndexOptions(filepath.Clean(f.Path), f.Rescan, convert, true, false, skipArchived)
		indOpt.SetUser(s.User())

		if len(indOpt.Path) > 1 {
			event.InfoMsg(i18n.MsgIndexingFiles, clean.Log(indOpt.Path))
		} else {
			event.InfoMsg(i18n.MsgIndexingOriginals)
		}

		ind := get.Index()
		lastRun, lastFound := ind.LastRun()
		indexStart := time.Now()

		// Start indexing.
		found, indexed := ind.Start(indOpt)

		// Only run purge and moments if necessary.
		forceUpdate := indOpt.Rescan || indexed > 0 || lastRun.IsZero()
		updateIndex := forceUpdate || len(found) != lastFound

		log.Infof("index: updated %s [%s]", english.Plural(indexed, "file", "files"), time.Since(indexStart))

		// Update index?
		if updateIndex {
			event.Publish("index.updating", event.Data{
				"uid":    indOpt.UID,
				"action": indOpt.Action,
				"step":   "folders",
			})

			RemoveFromFolderCache(entity.RootOriginals)

			event.Publish("index.updating", event.Data{
				"uid":    indOpt.UID,
				"action": indOpt.Action,
				"step":   "purge",
			})

			// Configure purge options.
			prgOpt := photoprism.PurgeOptions{
				Path:   filepath.Clean(f.Path),
				Ignore: found,
				Force:  forceUpdate,
			}

			// Start purging.
			prg := get.Purge()

			if files, photos, updated, err := prg.Start(prgOpt); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": txt.UpperFirst(err.Error())})
				return
			} else if updated > 0 {
				event.InfoMsg(i18n.MsgRemovedFilesAndPhotos, len(files), len(photos))
				forceUpdate = true
			}
		}

		// Update moments?
		if forceUpdate {
			event.Publish("index.updating", event.Data{
				"uid":    indOpt.UID,
				"action": indOpt.Action,
				"step":   "moments",
			})

			moments := get.Moments()

			if err := moments.Start(); err != nil {
				log.Warnf("moments: %s", err)
			}
		}

		elapsed := int(time.Since(start).Seconds())

		msg := i18n.Msg(i18n.MsgIndexingCompletedIn, elapsed)

		event.Success(msg)
		event.Publish("index.completed", event.Data{
			"uid":     indOpt.UID,
			"action":  indOpt.Action,
			"path":    path,
			"seconds": elapsed,
		})

		UpdateClientConfig()

		c.JSON(http.StatusOK, i18n.Response{Code: http.StatusOK, Msg: msg})
	})
}

// CancelIndexing stops indexing media files in the "originals" folder.
//
// DELETE /api/v1/index
func CancelIndexing(router *gin.RouterGroup) {
	router.DELETE("/index", func(c *gin.Context) {
		s := Auth(c, acl.ResourcePhotos, acl.ActionUpdate)

		if s.Abort(c) {
			return
		}

		conf := get.Config()

		if !conf.Settings().Features.Library {
			AbortFeatureDisabled(c)
			return
		}

		ind := get.Index()

		ind.Cancel()

		c.JSON(http.StatusOK, i18n.NewResponse(http.StatusOK, i18n.MsgIndexingCanceled))
	})
}
