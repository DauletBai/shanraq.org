# Two stages. The first has the whole Go toolchain and turns the source into a
# binary; the second keeps only what has to run. Everything the compiler needed
# stays behind in a layer the final image never carries.
FROM golang:1.26-alpine AS build

WORKDIR /src

# Dependencies are copied on their own first. A layer is rebuilt only when the
# files it was built from change, so editing a handler no longer re-downloads
# every module.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 is what makes the second stage possible: our SQLite driver is
# written in Go, so the binary needs no C library and runs in an empty image.
# -s -w drops the debug tables -- they help nobody in production.
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /blog .

FROM alpine:3.22

# Not root. A blog that is broken into should not be a machine that is broken
# into, and one line is the whole price of that.
RUN adduser -D -u 10001 blog

WORKDIR /app
COPY --from=build /blog /app/blog

# The database and the uploads live here. On a free plan this folder is wiped
# on every deploy -- see the lesson.
RUN mkdir -p /app/data && chown blog /app/data
ENV BLOG_DB_PATH=/app/data/blog.db
ENV BLOG_UPLOADS=/app/data/uploads

USER blog

# The port is documentation, not a rule: the host hands the real one over in
# PORT, and the blog reads it.
EXPOSE 8080

CMD ["/app/blog"]
