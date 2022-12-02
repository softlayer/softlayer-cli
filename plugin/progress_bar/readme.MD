# progress_bar
A very simple thread progress bar, which should work with any version of golang
A progress bar compatible with golang version 1.17.7 was needed, but the existing libraries only worked with versions 1.18 and 1.19, so it would have to work with any version of golang

## Usage

### Basic usage
```go
bar := progress_bar.NewProgressBar(100)
for i := 0; i < 100; i++ {
    bar.Add(1)
    time.Sleep(40 * time.Millisecond)
}
```

### Customization
There is a lot of customization that can be done: add a custom title, change the size of the progress bar, customize the progress bar symbols and additionally an option to debug and another option to print the current values.

You can instantiate and set the options in two ways, either works

```go
1.-
bar := progress_bar.NewProgressBar(numberElements).OptionTitle("Title").OptionSizeProgressBar(100).OptioProgressBarStart("(").OptioProgressBarEnd(")").OptioArrowBody("-").OptioArrowHead("}").OptionDebug(true).PrintBlankProgressBar()
```

```go
2.-
bar := progress_bar.NewProgressBar(numberElements)
bar.OptionTitle("Title")
bar.OptionSizeProgressBar(100)
bar.OptioProgressBarStart("(")
bar.OptioProgressBarEnd(")")
bar.OptioArrowBody("-")
bar.OptioArrowHead("}")
bar.OptionDebug(true)
bar.PrintBlankProgressBar()
```

```
bar.PrintOptions()
```

### Restrictions
- The size of the progress bar cannot be less than the initialized quantity, although there is a conditional that will equalize the size in case this happens
- The `PrintBlankProgressBar()` option must always be the last option set

### Possible improvements in the future
- You could add a custom output like: stout or io.Writer
- You can add custom colors for progress bar elements
- You can add a timer, to see the total elapsed time
