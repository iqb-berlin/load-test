local uri = ngx.var.uri

local fileName
local token

token, wsId, fileName = string.match(uri, "^/file/([^/]+)/ws_(%d+)/Resource/(.+)$")

local wsRequested = tonumber(wsId)

local filePath = '/ws_' .. wsRequested .. '/Resource/' .. fileName
ngx.var.file_path = filePath
ngx.req.set_uri(filePath)

