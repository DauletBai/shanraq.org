# Uploading images: somebody else's file on your disk

_Лид (summary):_ **The fortieth lesson of the Go course. A form with a file in it, and four things checked before that file reaches the disk: its size, its type by content, a name of our own instead of the one sent, and where to put it. And /static gets a shape at last: brand, css, images, js, uploads.**

## Why this matters

An article without a picture is half an article. Adding an upload is not hard, which is exactly why it is so often done wrong.

The file a reader sends is somebody else's bytes with somebody else's name, extension and declared type. None of those three words can be trusted, and each has a hole of its own: the name takes the file out of your folder, the extension hides a script behind a picture, and the size takes your disk and your memory.

Today all three are closed — and we see straight away that checking the type takes not one step but two: the first bytes say “looks like a picture”, not “is a picture”. Along the way `/static` is laid out the way real projects lay it out.

## The whole thing at once

A new folder, `go mod init sabaq37`:

```go
package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
)

// maxUpload is how many bytes we agree to read at all. Without this number
// anyone can take the server's memory with a single request.
const maxUpload = 1 << 20 // 1 MiB

// allowed lists the types we accept. The type is decided by the content, not
// by the extension and not by what the browser said.
var allowed = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/gif":  ".gif",
}

// maxSide caps how large a picture may claim to be. A file of a few dozen
// bytes can announce twenty thousand pixels a side, and decoding it asks for
// gigabytes of memory. The header is cheap to read, the pixels are not, so
// the size is checked before anything is decoded.
const maxSide = 8000

func main() {
	dir, _ := os.MkdirTemp("", "uploads")
	defer os.RemoveAll(dir)

	srv := httptest.NewServer(http.MaxBytesHandler(routes(dir), maxUpload))
	defer srv.Close()

	fmt.Println("== the file name from the form")
	for _, name := range []string{"photo.png", "../../../etc/passwd", "a/b/c.png", ""} {
		fmt.Printf("%-22q -> %q\n", name, filepath.Base(name))
	}

	fmt.Println()
	fmt.Println("== what the server accepts")
	fmt.Println("a real png:         ", send(srv.URL, "photo.png", smallPNG()))
	fmt.Println("html called png:    ", send(srv.URL, "photo.png", []byte("<html><script>alert(1)</script>")))
	fmt.Println("a name with a path: ", send(srv.URL, "../../evil.png", smallPNG()))
	fmt.Println("a png bomb:         ", send(srv.URL, "bomb.png", bomb()))
	fmt.Println("too large:          ", send(srv.URL, "big.png", bytes.Repeat(smallPNG(), 200000)))

	fmt.Println()
	fmt.Println("== what is left in the folder")
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		fmt.Println(" ", e.Name())
	}
}

func routes(dir string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /upload", func(w http.ResponseWriter, r *http.Request) {
		// 32 KiB stays in memory; the rest goes to a temporary file by itself.
		if err := r.ParseMultipartForm(32 << 10); err != nil {
			http.Error(w, "the file is too large", http.StatusRequestEntityTooLarge)
			return
		}
		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "no file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// The first 512 bytes decide what this really is. Neither the extension
		// nor the browser's Content-Type deserves any trust.
		head := make([]byte, 512)
		n, _ := io.ReadFull(file, head)
		kind := strings.SplitN(http.DetectContentType(head[:n]), ";", 2)[0]
		ext, ok := allowed[kind]
		if !ok {
			http.Error(w, "this is not an image: "+kind, http.StatusUnsupportedMediaType)
			return
		}
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			http.Error(w, "did not work", http.StatusInternalServerError)
			return
		}

		// The first bytes are a hint, not a proof: they are easy to forge by
		// writing a real picture’s opening into any file. So the picture’s own
		// header is decoded — it either reads or it does not — and the size it
		// announces is checked.
		cfg, _, err := image.DecodeConfig(file)
		if err != nil {
			http.Error(w, "this is not an image: "+err.Error(), http.StatusUnsupportedMediaType)
			return
		}
		if cfg.Width > maxSide || cfg.Height > maxSide {
			http.Error(w, fmt.Sprintf("the picture is too large: %dx%d", cfg.Width, cfg.Height), http.StatusRequestEntityTooLarge)
			return
		}
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			http.Error(w, "did not work", http.StatusInternalServerError)
			return
		}

		// The name is ours. What was sent is good for showing back to a person,
		// and then only after escaping.
		name := random() + ext
		path := filepath.Join(dir, name)
		dst, err := os.Create(path)
		if err != nil {
			http.Error(w, "did not work", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			// Half a file on disk is worse than none.
			os.Remove(path)
			http.Error(w, "did not work", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "saved as %s (sent as %q)", name, filepath.Base(header.Filename))
	})

	return mux
}

func random() string {
	raw := make([]byte, 12)
	rand.Read(raw)
	return base64.RawURLEncoding.EncodeToString(raw)
}

// smallPNG is a real picture one pixel in size.
func smallPNG() []byte {
	var buf bytes.Buffer
	png.Encode(&buf, image.NewRGBA(image.Rect(0, 0, 1, 1)))
	return buf.Bytes()
}

// bomb is an honest PNG by its first bytes: it announces twenty thousand
// pixels a side and weighs seventy-four bytes.
func bomb() []byte {
	raw := smallPNG()
	binary.BigEndian.PutUint32(raw[16:20], 20000) // the width in the IHDR header
	binary.BigEndian.PutUint32(raw[20:24], 20000) // the height
	binary.BigEndian.PutUint32(raw[29:33], crc32.ChecksumIEEE(raw[12:29]))
	return raw
}

func send(base, name string, body []byte) string {
	var buf bytes.Buffer
	form := multipart.NewWriter(&buf)
	part, err := form.CreateFormFile("image", name)
	if err != nil {
		return "error: " + err.Error()
	}
	part.Write(body)
	form.Close()

	resp, err := http.Post(base+"/upload", form.FormDataContentType(), &buf)
	if err != nil {
		return "error: " + err.Error()
	}
	defer resp.Body.Close()
	out, _ := io.ReadAll(resp.Body)
	return fmt.Sprintf("%d — %s", resp.StatusCode, strings.TrimSpace(string(out)))
}
```

