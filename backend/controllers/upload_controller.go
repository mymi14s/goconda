package controllers

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/beego/beego/v2/server/web"
)

type UploadController struct {
    BaseController
}

func (c *UploadController) Prepare() {
    c.MustAuth()
}

// @router /api/v1/upload [post]
func (c *UploadController) Upload() {
    f, h, err := c.GetFile("file")
    if err != nil {
        c.JSONError(400, "file is required")
        return
    }
    defer f.Close()

    uploadDir := web.AppConfig.DefaultString("upload::dir", "./uploads")
    if err := os.MkdirAll(uploadDir, 0o755); err != nil {
        c.JSONError(500, "could not create upload dir")
        return
    }

    name := filepath.Base(h.Filename)
    name = strings.ReplaceAll(name, "..", "")
    ts := time.Now().UnixNano()
    dst := filepath.Join(uploadDir, fmt.Sprintf("%d_%s", ts, name))

    out, err := os.Create(dst)
    if err != nil {
        c.JSONError(500, "could not save file")
        return
    }
    defer out.Close()

    if _, err := io.Copy(out, f); err != nil {
        c.JSONError(500, "could not write file")
        return
    }

    c.JSONOK(map[string]interface{}{
        "filename": h.Filename,
        "stored_as": dst,
        "size": h.Size,
    })
}
