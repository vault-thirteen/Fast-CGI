# Fast CGI Compatibility Layer
### For Go Programming Language
![FastCGI Logotype](img/Logo_GreyBg_330x200.png)

## <a name="section-0" id="section-0">Contents</a>
* [Description](#section-1)
* [Repository Structure](#section-2)
* [Usage](#section-3)
* [Why ?](#section-4)
* [Résumé](#section-5)
* [Feedback](#section-6)

## <a name="section-1" id="section-1">Description</a>

This repository stores documentation of CGI and FastCGI interfaces and 
programming tools which may be useful for building applications using the 
FastCGI interface in Go programming language, a.k.a. Golang.

These tools include a library written in Go language. This library provides 
objects, methods and functions to work with FastCGI interface in Go programming 
language.

The library implements methods for a FastCGI <b>client</b>, mostly.

Implementation of a FastCGI <b>server</b> can be found in a built-in standard library
of Go language – `net/http/fcgi` package, https://pkg.go.dev/net/http/fcgi.

The library provides a simple experimental web server to use with legacy 
scripts supporting the CGI and FastCGI interfaces. This server should not be 
used anywhere except some experiments with PHP scripts because it has several 
issues which are described further in the [Résumé](#section-5) section.

## <a name="section-2" id="section-2">Repository Structure</a>

**N.B.** *Due to some bugs in Go language, the structure of this repository is 
heavily modified to meet the requirements of the Google's Golang proxy server 
which often throws funny messages. The most funny of them are messages stating 
that a package was downloaded but not found. If you are, like me, used to a 
strict repository layout which separates source code into an `src` folder and 
other stuff – into other folders, this structure is forbidden by Google's 
Golang.* 

So, knowing all the above, the structure of this repository is as 
follows:

* [DOC](doc) folder contains the documentation.
* Other folders contain various parts of the library for Go language.

## <a name="section-3" id="section-3">Usage</a>

Usage examples can be found in the ['example'](example) folder.

The most simple usage example requires only a few lines of code and is as
simple as the following code:
```go
package main

import (
  "fmt"
  "os"

  "github.com/vault-thirteen/Fast-CGI/pkg/models/php"
)

func main() {
  var err error
  err = runSimplePhpScript(`D:\Scripts\script.php`)
  if err != nil {
    panic(err)
  }
}

func runSimplePhpScript(scriptFilePath string) (err error) {
  var stdOut, stdErr []byte
  stdOut, stdErr, err = pm.RunOnceSimplePhpScript("tcp", "127.0.0.1:9000", scriptFilePath)
  if err != nil {
    return err
  }

  if len(stdErr) > 0 {
    _, err = fmt.Fprintln(os.Stderr, string(stdErr))
    if err != nil {
      return err
    }
  }

  if len(stdOut) > 0 {
    _, err = fmt.Fprintln(os.Stdout, string(stdOut))
    if err != nil {
      return err
    }
  }

  return nil
}

```

For more complex tasks, the `Client` object and its methods can be used.

## <a name="section-4" id="section-4">Why ?</a>

<b>Reason 1</b>

> <b>FastCGI</b> interface provides a cross-platform and
cross-language interface for communication between totally different programming
languages. For example, with the help of <b>FastCGI</b> it is possible to call
a function of a PHP language from the code written in Go language, where Go is
a compiled language and PHP is an interpreted language with no runtime
"engine".

<b>Reason 2</b>

> <b>FastCGI</b> interface is still actual for simple tasks and
applications where client code is not a bottleneck, e.g. for applications using
a lot of complicated SQL queries, such as bulletin boards and online forums.

<b>Reason 3</b>

> There are a lot of programming languages which still support the
<b>FastCGI</b> interface. And many of these languages are mature and
are sufficient for their purposes.

<b>Reason 4</b>

> This library makes it possible to bring a new life into old projects
written in other programming languages that support the <b>FastCGI</b>
interface. It is not necessary to re-write PHP projects in Go from scratch.
This is why it is said in this repository, that this project is a compatibility
layer.


## <a name="section-5" id="section-5">Résumé</a>

1. The **FastCGI** interface itself is not so bad even with all its drawbacks.  

   *  It uses 16-bit fields for data transmission making it practically useless 
      for HTTP bodies being longer than 65535 bytes, but it does its job when 
      you need to connect totally different systems together. The main problem 
      lies much deeper than FastCGI.


2. The main problem is the **CGI** interface which is even older than 
**FastCGI** interface and has several serious issues.
<br><br>

   *  First of all, CGI was made for UNIX operating systems. It was using 
      so-called UNIX sockets to make it work faster than a turtle. Other 
      operating systems, such as Microsoft Windows, do not have UNIX sockets 
      and use TCP network protocol for inter-process communication. This 
      approach itself decreases the already slow protocol and makes it even 
      slower.
<br><br>

   *  CGI on UNIX uses forward slash symbols in all paths – in URL and in paths 
      inside an operating system. This creates difficulties for non-UNIX 
      operating systems.
<br><br>

   *  The `Extra Path` feature of the CGI interface is very dumb and dangerous 
      at the same time. Yes, this is **D&D**, but it is not about dragons this 
      time. First of all, it breaks URL parsing completely, because there is no 
      way to distinguish files and folders by only looking at their names. So, 
      the URL `http://some.machine/cgi-bin/display.pl/cgi/cgi_doc.txt`
      in real life may be a Perl script `display.pl` with a CGI "extra path", 
      or it may be a text file `cgi_doc.txt` sitting inside a folder with dot 
      symbol in its name. This "feature" makes CGI unusable in real-world 
      applications because it forces owners of web servers to invent stupid 
      mechanisms for guessing the real meaning of a URL. All this situation
      introduces at least two consequences: poor performance of a web server 
      and a huge hole in the security of an entire system. This `Extra Path` 
      feature should be disabled every time except those corner cases when you 
      really need it to experiment with some legacy code.

## <a name="section-6" id="section-6">Feedback</a>
If you have any feedback, you are free to direct it to this GitHub repository:
* Comments should be written to the `Discussions` section:
  [here](https://github.com/vault-thirteen/Fast-CGI/discussions)
* Bug reports should be directed to the `Issues` section:
  [here](https://github.com/vault-thirteen/Fast-CGI/issues)
