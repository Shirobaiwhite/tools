As of Go ver. 1.16, ioutil is deprecated and its methods were inherited by io and os.

The package io performs more sophisticated tasks related to input and output such as `WriteString()` which writes a string to a writer, `CopyN()` which copies exactly n bytes from src to dst; while the package os performs system-level tasks more such as `Chmod()` and `Mkdir()`.

 

Here is a definition of `WriteFile()` in the `os` package:
```Go

// WriteFile writes data to the named file, creating it if necessary.
// If the file does not exist, WriteFile creates it with permissions perm (before umask);
// otherwise WriteFile truncates it before writing, without changing permissions.

func WriteFile(name string, data []byte, perm FileMode) error {
	f, err := OpenFile(name, O_WRONLY|O_CREATE|O_TRUNC, perm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}
```

As we can see, there are basically 3 steps for `WriteFile()`. First, we open the file; then we write it; then we close the file. There isn’t much room for us to squeeze in something that can track how many bytes of data we have already written into the file. The `Write()` method in the same package displays the same properties:

```Go
// Write writes len(b) bytes from b to the File.
// It returns the number of bytes written and an error, if any.
// Write returns a non-nil error when n != len(b).

func (f *File) Write(b []byte) (n int, err error) {
	if err := f.checkValid("write"); err != nil {
		return 0, err
	}
	n, e := f.write(b)
	if n < 0 {
		n = 0
	}
	if n != len(b) {
		err = io.ErrShortWrite
	}

	epipecheck(f, e)

	if e != nil {
		err = f.wrapErr("write", e)
	}

	return n, err
}
```

Although we have a more detailed description of the number of bytes written, it is something we can only get when the function returns.

 

The problem could be resolved if we start a goroutine that constantly stats the file we are writing into. For example, the goroutine here is used to track the progress of the `WriteFile()` operation:

```Go
...
  go func() {
		for {
			f, err := os.Stat(file_path)
			if err == nil {
				progress := f.Size() * 100 / int64(len(b))
                ...
				if progress > 99 {
					return
				}
			}
			time.Sleep(15 * time.Second)
		}
	}()

	if err := os.WriteFile(file_path, b, 0644); err != nil {
		log.Println(err)
		return
	}
	defer os.Remove(file_path)
...
```

This was easily done because we have the length of b, the total bytes being written into the file at `file_path`. Note that goroutines end with the main function.

 

However, if we are using HTTP requests, this method will be no longer feasible because we cannot stat the size of the data posted(fetched). 

``` Go
...
    client := &http.Client{}
	data, err := os.Open(file_path)
	defer data.Close()
	defer os.Remove(file_path)
	if err != nil {
		log.Println(err)
		return err
	}
	req, err := http.NewRequest("PUT", url, data)
	...
	_, err = client.Do(req)
	...
...

```
An example will be the code above. The method for tracking how the HTTP client is Do()ing is to wrap the file in a reader.

```Go
type progressReader struct {
	r     io.Reader
	max   int
	sent  int
	atEOF bool
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.r.Read(p)
	pr.sent += n
	if err == io.EOF {
		pr.atEOF = true
	}
	pr.report()
	return n, err
}

func (pr *progressReader) report() {
	fmt.Printf("sent %d of %d bytes, %d%%\n", pr.sent, pr.max, int(math.Floor(float64(pr.sent)*100/float64(pr.max))))
	if pr.atEOF {
		fmt.Println("DONE")
	}
}
```

Then, if we change the `io.Reader` in the HTTP request to

```Go
&progressReader{r: data, max: int(file.Size())}
```
This case, we can track the progress each time `Read()` is called with `report()`. 
