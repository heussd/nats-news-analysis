group "default" {
  targets = [
    "feed-feeder",
    "article-url-feeder",
    "keyword-matcher",
    "raindrop-integration"
  ]
}

variable "VERSION" {
  default = "latest"
}

target "default" {
  dockerfile = "Dockerfile"
  annotations = [
    "org.opencontainers.image.source=https://github.com/heussd/nats-news-analysis"
  ]
  platforms = ["linux/arm64"]
}

target "golang" {
  context = "golang"
}
target "bash" {
  context = "bash"
}


target "feed-feeder" {
  inherits = ["default", "bash"]
  args = {
    MAIN = "feed-feeder"
  }
  tags = [
    "ghcr.io/heussd/nats-news-analysis/feed-feeder:latest",
    "ghcr.io/heussd/nats-news-analysis/feed-feeder:${VERSION}"
  ]
}

target "article-url-feeder" {
  inherits = ["default", "golang"]
  args = {
    MAIN = "article-url-feeder"
  }
  tags = [
    "ghcr.io/heussd/nats-news-analysis/article-url-feeder:latest",
    "ghcr.io/heussd/nats-news-analysis/article-url-feeder:${VERSION}"
  ]
}

target "keyword-matcher" {
  inherits = ["default", "golang"]

  args = {
    MAIN = "keyword-matcher"
  }
  tags = [
    "ghcr.io/heussd/nats-news-analysis/keyword-matcher:latest",
    "ghcr.io/heussd/nats-news-analysis/keyword-matcher:${VERSION}"
  ]
}

target "raindrop-integration" {
  inherits = ["default", "golang"]
  args = {
    MAIN = "raindrop-integration"
  }
  tags = [
    "ghcr.io/heussd/nats-news-analysis/raindrop:latest",
    "ghcr.io/heussd/nats-news-analysis/raindrop:${VERSION}"
  ]
}