The output:

```
== the file name from the form
"photo.png"            -> "photo.png"
"../../../etc/passwd"  -> "passwd"
"a/b/c.png"            -> "c.png"
""                     -> "."

== what the server accepts
a real png:          200 — saved as pM361aTkUKWtjqjC.png (sent as "photo.png")
html called png:     415 — this is not an image: text/html
a name with a path:  200 — saved as s035PI2jgJrUCVyT.png (sent as "evil.png")
a png bomb:          413 — the picture is too large: 20000x20000
too large:           413 — the file is too large

== what is left in the folder
  pM361aTkUKWtjqjC.png
  s035PI2jgJrUCVyT.png
```

## Taking it apart

### `multipart/form-data`: a form with a file in it

An ordinary form sends "name — value" pairs in one line. A file cannot go that way, so a form with a file looks different:

```html
<form method="post" action="/upload" enctype="multipart/form-data">
  <input type="file" name="image" accept="image/*">
  <button>upload</button>
</form>
```

The `enctype` is compulsory here. Leave it out and the browser sends only the name of the file with no content, and `FormFile` returns an error people then hunt for a long time.

In Go such a form is parsed by `ParseMultipartForm`, and `r.FormFile("image")` gives three things: the file, a header with the sent name and size, and an error.

### The cap goes before the reading, not after

`http.MaxBytesHandler` wraps the server and cuts the reading off at a given number of bytes. The output shows how that ends for a file that is too large: `413`.

Putting a cap **after** the reading is pointless: to learn the size the file has to be read first, which means memory or disk has already been spent on it. Checking `header.Size` is checking what the browser said, and it guarantees nothing.

`ParseMultipartForm(32 << 10)` is a second number and it is about something else: how much to keep in memory. Anything larger Go puts into a temporary file itself and clears up afterwards.

> **Picture it.** A turnstile at the entrance rather than a guard at the end of the corridor. The turnstile does not let through what does not fit; the guard at the end finds it when it is already inside.

### You make the file name up

The first part of the output is about the name. `filepath.Base` cuts everything but `passwd` off `../../../etc/passwd`, which is not bad. But leaning on it alone will not do: look at the last line — an empty name becomes `.`, and you would create such a file in silence.

So the name is not cleaned but **made up**: twelve random bytes plus the extension we chose ourselves. The name that was sent is good for showing back to a person, and then only after escaping, as in the XSS lesson.

Several troubles go away at once: the file cannot leave the folder, cannot overwrite somebody else's, cannot turn out to be `.htaccess` or `index.html`, and its name holds no spaces, no Cyrillic and no semicolons.

### The type is decided by the content

The second line of the output is the main one. A file named `photo.png` holding `<html><script>` is refused with `415`, because we asked about the content rather than the name:

```go
head := make([]byte, 512)
n, _ := io.ReadFull(file, head)
kind := strings.SplitN(http.DetectContentType(head[:n]), ";", 2)[0]
```

`http.DetectContentType` looks at the first bytes and answers by what is really there: a PNG begins with `\x89PNG`, a JPEG with `\xFF\xD8\xFF`. The five hundred and twelve bytes are not our invention but the size the algorithm is built for.

The extension cannot be trusted: it is written by whoever sends the file. Nor can the `Content-Type` from the form: that is written by the browser, and a browser need not be involved at all.

And **a white list rather than a black one**: we list what we accept, not what we forbid. It is the same rule as with sorting in the injection lesson, and it works for the same reason — a stranger's imagination is richer than yours.

After the check do not forget `file.Seek(0, io.SeekStart)`: five hundred and twelve bytes have been read, and without rewinding the file is saved with its head missing.

### The first bytes are a hint, not a proof

`DetectContentType` looks at the signature, and a signature can be written in: take the opening of a real png, append anything you like — and the check passes. So a second, honest one follows it: the standard library decodes the picture’s own header.

