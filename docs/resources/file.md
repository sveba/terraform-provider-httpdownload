---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "httpdownload_file Resource - terraform-provider-httpdownload"
subcategory: ""
description: |-
  Represents a download of a file from a HTTP server.
---

# httpdownload_file (Resource)

Represents a download of a file from a HTTP server.

## Example Usage

```terraform
resource "httpdownload_file" "hetzner100mb" {
  url  = "https://ash-speed.hetzner.com/100MB.bin"
  dest = "100MB.bin"
}

resource "httpdownload_file" "task" {
  url  = "https://github.com/go-task/task/releases/download/v3.42.0/task_linux_amd64.tar.gz"
  dest = "task.tar.gz"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `dest` (String) Destination file path.
- `url` (String) Url of the file to download.

### Read-Only

- `checksum` (String) Destination file path.
