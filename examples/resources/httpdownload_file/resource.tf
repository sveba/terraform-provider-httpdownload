resource "httpdownload_file" "hetzner100mb" {
  url  = "https://ash-speed.hetzner.com/100MB.bin"
  dest = "100MB.bin"
}

resource "httpdownload_file" "task" {
  url  = "https://github.com/go-task/task/releases/download/v3.42.0/task_linux_amd64.tar.gz"
  dest = "task.tar.gz"
}