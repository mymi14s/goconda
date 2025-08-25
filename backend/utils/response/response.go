package response

import "github.com/beego/beego/v2/server/web/context"

type JSON map[string]interface{}

func JSONOK(ctx *context.Context, data interface{}) {
    ctx.Output.SetStatus(200)
    ctx.Output.JSON(JSON{
        "success": true,
        "data":    data,
    }, false, false)
}

func JSONError(ctx *context.Context, code int, msg string) {
    ctx.Output.SetStatus(code)
    ctx.Output.JSON(JSON{
        "success": false,
        "error":   msg,
    }, false, false)
}
