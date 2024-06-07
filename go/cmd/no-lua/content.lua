ngx.header["X-Source"] = "Disk-nolua"
ngx.exec(ngx.var.file_path)
