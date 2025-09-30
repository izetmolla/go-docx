<h1 align="center">go-docx</h1>

<p align="center">

[![tests](https://github.com/izetmolla/go-docx/workflows/test/badge.svg)](https://github.com/izetmolla/go-docx/actions?query=workflow%3Atest)
[![goreport](https://goreportcard.com/badge/github.com/izetmolla/go-docx)](https://goreportcard.com/report/github.com/izetmolla/go-docx)
[![GoDoc reference ](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/izetmolla/go-docx)

</p>

<p align="center">
   <b>Replace placeholders inside docx documents with speed and confidence.</b>


   <p align="center"><sub>This project provides a simple and clean API to perform replacing of user-defined placeholders.
   Without the uncertainty that the placeholders may be ripped  apart by the WordprocessingML engine used to create the document.</sub></p>
</p>


<p align="center">
  <img src="https://github.com/izetmolla/go-docx/blob/refactoring/screenshot.png" alt="Example" width="900" />
</p>


* **Simple**: The API exposed is kept to a minimum in order to stick to the purpose.
* **Fast**: go-docx is fast since it operates directly on the byte contents instead mapping the XMLs to a custom data struct.
* **Zero dependencies**: go-docx is build with the go stdlib only, no external dependencies.
* **Serverless Ready**: Built-in support for serverless processing with MinIO-to-MinIO workflows.

---

### Table of Contents
- [Purpose](#-purpose)
- [Getting Started](#-getting-started)
- [Serverless Processing](#-serverless-processing)
- [Examples](#-examples)
- [Terminology](#-terminology)
- [How it works](#-how-it-works)
- [License](#-license)

---

### ➤ Purpose
The task at hand was to replace a set of user-defined placeholders inside a docx document archive with calculated values.
All current implementations in Golang which solve this problem use a naive approach by attempting to `strings.Replace()` the placeholders.

Due to the nature of the WordprocessingML specification, a placeholder which is defined as `{the-placeholder}` may be ripped apart inside the resulting XML.
The placeholder may then be in two fragments for example `{the-` and `placeholder}` which are spaced apart inside the XML.

The naive approach therefore is not always working. To provide a way to replace placeholders, even if they are fragmented, is the purpose of this library.

### ➤ Getting Started
All you need is to `go get github.com/izetmolla/go-docx`

```go
func main() {
        // replaceMap is a key-value map whereas the keys
	// represent the placeholders without the delimiters
	replaceMap := docx.PlaceholderMap{
		"key":                         "REPLACE some more",
		"key-with-dash":               "REPLACE",
		"key-with-dashes":             "REPLACE",
		"key with space":              "REPLACE",
		"key_with_underscore":         "REPLACE",
		"multiline":                   "REPLACE",
		"key.with.dots":               "REPLACE",
		"mixed-key.separator_styles#": "REPLACE",
		"yet-another_placeholder":     "REPLACE",
	}

        // read and parse the template docx
	doc, err := docx.Open("template.docx")
	if err != nil {
	    panic(err)
	}

        // replace the keys with values from replaceMap
	err = doc.ReplaceAll(replaceMap)
	if err != nil {
	    panic(err)
	}

        // read the replacement image
	bs, err := os.ReadFile("./test/cameraman.jpg")
	if err != nil {
	    panic(err)
	}

        // replace the lenna image with the cameraman image
        // image attributes, e.g. size, position, keep unchanged
	err = doc.SetFile("word/media/image1.jpg", bs)
	if err != nil {
	    panic(err)
	}

        // write out a new file
	err = doc.WriteToFile("replaced.docx")
	if err != nil {
		panic(err)
	}
}
```

#### Placholders
Placeholders are delimited with `{` and `}`, nesting of placeholders is not possible.
Placeholders can be changed using `ChangeOpenCloseDelimiter()`.

#### Styling
The way this lib works is that a placeholder is just a list of fragments. When detecting the placeholders inside the XML, it looks for the OpenDelimiter and CloseDelimiter.
The first fragment found (e.g. `{foo` of placeholder `{foo-bar}`) will be replaced with the value from the `ReplaceMap`.

This means that technically you can style only the OpenDelimiter inside the Word-Document and the whole value will be styled like that after replacing.
Although I do not recommend to do that as the WordprocessingML spec is somewhat fragile in this case. So it's best to just style the whole placeholder.

But, for whatever reason there might be, you can do that.

#### Image replace
Image replacing is slightly different from text replacing. To replace an image, you need to know its path within the docx archive, rather than using a placeholder.

To find the path, use unzip or similar tools to extract the contents of a docx file, then locate the image to be replaced inside the `word/media/` folder. Assume the path is `word/media/image1.jpg`, then you can use the `SetFile()` to overwrite the old image with a new one. It should be noted that:
- The image format (encoding) should keep the same during the replacement.
- Since the metadata of the image is not changed, only the image file itself is replaced, the new image will appear in its original location, with its original size. In other words, the image attributes keep unchanged.

### ➤ Serverless Processing

For serverless environments (AWS Lambda, Google Cloud Functions, Azure Functions, etc.) where no file system is available, go-docx provides dedicated functions for processing documents entirely in memory.

#### Basic Serverless Processing

```go
// Process template bytes with data (no file system involved)
processedBytes, err := docx.CompleteTemplateFromBytesToBytes(templateData, data)
if err != nil {
    log.Fatal(err)
}
```

#### With Custom Functions

```go
// Define custom functions for data processing
funcMap := map[string]interface{}{
    "uppercase": func(s string) string {
        return strings.ToUpper(s)
    },
    "format_currency": func(amount string) string {
        return fmt.Sprintf("USD %s", amount)
    },
}

// Process with custom functions
processedBytes, err := docx.CompleteTemplateFromBytesToBytesWithFuncs(templateData, data, funcMap)
if err != nil {
    log.Fatal(err)
}
```

#### MinIO-to-MinIO Example

```go
// Download template bytes from MinIO
templateBytes, err := minioClient.GetObject(bucketName, "templates/report-template.docx", minio.GetObjectOptions{})
if err != nil {
    log.Fatal(err)
}
defer templateBytes.Close()

templateData, err := ioutil.ReadAll(templateBytes)
if err != nil {
    log.Fatal(err)
}

// Process template bytes with data (no file system involved)
processedBytes, err := docx.CompleteTemplateFromBytesToBytes(templateData, data)
if err != nil {
    log.Fatal(err)
}

// Upload processed bytes back to MinIO
minioClient.PutObject(
    bucketName, 
    "processed/report.docx", 
    bytes.NewReader(processedBytes), 
    int64(len(processedBytes)), 
    minio.PutObjectOptions{},
)
```

**Benefits of Serverless Processing:**
- ✅ **No File System**: Perfect for serverless environments
- ✅ **Memory Efficient**: Processes documents entirely in memory
- ✅ **Scalable**: Can handle high concurrency in serverless environments
- ✅ **Cost Effective**: Pay only for processing time
- ✅ **Fast**: No disk I/O overhead

### ➤ Examples

The repository includes several examples to help you get started:

#### Basic Example
```bash
cd examples/simple
go run main.go
```

#### Complex Example
```bash
cd examples/complex
go run main.go -template template.docx -out output.docx
```

#### Serverless Example (MinIO-to-MinIO)
```bash
cd examples/serverless
go mod tidy
go run main.go
```

**Prerequisites for Serverless Example:**
1. Start MinIO server:
```bash
docker run -p 9000:9000 -p 9001:9001 \
  -e "MINIO_ROOT_USER=minioadmin" \
  -e "MINIO_ROOT_PASSWORD=minioadmin" \
  minio/minio server /data --console-address ":9001"
```

2. Upload a DOCX template to MinIO at `templates/report-template.docx`

#### Benchmark Example
```bash
cd examples/benchmark
go test -bench=.
```

### ➤ Terminology
To not cause too much confusion, here is a list of terms which you might come across.

* **Parser**: Every file which this lib handles (document, footers and headers) has their own parser attached since everything is relative to the underlying byte-slice (aka. file).
* **Position**: A Position is just a `Start` and `End` offset, relative to the byte slice of the document of a parser.
* **Run**: Describes the pair `<w:r>` and `</w:r>` and thus has two `Positions` for the open and close tag. Since they are Positions, they have a `Start` and `End` Position which point to `<` and `>` of the tag. A run also consists of a `TagPair`.

* **Placeholder**: A Placeholder is basically just a list of `PlaceholderFragments` representing a full placeholder extracted by a `Parser`.
* **PlaceholderFragment**: A PlaceholderFragment is a parsed fragment of a placeholder since those will most likely be ripped apart by WordprocessingML. The Placeholder `{foo-bar-baz}` might ultimately consist of 5 fragments ( `{`, `foo-`, `bar-`, `baz`, `}`).
The fragment is at the heart of replacing. It knows to which `Run` it belongs to and has methods of manipulating these byte-offsets. Additionally it has a `Position` which describes the offset inside the `TagPair` since the fragments don't always start at the beginning of one (e.g. `<w:t>some text {fragment-start</w:t>`)

### ➤ How it works
This section will give you a short overview of what's actually going on.
And honenstly.. it's a much needed reference for my future self :D.

#### Overview
The project does rely on some invariants of the WordprocessingML spec which defines the docx structure.
A good overview over the spec can be found on: [officeopenxml.com](http://officeopenxml.com/anatomyofOOXML.php).

Since this project aims to work only on text within the document, it currently only focuses on the **runs** (`<w:r>` element).
A run *always* encloses a **text** (`<w:t>` element) thus finding all runs inside the docx is the first step. Keep in
mind that a run does not need to have a text element. It can also contain an image for example. But all text
literals will always be inside a run, within their respective text tags.

To illustrate that, here is how this might look inside the document.xml.

```xml
 <w:p>
    <w:r>
        <w:t>{key-with-dashes}</w:t>
    </w:r>
</w:p>
```
One can clearly see that replacing the `{key-with-dashes}` placeholder is quite simple.
Just do a `strings.Replace()`, right? **Wrong!**

Although this might work on 70-80% of the time, it will not work reliably.
The reason is how the WordprocessingML spec is set-up. It will fragment text-literals
based on many different circumstances.

For example if you added half of the placeholder, saved
and quit Word, and then add the second half of the placeholder, it might happen (in order to preserve session history), that the placeholder will look something like that (simplified).

```xml
 <w:p>
    <w:r>
        <w:t>{key-</w:t>
    </w:r>
    <w:r>
        <w:t>with-dashes}</w:t>
    </w:r>
</w:p>
```

As you can clearly see, doing a simple replace doesn't do it in this case.

#### Premises
In order to achive the goal of reliably replacing values inside a docx archive, the following
premises are considered:

* Text literals are always inside `<w:t>` tags
* `<w:t>` tags only occur inside `<w:r>` tags
* All placeholders are delimited with predefined runes (`{` and `}` in this case)
* Placeholders cannot be nested (e.g. `{foo {bar}}`)

#### Order of operations
Here I will outline what happens in order to achieve the said goal.

1. Open the *.docx file specified and extract all files in which replacement should take place.
 Currently, there files extracted are `word/document.xml`, `word/footer<X>.xml` and `word/header<X>.xml`.
 Any content which resides in different files requires a modification.

2. First XML pass. Iterate over a given file (e.g. the document.xml) and find all `<w:r>` and `</w:r>` tags inside
the bytes of the file. Remember the positions given by the custom `io.Reader` implementation.
**Note** Singleton tags are handled correctly (e.g. `<w:r/>`).

3. Second XML pass. Basically the same as the first pass, just this time the text tags (`<w:t>`) inside
the found runs are extracted.

4. Placeholder extraction. At this point all text literals are known by their offset inside the file.
Using the premise that no placeholder nesting is allowed, the placeholder fragments can be
extracted from the text runs. At the end a placeholder may be described by X fragments.
The result of the extraction is the knowledge of which placeholders are located inside the document and
at which positions the fragments start and end.

5. Making use of the positions and replace some content. This is the stage where all the placeholders
need to be replaced by their expected values given in a `PlaceholderMap`. The process can rougly
be outlined in two steps:
    * The first fragment of the placeholder (e.g. `{foo-`) is replaced by the actual value. This
    also explains why one only has to style the first fragment inside the document. As you cannot
    see the fragments it is still a good idea to style the whole placeholder as needed.
    * All other fragments of the placeholders are cut out, removing the leftovers.

All the steps taken in 5. require cumbersome shifting of the offsets. This is the tricky part
where the most debugging happened (gosh, so many offsets). The given explanation is definitely enough
to grasp the concept, leaving out the messy bits.

### ➤ License
This software is licensed under the [MIT license](https://github.com/izetmolla/go-docx/blob/develop/LICENSE).