```go
cfg, _, err := image.DecodeConfig(file)
```

`image.DecodeConfig` reads only the header — the width, the height and the colour model — and never touches the pixels. It either parses the file or returns an error, and a forged opening no longer gets through: a png has to match its checksum, a jpeg its segment structure.

It also answers a question the file size leaves open. The output has a line `a png bomb: 413`. The file weighs seventy-four bytes and looks genuine by every sign, but its header says `20000x20000`. Decoding such a picture means asking the system for gigabytes of memory on a single request. That is why the announced size is checked **before** anything is decoded:

```go
if cfg.Width > maxSide || cfg.Height > maxSide {
```

The same explains why `webp` is not on the list: the standard library has no decoder for it, and accepting a format we cannot check means trusting a stranger’s word again.

And after the check the file is rewound once more: the header has been read, and the whole file has to land on disk.

### The `/static` tree: `brand`, `css`, `images`, `js`, `uploads`

From this step on the blog keeps its files in folders, and it stays that way. While there were three of them they sat side by side quite happily — making folders for three files is tidying where there is no mess yet. Today a fourth kind of file appears, somebody else's, and the rule changes.

It changes once and for the future. Fonts, icons, a script for the search, placeholder pictures and article covers are all coming, and each kind will take a folder of its own rather than joining one heap. Sorting forty files afterwards costs more than making five folders now:

```
static/
    brand/      the logo and the icons: what makes the site recognisable
        logo.svg
    css/
        style.css
    js/
        main.js
    images/     the site's own pictures: illustrations, placeholders
        empty.png
    uploads/    what readers sent
        2026/09/xxxx.png
```

The rule is simple: **ours and theirs live in different folders**. `brand`, `css`, `js`, `images` are ours and travel with the program. `uploads` is other people's and appears while it runs.

The split is not about tidiness. What is ours can be embedded into the program and served with a long cache; what is theirs cannot — `embed` works at build time and an upload does not exist then. And the safety rules differ: our files we wrote ourselves, theirs we did not.

Inside `uploads` it pays to sort by year and month: a thousand files in one folder is already awkward, a hundred thousand is already slow.

### Serve other people's files carefully

Serving what was uploaded is a separate place where mistakes are made.

An `http.FileServer` over the `uploads` folder is fine, but the `X-Content-Type-Options: nosniff` header from the previous lesson is compulsory here: without it a browser may decide your "png" is really HTML and run it.

Beyond that it depends. Large sites serve uploads from a domain of their own, so that a script that does get executed is "foreign" to the main site and sees neither its cookies nor its pages. For a blog that is too much; `nosniff` and a white list of types are the necessary minimum.

## The lesson map

![Lesson map: size, type, name and place](/static/course/go/map-upload-en.svg)

## Say it in your own words

Without looking, answer out loud or on paper. The answers are at the end of the lesson.

1. Why is the size capped by a wrapper rather than by checking `header.Size`?
2. Why is the file name made up afresh when `filepath.Base` has already cut the path off?
3. Why is the type decided by the first bytes rather than by the extension?

## Exercise

**Required.** Add a cover to an article in the blog. The form gets an `enctype`, the handler gets a cap, a check of the type by content, a decode of the picture’s header with a limit on the side, and a name of its own. Files land in `static/uploads`, and `/static` spreads into folders at the same time: `brand`, `css`, `js`, `images`, `uploads`.

All of it is done in [step-25](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-25) — compare against it once you have written your own.

**Optional.**

- Send a file named `../../main.go` with `curl` and make sure it lands in `uploads` under its new name.
- Sort `uploads` by year and month, and write down what to do with what is already there.
- Shrink an upload to a sensible width before saving it. You will need the `image` package from the standard library.

## Where this goes in your blog

An article gets a cover, and the blog gets its first folder whose contents we did not write. From now on `/static` lives by a rule: ours apart, theirs apart.

Debts. Old files are not deleted along with an article. Pictures are not re-encoded: whatever was sent is what sits there, camera metadata and all. And we serve them from the same domain as the pages.

## Answers

1. Because `header.Size` is what the browser said, and the real size can only be learned by reading the file, which means memory or disk has already been spent on it. The wrapper cuts the reading off at the byte you name.
2. Because `Base` solves only the problem with the path. What is left is an empty name that becomes `.`, a collision with somebody else's file, dangerous names such as `index.html`, and whatever else the sender thinks of. A name of your own removes that whole class of questions at once.
3. Because the extension is written by whoever sends the file, and so is the declared `Content-Type`. The first bytes belong to the file itself rather than to words about it. But they are a hint too: the real answer comes from decoding the picture’s header.

## Sources

- [Go: the mime/multipart package](https://pkg.go.dev/mime/multipart)
- [Go: DetectContentType](https://pkg.go.dev/net/http#DetectContentType)
- [Go: image.DecodeConfig](https://pkg.go.dev/image#DecodeConfig)
- [OWASP: file upload](https://cheatsheetseries.owasp.org/cheatsheets/File_Upload_Cheat_Sheet.html)
